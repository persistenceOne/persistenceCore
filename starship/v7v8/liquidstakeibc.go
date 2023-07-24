package main

import (
	"context"

	starship "github.com/cosmology-tech/starship/clients/go/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	lsibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	lscosmostypes "github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
)

func (s *TestSuite) RunLiquidstakeibcTests() {
	persistence := s.GetChainClient("test-core-2")

	s.VerifyMigratedState(persistence)
}

func (s *TestSuite) VerifyMigratedState(persistence *starship.ChainClient) {
	s.T().Log("verifying lscosmos migrated state")
	gaia := s.GetChainClient("test-gaia-1")
	uatom := gaia.MustGetChainDenom()

	lscosmosState, err := lscosmostypes.
		NewQueryClient(persistence.Client).
		AllState(context.Background(), &lscosmostypes.QueryAllStateRequest{})
	s.Require().NoError(err)

	lsibcState, err := lsibctypes.
		NewQueryClient(persistence.Client).
		HostChains(context.Background(), &lsibctypes.QueryHostChainsRequest{})
	s.Require().NoError(err)
	s.Require().GreaterOrEqual(1, len(lsibcState.HostChains))
	s.Require().Equal(&lsibctypes.QueryHostChainsResponse{
		HostChains: []*lsibctypes.HostChain{{
			ChainId:      lscosmosState.Genesis.HostChainParams.ChainID,
			ConnectionId: lscosmosState.Genesis.HostChainParams.ConnectionID,
			ChannelId:    lscosmosState.Genesis.HostChainParams.TransferChannel,
			PortId:       lscosmosState.Genesis.HostChainParams.TransferPort,
			HostDenom:    lscosmosState.Genesis.HostChainParams.BaseDenom,
			Params: &lsibctypes.HostChainLSParams{
				DepositFee:    lscosmosState.Genesis.HostChainParams.PstakeParams.PstakeDepositFee,
				RestakeFee:    lscosmosState.Genesis.HostChainParams.PstakeParams.PstakeRestakeFee,
				UnstakeFee:    lscosmosState.Genesis.HostChainParams.PstakeParams.PstakeUnstakeFee,
				RedemptionFee: lscosmosState.Genesis.HostChainParams.PstakeParams.PstakeRedemptionFee,
			},
			DelegationAccount: &lsibctypes.ICAAccount{
				Address:      lscosmosState.Genesis.DelegationState.HostChainDelegationAddress,
				Owner:        lscosmosState.Genesis.HostAccounts.DelegatorAccountOwnerID,
				ChannelState: lsibctypes.ICAAccount_ICA_CHANNEL_CREATED,
				Balance:      sdk.NewInt64Coin(uatom, 0),
			},
			RewardsAccount: &lsibctypes.ICAAccount{
				Address:      lscosmosState.Genesis.HostChainRewardAddress.Address,
				Owner:        lscosmosState.Genesis.HostAccounts.RewardsAccountOwnerID,
				ChannelState: lsibctypes.ICAAccount_ICA_CHANNEL_CREATED,
				Balance:      sdk.NewInt64Coin(uatom, 0),
			},
			Validators: []*lsibctypes.Validator{{
				OperatorAddress: lscosmosState.Genesis.AllowListedValidators.AllowListedValidators[0].ValidatorAddress,
				Status:          "BOND_STATUS_BONDED",
				Weight:          lscosmosState.Genesis.AllowListedValidators.AllowListedValidators[0].TargetWeight,
				DelegatedAmount: lscosmosState.Genesis.DelegationState.HostAccountDelegations[0].Amount.Amount,
				ExchangeRate:    sdk.NewDec(1),
				UnbondingEpoch:  0,
			}},
			MinimumDeposit:     lsibcState.HostChains[0].MinimumDeposit,
			CValue:             lsibcState.HostChains[0].CValue,
			LastCValue:         lsibcState.HostChains[0].LastCValue,
			UnbondingFactor:    lsibcState.HostChains[0].UnbondingFactor,
			Active:             false,
			AutoCompoundFactor: lsibcState.HostChains[0].AutoCompoundFactor,
		}},
	}, lsibcState)

	lsibcParams, err := lsibctypes.
		NewQueryClient(persistence.Client).
		Params(context.Background(), &lsibctypes.QueryParamsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(lscosmosState.Genesis.HostChainParams.PstakeParams.PstakeFeeAddress, lsibcParams.Params.FeeAddress)
}
