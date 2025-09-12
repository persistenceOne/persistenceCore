package interchaintest

import (
	"context"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"strconv"
	"testing"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/interchaintest/v10"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	"github.com/cosmos/interchaintest/v10/ibc"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v14/interchaintest/helpers"
)

// TestTokenizeSendVote checks that once shares are tokenized, the tokens can be
// sent to other party and used for voting, however, not counted in tally until bonded.
func TestTokenizeSendVote(t *testing.T) {
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
	firstUserBondAmount := math.NewInt(100000)
	firstUserBondCoins := sdk.NewCoin(testDenom, firstUserBondAmount)
	_, err := chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "delegate", validators[0].OperatorAddress, firstUserBondCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err)

	// Tokenize all shares - first user
	_, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "tokenize-share", validators[0].OperatorAddress, firstUserBondCoins.String(), firstUser.FormattedAddress(),
		"--gas=500000",
	)
	require.NoError(t, err)

	sharesBalance, err := chain.GetBalance(ctx, firstUser.FormattedAddress(), validators[0].OperatorAddress+"/1")
	require.NoError(t, err)
	require.Equal(t, firstUserBondAmount, sharesBalance, "shares balance must match initially bonded amount")

	// Send tokenized shares from first user to second user (only 1stake in tokenized denom)

	sharesToSend := ibc.WalletAmount{
		Address: secondUser.FormattedAddress(), // recipient
		Denom:   validators[0].OperatorAddress + "/1",
		Amount:  math.NewInt(1),
	}

	err = chainNode.SendFunds(ctx, firstUser.KeyName(), sharesToSend)
	require.NoError(t, err)

	sharesBalance, err = chain.GetBalance(ctx, secondUser.FormattedAddress(), validators[0].OperatorAddress+"/1")
	require.NoError(t, err)
	require.Equal(t, sharesToSend.Amount, sharesBalance, "second user's shares balance must match sent shares")

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

	proposalID, err := strconv.ParseUint(proposalTx.ProposalID, 10, 64)
	require.NoError(t, err, "error parsing proposal id")

	_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+10, proposalID, govv1beta1.StatusVotingPeriod)
	require.NoError(t, err, "proposal status did not change to voting in expected number of blocks")

	// At this point, second user has 1stake tokenized, put up a proposal for 100stake initial deposit,
	// and now he votes once the proposal in the voting period.

	err = chainNode.VoteOnProposal(ctx, secondUser.KeyName(), proposalID, helpers.ProposalVoteYes)
	require.NoError(t, err)

	// The vote is not being reflected in the tally for now
	tally := helpers.QueryProposalTally(t, ctx, chainNode, proposalTx.ProposalID)
	require.Equal(t, math.ZeroInt(), tally.YesCount, "second user's tokenized shares don't count in tally")

	// Redeem all shares - second user
	redeemCoints := sdk.NewCoin(validators[0].OperatorAddress+"/1", sharesBalance)
	_, err = chainNode.ExecTx(ctx, secondUser.KeyName(),
		"staking", "redeem-tokens", redeemCoints.String(),
		"--gas=500000",
	)
	require.NoError(t, err)

	sharesBalance, err = chain.GetBalance(ctx, secondUser.FormattedAddress(), validators[0].OperatorAddress+"/1")
	require.NoError(t, err)
	require.Equal(t, math.ZeroInt(), sharesBalance, "second user's shares balance must be 0")
	secondUserBondCoins := sdk.NewCoin(testDenom, sharesToSend.Amount)

	// The vote will be reflected in the tally (on behalf of second user - their shares were just bonded)
	tally = helpers.QueryProposalTally(t, ctx, chainNode, proposalTx.ProposalID)
	require.Equal(t, sharesToSend.Amount, tally.YesCount, "second user's bonded amount counted towards Yes")

	// Tokenize all shares - second user
	_, err = chainNode.ExecTx(ctx, secondUser.KeyName(),
		"staking", "tokenize-share", validators[0].OperatorAddress, secondUserBondCoins.String(), secondUser.FormattedAddress(),
		"--gas=500000",
	)
	require.NoError(t, err)

	sharesBalance, err = chain.GetBalance(ctx, secondUser.FormattedAddress(), validators[0].OperatorAddress+"/2")
	require.NoError(t, err)
	require.Equal(t, secondUserBondCoins.Amount, sharesBalance, "shares balance must match initially bonded amount")

	// Tokenized amount has been cleared up from the tally (second user now has liquid shares with his own denom)
	tally = helpers.QueryProposalTally(t, ctx, chainNode, proposalTx.ProposalID)
	require.Equal(t, math.ZeroInt(), tally.YesCount, "second user's bonded amount not shown in tally")

	// First user tries to vote with No and larger bond (still as liquid shares)
	err = chainNode.VoteOnProposal(ctx, firstUser.KeyName(), proposalID, helpers.ProposalVoteNoWithVeto)
	require.NoError(t, err)

	// His vote is not reflected in the tally (has no bond but liquid shares)
	tally = helpers.QueryProposalTally(t, ctx, chainNode, proposalTx.ProposalID)
	require.Equal(t, math.ZeroInt(), tally.YesCount, "second user's bonded amount not shown in tally")
	require.Equal(t, math.ZeroInt(), tally.NoWithVetoCount, "first user's bonded amount not counted towards NoWithVeto (since it's liquid)")

	// Redeem all shares - first user
	firstUserSharesLeftAmount := firstUserBondAmount.Sub(sharesToSend.Amount)
	redeemCoints = sdk.NewCoin(validators[0].OperatorAddress+"/1", firstUserSharesLeftAmount)
	_, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "redeem-tokens", redeemCoints.String(),
		"--gas=500000",
	)
	require.NoError(t, err)

	sharesBalance, err = chain.GetBalance(ctx, firstUser.FormattedAddress(), validators[0].OperatorAddress+"/1")
	require.NoError(t, err)
	require.Equal(t, math.ZeroInt(), sharesBalance, "first user's shares balance must be 0")

	// Votes from the first user now successfully reflected in the tally towards NoWithVeto
	tally = helpers.QueryProposalTally(t, ctx, chainNode, proposalTx.ProposalID)
	require.Equal(t, math.ZeroInt(), tally.YesCount, "second user's bonded amount not shown in tally")
	require.Equal(t, firstUserSharesLeftAmount, tally.NoWithVetoCount, "first user's bonded amount counted towards NoWithVeto now")

	// Redeem all shares - second user
	redeemCoints = sdk.NewCoin(validators[0].OperatorAddress+"/2", secondUserBondCoins.Amount)
	_, err = chainNode.ExecTx(ctx, secondUser.KeyName(),
		"staking", "redeem-tokens", redeemCoints.String(),
		"--gas=500000",
	)
	require.NoError(t, err)

	sharesBalance, err = chain.GetBalance(ctx, secondUser.FormattedAddress(), validators[0].OperatorAddress+"/2")
	require.NoError(t, err)
	require.Equal(t, math.ZeroInt(), sharesBalance, "second user's shares balance must be 0")

	// Votes from both users should now be reflected in the tally according to their vote options
	tally = helpers.QueryProposalTally(t, ctx, chainNode, proposalTx.ProposalID)
	require.Equal(t, secondUserBondCoins.Amount, tally.YesCount, "second user's bonded amount counted towards Yes now")
	require.Equal(t, firstUserSharesLeftAmount, tally.NoWithVetoCount, "first user's bonded amount counted towards NoWithVeto")
}
