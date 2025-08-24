package interchaintest

import (
	"context"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v13/interchaintest/helpers"
)

// TestMultiTokenizeVote case checks what happens with a voting tally when two
// separate LSM users are tokenizing and voting independently.
func TestMultiTokenizeVote(t *testing.T) {
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
	ic, chain := CreateChain(t, ctx, validatorsCount, 0, votingGenesisOverridesKV...)
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
	secondUserFunds := math.NewInt(1000)
	secondUser := interchaintest.GetAndFundTestUsers(t, ctx, secondUserName(t.Name()), secondUserFunds, chain)[0]

	// Get list of validators
	validators := helpers.QueryAllValidators(t, ctx, chainNode)
	require.Len(t, validators, validatorsCount, "validator returned must match count of validators created")

	// Bond first user
	firstUserBondAmount := sdk.NewInt(100000)
	firstUserBondCoins := sdk.NewCoin(testDenom, firstUserBondAmount)
	_, err := chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "delegate", validators[0].OperatorAddress, firstUserBondCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err)

	// Bond second user
	secondUserBondAmount := sdk.NewInt(1)
	secondUserBondCoins := sdk.NewCoin(testDenom, secondUserBondAmount)
	_, err = chainNode.ExecTx(ctx, secondUser.KeyName(),
		"staking", "delegate", validators[0].OperatorAddress, secondUserBondCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err)

	// Tokenize all shares - second user
	_, err = chainNode.ExecTx(ctx, secondUser.KeyName(),
		"staking", "tokenize-share", validators[0].OperatorAddress, secondUserBondCoins.String(), secondUser.FormattedAddress(),
		"--gas=500000",
	)
	require.NoError(t, err)

	sharesBalance, err := chain.GetBalance(ctx, secondUser.FormattedAddress(), validators[0].OperatorAddress+"/1")
	require.NoError(t, err)
	require.Equal(t, secondUserBondAmount, sharesBalance, "shares balance must match initially bonded amount")

	// Submitting a proposal from the second user
	height, err := chain.Height(ctx)
	require.NoError(t, err, "error fetching height before submit proposal")

	proposal := cosmos.TextProposal{
		Deposit:     "100" + testDenom,
		Title:       "Awesome proposal",
		Description: "Hello world!",
	}

	proposalTxHash, err := helpers.LegacyTextProposal(ctx, secondUser.KeyName(), chainNode, proposal)
	require.NoError(t, err, "error submitting text proposal tx")

	proposalTx, err := helpers.QueryProposalTx(ctx, chainNode, proposalTxHash)
	require.NoError(t, err, "error reading text proposal result")

	proposalID, err := strconv.ParseInt(proposalTx.ProposalID, 10, 64)
	require.NoError(t, err, "error parsing proposal id")

	_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+10, proposalID, govv1beta1.StatusVotingPeriod)
	require.NoError(t, err, "proposal status did not change to voting in expected number of blocks")

	// At this point, second user has 1stake tokenized, put up a proposal for 100stake initial deposit,
	// and now he votes once the proposal in the voting period.

	err = chainNode.VoteOnProposal(ctx, secondUser.KeyName(), proposalID, helpers.ProposalVoteYes)
	require.NoError(t, err)

	// The vote is not being reflected in the tally for now
	tally := helpers.QueryProposalTally(t, ctx, chainNode, proposalTx.ProposalID)
	require.Equal(t, math.ZeroInt(), tally.YesCount, "second user's tokenized shares don't show up in tally")

	// Redeem all shares - second user
	redeemCoints := sdk.NewCoin(validators[0].OperatorAddress+"/1", secondUserBondAmount)
	_, err = chainNode.ExecTx(ctx, secondUser.KeyName(),
		"staking", "redeem-tokens", redeemCoints.String(),
		"--gas=500000",
	)
	require.NoError(t, err)

	sharesBalance, err = chain.GetBalance(ctx, secondUser.FormattedAddress(), validators[0].OperatorAddress+"/1")
	require.NoError(t, err)
	require.Equal(t, math.ZeroInt(), sharesBalance, "second user's shares balance must be 0")

	// The vote will be reflected in the tally now (on behalf of second user's bond)
	tally = helpers.QueryProposalTally(t, ctx, chainNode, proposalTx.ProposalID)
	require.Equal(t, secondUserBondAmount, tally.YesCount, "second user's bonded amount counted towards Yes")

	// First user tries to vote with NoWithVeto using his own bond
	err = chainNode.VoteOnProposal(ctx, firstUser.KeyName(), proposalID, helpers.ProposalVoteNoWithVeto)
	require.NoError(t, err)

	// His vote is reflected in the tally (on behalf of the delegation)
	tally = helpers.QueryProposalTally(t, ctx, chainNode, proposalTx.ProposalID)
	require.Equal(t, firstUserBondAmount, tally.NoWithVetoCount, "first user's bonded amount counted towards NoWithVeto")

	// Tokenize all shares - first user
	_, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "tokenize-share", validators[0].OperatorAddress, firstUserBondCoins.String(), firstUser.FormattedAddress(),
		"--gas=500000",
	)
	require.NoError(t, err)

	sharesBalance, err = chain.GetBalance(ctx, firstUser.FormattedAddress(), validators[0].OperatorAddress+"/2")
	require.NoError(t, err)
	require.Equal(t, firstUserBondAmount, sharesBalance, "shares balance must match initially bonded amount")

	// Tokenize all shares - second user
	_, err = chainNode.ExecTx(ctx, secondUser.KeyName(),
		"staking", "tokenize-share", validators[0].OperatorAddress, secondUserBondCoins.String(), secondUser.FormattedAddress(),
		"--gas=500000",
	)
	require.NoError(t, err)

	sharesBalance, err = chain.GetBalance(ctx, secondUser.FormattedAddress(), validators[0].OperatorAddress+"/3")
	require.NoError(t, err)
	require.Equal(t, secondUserBondAmount, sharesBalance, "shares balance must match initially bonded amount")

	// No votes displayed in the voting tally as both delegators have tokenized their shares
	tally = helpers.QueryProposalTally(t, ctx, chainNode, proposalTx.ProposalID)
	require.Equal(t, math.ZeroInt(), tally.YesCount, "second user's tokenized shares amount doesn't show in the tally")
	require.Equal(t, math.ZeroInt(), tally.NoWithVetoCount, "first user's tokenized shares amount doesn't show in the tally")

	// Redeem all shares - first user
	redeemCoints = sdk.NewCoin(validators[0].OperatorAddress+"/2", firstUserBondAmount)
	_, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "redeem-tokens", redeemCoints.String(),
		"--gas=500000",
	)
	require.NoError(t, err)

	sharesBalance, err = chain.GetBalance(ctx, firstUser.FormattedAddress(), validators[0].OperatorAddress+"/2")
	require.NoError(t, err)
	require.Equal(t, math.ZeroInt(), sharesBalance, "second user's shares balance must be 0")

	// Check that first user's votes appear in the tally now
	tally = helpers.QueryProposalTally(t, ctx, chainNode, proposalTx.ProposalID)
	require.Equal(t, math.ZeroInt(), tally.YesCount, "second user's tokenized shares amount doesn't show in the tally")
	require.Equal(t, firstUserBondAmount, tally.NoWithVetoCount, "first user's bonded amount counted towards NoWithVeto")

	// Redeem all shares - second user
	redeemCoints = sdk.NewCoin(validators[0].OperatorAddress+"/3", secondUserBondAmount)
	_, err = chainNode.ExecTx(ctx, secondUser.KeyName(),
		"staking", "redeem-tokens", redeemCoints.String(),
		"--gas=500000",
	)
	require.NoError(t, err)

	sharesBalance, err = chain.GetBalance(ctx, secondUser.FormattedAddress(), validators[0].OperatorAddress+"/3")
	require.NoError(t, err)
	require.Equal(t, math.ZeroInt(), sharesBalance, "second user's shares balance must be 0")

	// Check that second user's votes appear in the tally now
	tally = helpers.QueryProposalTally(t, ctx, chainNode, proposalTx.ProposalID)
	require.Equal(t, secondUserBondAmount, tally.YesCount, "second user's bonded amount counted towards Yes")
	require.Equal(t, firstUserBondAmount, tally.NoWithVetoCount, "first user's bonded amount counted towards NoWithVeto")
}
