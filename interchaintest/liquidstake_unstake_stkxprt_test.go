package interchaintest

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	liquidstaketypes "github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v11/interchaintest/helpers"
)

// TestLiquidStakeUnstakeStkXPRT runs the flow of stkXPRT unstaking.
func TestLiquidStakeUnstakeStkXPRT(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// override SDK bech prefixes with chain specific
	helpers.SetConfig()

	ctx, cancelFn := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancelFn()
	})

	// create a single chain instance with 8 validators
	validatorsCount := 8

	// important overrides: fast voting for quick proposal passing
	ic, chain := CreateChain(t, ctx, validatorsCount, 0, fastVotingGenesisOverridesKV...)
	chainNode := chain.Nodes()[0]
	testDenom := chain.Config().Denom

	require.NotNil(t, ic)
	require.NotNil(t, chain)

	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Allocate user with funds
	chainUserFunds := int64(10_000_000_000)
	chainUser := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), chainUserFunds, chain)[0]

	// Get list of validators
	validators := helpers.QueryAllValidators(t, ctx, chainNode)
	require.Len(t, validators, validatorsCount, "validators returned must match count of validators created")

	// Updating liquidstake params for a new chain
	height, err := chain.Height(ctx)
	require.NoError(t, err, "error fetching height before submitting a proposal")

	whitelistedValidators := make([]liquidstaketypes.WhitelistedValidator, 0, len(validators))
	for _, val := range validators {
		whitelistedValidators = append(whitelistedValidators, liquidstaketypes.WhitelistedValidator{
			ValidatorAddress: val.OperatorAddress,
			TargetWeight:     math.NewInt(1),
		})
	}

	msgUpdateParams, err := codectypes.NewAnyWithValue(&liquidstaketypes.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress("gov").String(),
		Params: liquidstaketypes.Params{
			LiquidBondDenom:       liquidstaketypes.DefaultLiquidBondDenom,
			WhitelistedValidators: whitelistedValidators,
			LsmDisabled:           false,
			UnstakeFeeRate:        liquidstaketypes.DefaultUnstakeFeeRate,
			MinLiquidStakeAmount:  liquidstaketypes.DefaultMinLiquidStakeAmount,
		},
	})

	require.NoError(t, err, "failed to pack liquidstaketypes.MsgUpdateParams")

	broadcaster := cosmos.NewBroadcaster(t, chain)
	broadcaster.ConfigureFactoryOptions(func(factory tx.Factory) (_ tx.Factory) {
		return factory.WithGas(1_000_000)
	})
	txResp, err := cosmos.BroadcastTx(
		ctx,
		broadcaster,
		chainUser,
		&govv1.MsgSubmitProposal{
			InitialDeposit: []sdk.Coin{sdk.NewCoin(chain.Config().Denom, sdk.NewInt(500_000_000))},
			Proposer:       chainUser.FormattedAddress(),
			Title:          "LiquidStake Params Update",
			Summary:        "Sets whitelisted validators for liquidstake",
			Messages:       []*codectypes.Any{msgUpdateParams},
		},
	)
	require.NoError(t, err, "error submitting liquidstake params update tx")

	err = testutil.WaitForBlocks(ctx, 1, chain)
	require.NoError(t, err)

	upgradeTx, err := helpers.QueryProposalTx(context.Background(), chain.Nodes()[0], txResp.TxHash)
	require.NoError(t, err, "error checking proposal tx")

	err = chain.VoteOnProposalAllValidators(ctx, upgradeTx.ProposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit votes")

	_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+15, upgradeTx.ProposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	err = testutil.WaitForBlocks(ctx, 2, chain)
	require.NoError(t, err)

	// Liquid stake XPRT

	chainUserLiquidStakeAmount := sdk.NewInt(8_000_000)
	chainUserLiquidStakeCoins := sdk.NewCoin(testDenom, chainUserLiquidStakeAmount)
	txHash, err := chainNode.ExecTx(ctx, chainUser.KeyName(),
		"liquidstake", "liquid-stake", chainUserLiquidStakeCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err)

	stkXPRTBalance, err := chain.GetBalance(ctx, chainUser.FormattedAddress(), "stk/uxprt")
	require.NoError(t, err)
	require.Equal(t, chainUserLiquidStakeAmount, stkXPRTBalance, "stkXPRT balance must match the liquid-staked amount")

	// Try to unstake stkXPRT

	unstakeCoins := sdk.NewCoin("stk/uxprt", stkXPRTBalance)
	txHash, err = chainNode.ExecTx(ctx, chainUser.KeyName(),
		"liquidstake", "liquid-unstake", unstakeCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err)

	// Check token balances afterwards

	stkXPRTBalance, err = chain.GetBalance(ctx, chainUser.FormattedAddress(), "stk/uxprt")
	require.NoError(t, err)
	require.Equal(t, math.ZeroInt(), stkXPRTBalance, "user's stkXPRT balance must be 0")

	// Query the created unbonding delegations in favour of the user

	unbondingDelegations := helpers.QueryUnbondingDelegations(t, ctx, chainNode, chainUser.FormattedAddress())
	require.Len(t, unbondingDelegations, len(validators))
	for idx, ubd := range unbondingDelegations {
		require.Len(t, ubd.Entries, 1)
		require.Equal(t, validators[idx].OperatorAddress, ubd.ValidatorAddress, "unbonding delegation must match validator address at idx")
		require.Equal(t, chainUser.FormattedAddress(), ubd.DelegatorAddress, "unbonding delegation must have user as delegator")
		require.Equal(t, math.NewInt(1_000_000), ubd.Entries[0].Balance, "balance of unbonding delegation to match for stkXPRT unbonding piece")
	}
}
