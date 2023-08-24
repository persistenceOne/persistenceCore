package interchaintest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceCore/v8/interchaintest/helpers"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
)

// TestBondTokenize executes scenario of bonding and tokenizing.
func TestBondTokenize(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	ctx, cancelFn := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancelFn()
	})

	// create a single chain instance with 4 validators
	validatorsCount := 4
	ic, chain := CreateChain(t, ctx, validatorsCount, 0, cosmos.GenesisKV{
		Key:   "app_state.staking.params.validator_bond_factor",
		Value: "250",
	})
	chainNode := chain.Nodes()[0]
	testDenom := chain.Config().Denom

	require.NotNil(t, ic)
	require.NotNil(t, chain)

	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Allocate two chain users with funds
	firstUserFunds := int64(10_000_000_000)
	firstUser := interchaintest.GetAndFundTestUsers(t, ctx, firstUserName(t.Name()), firstUserFunds, chain)[0]
	secondUserFunds := int64(2_000_000)
	secondUser := interchaintest.GetAndFundTestUsers(t, ctx, secondUserName(t.Name()), secondUserFunds, chain)[0]

	// Get list of validators
	validators := helpers.QueryAllValidators(t, ctx, chainNode)
	require.Len(t, validators, validatorsCount, "validator returned must match count of validators created")

	// Delegate from first user
	firstUserDelegationAmount := sdk.NewInt(1_000_000_000)
	firstUserDelegationCoins := sdk.NewCoin(testDenom, firstUserDelegationAmount)
	_, err := chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "delegate", validators[0].OperatorAddress, firstUserDelegationCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err)

	// Delegate from second user
	secondUserDelegationAmount := sdk.NewInt(1_000_000)
	secondUserDelegationCoins := sdk.NewCoin(testDenom, secondUserDelegationAmount)
	_, err = chainNode.ExecTx(ctx, secondUser.KeyName(),
		"staking", "delegate", validators[0].OperatorAddress, secondUserDelegationCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err)

	delegation := helpers.QueryDelegation(t, ctx, chainNode, secondUser.FormattedAddress(), validators[0].OperatorAddress)
	require.Equal(t, sdk.NewDecFromInt(secondUserDelegationCoins.Amount), delegation.Shares)
	require.False(t, delegation.ValidatorBond)

	// Try to tokenize shares from first user, it won't work because there is no minimal bond
	tokenizeCoins := sdk.NewCoin(testDenom, sdk.NewInt(250_000_000))
	txHash, err := chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "tokenize-share", validators[0].OperatorAddress, tokenizeCoins.String(), firstUser.FormattedAddress(),
		"--gas=500000",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.ErrorContains(t, err, "insufficient validator bond shares")

	// Mark second user bond as validator bond
	_, err = chainNode.ExecTx(ctx, secondUser.KeyName(),
		"staking", "validator-bond", validators[0].OperatorAddress,
		"--gas=500000",
	)
	require.NoError(t, err)

	delegation = helpers.QueryDelegation(t, ctx, chainNode, secondUser.FormattedAddress(), validators[0].OperatorAddress)
	require.Equal(t, sdk.NewDecFromInt(secondUserDelegationCoins.Amount), delegation.Shares)
	require.True(t, delegation.ValidatorBond)

	validator := helpers.QueryValidator(t, ctx, chainNode, validators[0].OperatorAddress)
	require.Equal(t,
		secondUserDelegationAmount.Int64(),
		validator.ValidatorBondShares.TruncateInt().Int64(),
		"validator bond shares must match bonded amount",
	)

	// Try to tokenize shares from first user again, it should work now
	txHash, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "tokenize-share", validators[0].OperatorAddress, tokenizeCoins.String(), firstUser.FormattedAddress(),
		"--gas=500000",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err)

	sharesBalance, err := chain.GetBalance(ctx, firstUser.FormattedAddress(), validators[0].OperatorAddress+"/1")
	require.NoError(t, err)
	require.Equal(t, tokenizeCoins.Amount.Int64(), sharesBalance, "shares balance must match tokenized amount")

	// Try to tokenize more shares from first user, it will not work because of small bond
	txHash, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "tokenize-share", validators[0].OperatorAddress, tokenizeCoins.String(), firstUser.FormattedAddress(),
		"--gas=500000",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.ErrorContains(t, err, "insufficient validator bond shares")

	// Delegate from second user more
	txHash, err = chainNode.ExecTx(ctx, secondUser.KeyName(),
		"staking", "delegate", validators[0].OperatorAddress, secondUserDelegationCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err)

	delegation = helpers.QueryDelegation(t, ctx, chainNode, secondUser.FormattedAddress(), validators[0].OperatorAddress)
	secondUserDelegationCoinsDouble := sdk.NewDecFromInt(secondUserDelegationCoins.Amount).MulInt64(2)
	require.Equal(t, secondUserDelegationCoinsDouble, delegation.Shares, "expected updated delegation")
	require.True(t, delegation.ValidatorBond)

	// Try to tokenize more shares from first user, it must work now
	txHash, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "tokenize-share", validators[0].OperatorAddress, tokenizeCoins.String(), firstUser.FormattedAddress(),
		"--gas=500000",
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err)

	sharesBalance, err = chain.GetBalance(ctx, firstUser.FormattedAddress(), validators[0].OperatorAddress+"/2")
	require.NoError(t, err)
	require.Equal(t, tokenizeCoins.Amount.Int64(), sharesBalance, "shares balance must match tokenized amount")

	validator = helpers.QueryValidator(t, ctx, chainNode, validators[0].OperatorAddress)
	doubleTokenizedAmount := sdk.NewDecFromInt(tokenizeCoins.Amount.MulRaw(2))
	require.Equal(t, doubleTokenizedAmount, validator.LiquidShares, "validator's liquid shares amount must match tokenized amount x2")
}
