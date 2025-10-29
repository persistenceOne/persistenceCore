package interchaintest

import (
	"context"
	"cosmossdk.io/math"
	"encoding/json"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/interchaintest/v10"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	"github.com/cosmos/interchaintest/v10/ibc"
	"github.com/cosmos/interchaintest/v10/testutil"
	"github.com/persistenceOne/persistenceCore/v16/interchaintest/helpers"
	liquidstaketypes "github.com/persistenceOne/pstake-native/v5/x/liquidstake/types"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

// TestLiquidStakeStkXPRT runs the flow of liquid XPRT staking using
// liquidstake module, including LSM-LP flow when stake gets locked into Superfluid LP.
func TestLiquidStakeStkXPRT(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	// t.Parallel()

	// override SDK bech prefixes with chain specific
	helpers.SetConfig()

	ctx, cancelFn := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancelFn()
	})

	// create a single chain instance with 4 validators
	validatorsCount := 4

	overridesKV := append([]cosmos.GenesisKV{}, fastVotingGenesisOverridesKV...)
	// By default the module is in paused state.
	// set it unpaused.
	overridesKV = append(overridesKV, cosmos.GenesisKV{
		Key:   "app_state.liquidstake.params.module_paused",
		Value: false,
	})

	// important overrides: fast voting for quick proposal passing
	ic, chain := CreateChain(t, ctx, validatorsCount, 0, overridesKV...)
	chainNode := chain.Nodes()[0]
	testDenom := chain.Config().Denom

	require.NotNil(t, ic)
	require.NotNil(t, chain)

	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Allocate two chain users with funds
	firstUserFunds := math.NewInt(10_000_000_000)
	firstUser := interchaintest.GetAndFundTestUsers(t, ctx, firstUserName(t.Name()), firstUserFunds, chain)[0]
	secondUserFunds := math.NewInt(1_000_000)
	secondUser := interchaintest.GetAndFundTestUsers(t, ctx, secondUserName(t.Name()), secondUserFunds, chain)[0]

	instantiateMsg, err := json.Marshal(helpers.SuperFluidInstantiateMsg{
		VaultAddress: "persistence1z0uz82yle9tavl4qpq86a34z4hn7gsdd8n56t3qzr0nf4nwptv8q3h274d",
		Owner:        firstUser.FormattedAddress(),
		AllowedLockableTokens: []helpers.AssetInfo{{
			NativeToken: helpers.NativeTokenInfo{
				Denom: "stk/uxprt",
			},
		}},
	})
	require.NoError(t, err)

	_, lpContractAddr := helpers.SetupContract(
		t, ctx, chain, firstUser.KeyName(),
		"contracts/dexter_superfluid_lp.wasm",
		string(instantiateMsg),
	)

	t.Logf("Deployed Superfluid LP contract: %s", lpContractAddr)

	lockedLST := helpers.GetTotalAmountLocked(t, ctx, chainNode, lpContractAddr, firstUser.FormattedAddress())
	require.Equal(t, math.ZeroInt(), lockedLST, "no locked LST expected")

	// Get list of validators
	validators := helpers.QueryAllValidators(t, ctx, chainNode)
	require.Len(t, validators, validatorsCount, "validators returned must match count of validators created")

	// Updating liquidstake params for a new chain
	height, err := chain.Height(ctx)
	require.NoError(t, err, "error fetching height before submitting a proposal")

	msgUpdateParams, err := codectypes.NewAnyWithValue(&liquidstaketypes.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress("gov").String(),
		Params: liquidstaketypes.Params{
			LiquidBondDenom:       liquidstaketypes.DefaultLiquidBondDenom,
			LsmDisabled:           false,
			UnstakeFeeRate:        liquidstaketypes.DefaultUnstakeFeeRate,
			MinLiquidStakeAmount:  liquidstaketypes.DefaultMinLiquidStakeAmount,
			CwLockedPoolAddress:   lpContractAddr,
			FeeAccountAddress:     liquidstaketypes.DummyFeeAccountAcc.String(),
			AutocompoundFeeRate:   liquidstaketypes.DefaultAutocompoundFeeRate,
			WhitelistAdminAddress: firstUser.FormattedAddress(),
			ModulePaused:          false,
		},
	})

	require.NoError(t, err, "failed to pack liquidstaketypes.MsgUpdateParams")

	broadcaster := cosmos.NewBroadcaster(t, chain)
	txResp, err := cosmos.BroadcastTx(
		ctx,
		broadcaster,
		firstUser,
		&govv1.MsgSubmitProposal{
			InitialDeposit: []sdk.Coin{sdk.NewCoin(chain.Config().Denom, math.NewInt(500_000_000))},
			Proposer:       firstUser.FormattedAddress(),
			Title:          "LiquidStake Params Update",
			Summary:        "Sets params for liquidstake",
			Messages:       []*codectypes.Any{msgUpdateParams},
		},
	)
	require.NoError(t, err, "error submitting liquidstake params update tx")

	upgradeTx, err := helpers.QueryProposalTx(context.Background(), chain.Nodes()[0], txResp.TxHash)
	require.NoError(t, err, "error checking proposal tx")

	proposalID, err := strconv.ParseUint(upgradeTx.ProposalID, 10, 64)
	require.NoError(t, err, "error parsing proposal id")

	proposal, err := chain.GovQueryProposalV1(ctx, proposalID)
	require.NoError(t, err, "error getting proposal")
	t.Log(proposal)
	require.Equal(t, govv1.StatusVotingPeriod, proposal.Status, "proposal status equal check")

	err = chain.VoteOnProposalAllValidators(ctx, proposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit votes")

	_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+15, proposalID, govv1beta1.StatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	err = testutil.WaitForBlocks(ctx, 2, chain)
	require.NoError(t, err)

	// Update whitelisted validators list from the first user (just for convenience)
	txResp, err = cosmos.BroadcastTx(
		ctx,
		broadcaster,
		firstUser,
		&liquidstaketypes.MsgUpdateWhitelistedValidators{
			Authority: firstUser.FormattedAddress(),
			WhitelistedValidators: []liquidstaketypes.WhitelistedValidator{{
				ValidatorAddress: validators[0].OperatorAddress,
				TargetWeight:     math.NewInt(10000),
			}},
		},
	)
	require.NoError(t, err, "error submitting liquidstake validators whitelist update tx")
	require.Equal(t, uint32(0), txResp.Code, txResp.RawLog)
	// Liquid stake XPRT from the first user (5 XPRT)

	firstUserLiquidStakeAmount := math.NewInt(5_000_000)
	firstUserLiquidStakeCoins := sdk.NewCoin(testDenom, firstUserLiquidStakeAmount)
	txHash, err := chainNode.ExecTx(ctx, firstUser.KeyName(),
		"liquidstake", "liquid-stake", firstUserLiquidStakeCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err)

	stkXPRTBalance, err := chain.GetBalance(ctx, firstUser.FormattedAddress(), "stk/uxprt")
	require.NoError(t, err)
	require.Equal(t, firstUserLiquidStakeAmount, stkXPRTBalance, "stkXPRT balance must match the liquid-staked amount")

	// Lock some liquid stkXPRT tokens into LP contract manually, using a direct CW call

	tokensToLock := sdk.NewCoin("stk/uxprt", math.NewInt(1_000_000))

	msg := &helpers.LockLstAssetMsg{
		Asset: helpers.Asset{
			Amount: tokensToLock.Amount,
			Info: helpers.AssetInfo{
				NativeToken: helpers.NativeTokenInfo{
					Denom: tokensToLock.Denom,
				},
			},
		},
	}

	callData, err := json.Marshal(&helpers.ExecMsg{
		LockLstAsset: msg,
	})
	require.NoError(t, err, "failed to marshal ExecMsg")

	txHash = helpers.ExecuteMsgWithAmount(t, ctx, chain, firstUser, lpContractAddr, tokensToLock.String(), string(callData))
	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err)

	lockedLST = helpers.GetTotalAmountLocked(t, ctx, chainNode, lpContractAddr, firstUser.FormattedAddress())
	require.Equal(t, tokensToLock.Amount, lockedLST, "expected LST tokens to be locked")

	stkXPRTBalance, err = chain.GetBalance(ctx, firstUser.FormattedAddress(), "stk/uxprt")
	require.NoError(t, err)
	require.Equal(t, firstUserLiquidStakeAmount.Sub(tokensToLock.Amount), stkXPRTBalance, "first user's stkXPRT balance must be reduced by locked stkXPRT")

	// Delegate from the first user to get a delegation that could be used to obtain non-liquid stkXPRT

	firstUserDelegationAmount := math.NewInt(5_000_000)
	firstUserDelegationCoins := sdk.NewCoin(testDenom, firstUserDelegationAmount)

	txHash, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "delegate", validators[0].OperatorAddress, firstUserDelegationCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err)

	delegation := helpers.QueryDelegation(t, ctx, chainNode, firstUser.FormattedAddress(), validators[0].OperatorAddress)
	require.Equal(t, math.LegacyNewDecFromInt(firstUserDelegationCoins.Amount), delegation.Shares)
	require.False(t, delegation.ValidatorBond)

	firstUserXPRTBalanceBeforeLock, err := chain.GetBalance(ctx, firstUser.FormattedAddress(), testDenom)
	require.NoError(t, err)

	// Lock more liquid stkXPRT tokens into LP contract, as well as stake (through implicit LSM)
	// using pStake's liquidstake module
	tokensToLock2 := sdk.NewCoin(testDenom, math.NewInt(1_000_000))
	stakeToLP := sdk.NewCoin(testDenom, math.NewInt(2_000_000))

	txHash, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"liquidstake", "stake-to-lp", validators[0].OperatorAddress, stakeToLP.String(), tokensToLock2.String(),
		"--gas=auto",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err)

	// check that delegation has reduced by stakeToLP amount
	delegation = helpers.QueryDelegation(t, ctx, chainNode, firstUser.FormattedAddress(), validators[0].OperatorAddress)
	require.Equal(t, firstUserDelegationCoins.Amount.Sub(stakeToLP.Amount).ToLegacyDec(), delegation.Shares)

	firstUserXPRTBalanceAfterLock, err := chain.GetBalance(ctx, firstUser.FormattedAddress(), testDenom)
	require.NoError(t, err)
	require.Equal(t,
		firstUserXPRTBalanceBeforeLock.Sub(tokensToLock2.Amount).Add(math.NewInt(1)), // fix a blip from LSM rewards
		firstUserXPRTBalanceAfterLock,
		"first user's XPRT balance must be reduced by locked XPRT during stake-to-lp",
	)

	// Check total expected locked stkXPRT in LP: two deposits of liquid stkXPRT in different ways
	// and one stake transfer through LSM-LP flow (using stake-to-lp).
	totalLockedExpected := tokensToLock.Amount.Add(tokensToLock2.Amount).Add(stakeToLP.Amount)
	// totalLockedExpected = totalLockedExpected.Sub(math.NewInt(1)) // some dust lost due to stk math

	lockedLST = helpers.GetTotalAmountLocked(t, ctx, chainNode, lpContractAddr, firstUser.FormattedAddress())
	require.Equal(t, totalLockedExpected, lockedLST, "expected LST tokens to add up")

	// Send some stkXPRT tokens from first user to second user

	tokensToSend := ibc.WalletAmount{
		Address: secondUser.FormattedAddress(), // recipient
		Denom:   "stk/uxprt",
		Amount:  math.NewInt(1_000_000),
	}

	err = chainNode.SendFunds(ctx, firstUser.KeyName(), tokensToSend)
	require.NoError(t, err)

	stkXPRTBalance, err = chain.GetBalance(ctx, secondUser.FormattedAddress(), "stk/uxprt")
	require.NoError(t, err)
	require.Equal(t, tokensToSend.Amount, stkXPRTBalance, "second user's stkXPRT balance must match sent stk tokens")

	// Try to unstake stkXPRT from second user

	unstakeCoins := sdk.NewCoin("stk/uxprt", stkXPRTBalance)
	txHash, err = chainNode.ExecTx(ctx, secondUser.KeyName(),
		"liquidstake", "liquid-unstake", unstakeCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err)

	// Check token balances afterwards

	stkXPRTBalance, err = chain.GetBalance(ctx, secondUser.FormattedAddress(), "stk/uxprt")
	require.NoError(t, err)
	require.Equal(t, math.ZeroInt(), stkXPRTBalance, "second user's stkXPRT balance must be 0")

	xprtBalance, err := chain.GetBalance(ctx, secondUser.FormattedAddress(), "uxprt")
	require.NoError(t, err)
	require.Equal(t, xprtBalance.Int64(), secondUserFunds.Int64(), "second user's XPRT balance must be untouched")

	// Query the created unbonding delegation in favour of second user

	unbondingDelegation := helpers.QueryUnbondingDelegation(t, ctx, chainNode, secondUser.FormattedAddress(), validators[0].OperatorAddress)
	require.Len(t, unbondingDelegation.Entries, 1)
	require.Equal(t, secondUser.FormattedAddress(), unbondingDelegation.DelegatorAddress, "unbonding delegation must have second user as delegator")
	expectedUnbondingBalance := tokensToSend.Amount
	require.Equal(t, expectedUnbondingBalance, unbondingDelegation.Entries[0].Balance, "balance of unbonding delegation to match for stkXPRT unbonding")
}
