package interchaintest

import (
	"context"
	"encoding/json"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/interchaintest/v10"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	"github.com/cosmos/interchaintest/v10/ibc"
	"github.com/cosmos/interchaintest/v10/testutil"
	liquidstaketypes "github.com/persistenceOne/pstake-native/v4/x/liquidstake/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v14/interchaintest/helpers"
)

// TestPauseLiquidStakeStkXPRT runs the flow of liquid XPRT staking while pausing the module.
func TestPauseLiquidStakeStkXPRT(t *testing.T) {
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
	overridesKV = append(overridesKV, cosmos.GenesisKV{
		Key:   "app_state.liquidstake.params.module_paused",
		Value: false,
	})
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

	// Pause the module

	txHash, err := chainNode.ExecTx(ctx, firstUser.KeyName(),
		"liquidstake", "pause-module", "true",
		"--gas=auto",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err)

	// Liquid stake XPRT from the first user (5 XPRT)

	firstUserLiquidStakeAmount := math.NewInt(5_000_000)
	firstUserLiquidStakeCoins := sdk.NewCoin(testDenom, firstUserLiquidStakeAmount)
	_, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"liquidstake", "liquid-stake", firstUserLiquidStakeCoins.String(),
		"--gas=auto",
	)
	// uh-oh!
	require.ErrorContains(t, err, "module functions have been paused")

	// Unpause the module

	txHash, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"liquidstake", "pause-module", "false",
		"--gas=auto",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err)

	txHash, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"liquidstake", "liquid-stake", firstUserLiquidStakeCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err)

	stkXPRTBalance, err := chain.GetBalance(ctx, firstUser.FormattedAddress(), "stk/uxprt")
	require.NoError(t, err)
	require.Equal(t, firstUserLiquidStakeAmount, stkXPRTBalance, "stkXPRT balance must match the liquid-staked amount")

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
	require.Equal(t, xprtBalance.Int64(), secondUserFunds, "second user's XPRT balance must be untouched")

	// Query the created unbonding delegation in favour of second user

	unbondingDelegation := helpers.QueryUnbondingDelegation(t, ctx, chainNode, secondUser.FormattedAddress(), validators[0].OperatorAddress)
	require.Len(t, unbondingDelegation.Entries, 1)
	require.Equal(t, secondUser.FormattedAddress(), unbondingDelegation.DelegatorAddress, "unbonding delegation must have second user as delegator")
	expectedUnbondingBalance := tokensToSend.Amount
	require.Equal(t, expectedUnbondingBalance, unbondingDelegation.Entries[0].Balance, "balance of unbonding delegation to match for stkXPRT unbonding")
}
