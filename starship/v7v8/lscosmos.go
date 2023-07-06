package main

import (
	"context"
	"fmt"
	"time"

	starship "github.com/cosmology-tech/starship/clients/go/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	lscosmostypes "github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

func (s *TestSuite) RunLSCosmosTests() {
	persistence := s.GetChainClient("test-core-1")
	gaia := s.GetChainClient("test-gaia-1")
	uatom := gaia.MustGetChainDenom()
	c := s.GetTransferChannel(persistence, gaia.ChainID)

	// Setup lscosmos
	s.T().Log("submitting fee change proposal")
	s.SubmitFeeChangeProposal(persistence)
	s.T().Log("running jump-start tx")
	s.JumpStart(persistence, gaia, c)

	s.T().Log("transferring 100 ATOMs to persistence address")
	s.IBCTransferTokens(gaia, persistence.Address, c.PortId, c.ChannelId, uatom, 100000000)

	// Liquid stake
	s.T().Log("liquid staking 50 ATOMs")
	s.LiquidStake(persistence, c.PortId, c.ChannelId, uatom, 50000000)

	// Liquid unstake
	s.T().Log("liquid unstaking 5 ATOMs")
	s.LiquidUnstake(persistence, uatom, 5000000)
}

func (s *TestSuite) SubmitFeeChangeProposal(persistence *starship.ChainClient) {
	content := &lscosmostypes.PstakeFeeAddressChangeProposal{
		Title:            "pstake fee address change proposal",
		Description:      "pstake fee address change proposal",
		PstakeFeeAddress: persistence.Address,
	}
	s.SubmitAndPassProposal(persistence, content, "pstake fee address change")
}

func (s *TestSuite) JumpStart(persistence, gaia *starship.ChainClient, c *ibcchanneltypes.IdentifiedChannel) {
	valAddr := s.GetChainValAddr(gaia)
	uatom := gaia.MustGetChainDenom()
	msg := &lscosmostypes.MsgJumpStart{
		PstakeAddress:   persistence.Address,
		ChainID:         gaia.ChainID,
		TransferChannel: c.Counterparty.ChannelId,
		TransferPort:    c.Counterparty.PortId,
		ConnectionID:    c.ConnectionHops[0],
		BaseDenom:       uatom,
		MintDenom:       fmt.Sprintf("stk/%s", uatom),
		MinDeposit:      sdk.NewInt(1),
		AllowListedValidators: lscosmostypes.AllowListedValidators{
			AllowListedValidators: []lscosmostypes.AllowListedValidator{{
				ValidatorAddress: valAddr,
				TargetWeight:     sdk.NewDec(1),
			}},
		},
		PstakeParams: lscosmostypes.PstakeParams{
			PstakeDepositFee:    sdk.NewDec(0),
			PstakeRestakeFee:    sdk.NewDecWithPrec(5, 2),
			PstakeUnstakeFee:    sdk.NewDec(0),
			PstakeRedemptionFee: sdk.NewDecWithPrec(1, 1),
			PstakeFeeAddress:    persistence.Address,
		},
		HostAccounts: lscosmostypes.HostAccounts{
			DelegatorAccountOwnerID: "lscosmos_pstake_delegation_account",
			RewardsAccountOwnerID:   "lscosmos_pstake_reward_account",
		},
	}
	s.SendMsgAndWait(persistence, msg, "pstake jump-start")
	s.WaitForLSCosmosToBeEnabled(persistence)
}

func (s *TestSuite) LiquidStake(chain *starship.ChainClient, port, channel, denom string, amt int64) {
	ibcAtom := s.GetIBCDenom(chain, port, channel, denom)
	lsAmount := sdk.NewCoin(ibcAtom, sdk.NewInt(amt))
	s.SendMsgAndWait(chain, &lscosmostypes.MsgLiquidStake{
		DelegatorAddress: chain.Address,
		Amount:           lsAmount,
	}, "liquid stake")

	s.T().Log("wait for delegation")
	delAmount := s.WaitForDelegation(chain)
	s.Require().Equal(lsAmount.Amount, delAmount.Amount, "liquid staked amount do not match")
}

func (s *TestSuite) LiquidUnstake(chain *starship.ChainClient, denom string, amt int64) {
	mintDenom := fmt.Sprintf("stk/%s", denom)
	lsAmount := sdk.NewCoin(mintDenom, sdk.NewInt(amt))
	s.SendMsgAndWait(chain, &lscosmostypes.MsgLiquidUnstake{
		DelegatorAddress: chain.Address,
		Amount:           lsAmount,
	}, "liquid unstake")

	s.T().Log("wait for undelegation")
	undelAmount := s.WaitForUnDelegation(chain, sdk.NewInt(amt))
	s.Require().Equal(lsAmount, undelAmount, "total undelegation amount do not match")
}

func (s *TestSuite) WaitForLSCosmosToBeEnabled(chain *starship.ChainClient) {
	s.Require().Eventuallyf(
		func() bool {
			res, err := lscosmostypes.
				NewQueryClient(chain.Client).
				ModuleState(context.Background(), &lscosmostypes.QueryModuleStateRequest{})
			s.Require().NoError(err)
			return res.ModuleState
		},
		300*time.Second,
		time.Second,
		"waited for too long, lscosmos module-state is still false",
	)
}

func (s *TestSuite) WaitForDelegation(chain *starship.ChainClient) sdk.Coin {
	var delAmount *sdk.Coin
	s.Require().Eventuallyf(
		func() bool {
			res, err := lscosmostypes.
				NewQueryClient(chain.Client).
				DelegationState(context.Background(), &lscosmostypes.QueryDelegationStateRequest{})
			s.Require().NoError(err)
			if res == nil {
				return false
			}
			dels := res.DelegationState.HostAccountDelegations
			if len(dels) == 0 {
				return false
			}
			delAmount = &dels[0].Amount
			return true
		},
		300*time.Second,
		time.Second,
		"waited for too long, tokens not yet liquid staked",
	)
	s.Require().NotNil(delAmount)
	return *delAmount
}

func (s *TestSuite) WaitForUnDelegation(chain *starship.ChainClient, expectedAmount sdk.Int) sdk.Coin {
	var undelAmount *sdk.Coin
	s.Require().Eventuallyf(
		func() bool {
			res, err := lscosmostypes.
				NewQueryClient(chain.Client).
				DelegationState(context.Background(), &lscosmostypes.QueryDelegationStateRequest{})
			s.Require().NoError(err)
			if res == nil {
				return false
			}
			undels := res.DelegationState.HostAccountUndelegations
			if len(undels) == 0 {
				return false
			}
			undelAmount = &undels[0].TotalUndelegationAmount
			// If completion time is set (i.e. not 0), it means tokens undelegated successfully
			return !undels[0].CompletionTime.IsZero()
		},
		300*time.Second,
		time.Second,
		"waited for too long, tokens not yet liquid unstaked",
	)
	s.Require().NotNil(undelAmount)
	return *undelAmount
}

func (s *TestSuite) GetChainValAddr(chain *starship.ChainClient) string {
	accAddr, err := chain.Client.DecodeBech32AccAddr(chain.Address)
	s.Require().NoError(err)
	valAddr, err := chain.Client.EncodeBech32ValAddr(sdk.ValAddress(accAddr))
	s.Require().NoError(err)
	return valAddr
}
