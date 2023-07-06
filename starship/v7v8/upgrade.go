package main

import (
	"context"
	"time"

	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

func (s *TestSuite) Upgrade() {
	persistence := s.GetChainClient("test-core-1")

	currHeight, err := persistence.GetHeight()
	s.Require().NoError(err)

	upgradeName := "v8"
	upgradeHeight := currHeight + 50

	s.T().Logf("submitting v8 upgrade proposal, upgrade height: %d, current height: %d", upgradeHeight, currHeight)
	content := upgradetypes.NewSoftwareUpgradeProposal(
		"persistence v8 upgrade test",
		"persistence v8 upgrade test",
		upgradetypes.Plan{
			Name:   upgradeName,
			Height: upgradeHeight,
			Info:   "",
		},
	)
	proposalID := s.SubmitAndVoteProposal(persistence, content, "upgrade to v8")
	s.T().Logf("proposal submitted: %d", proposalID)

	expectedTimeToUpgradeHeight := time.Duration(upgradeHeight-currHeight-5) * time.Second // keeping margin for 5 blocks
	// sleeping here because WaitForHeight hits status api every second to check height
	// and gets this error after many repetitive calls
	// post failed: Post "http://localhost:26657": EOF
	s.T().Logf("Wating for %d seconds", expectedTimeToUpgradeHeight)
	time.Sleep(expectedTimeToUpgradeHeight)

	s.T().Log("waiting for upgrade height")
	s.WaitForHeight(persistence, upgradeHeight)

	s.T().Log("checking proposal status")
	res, err := govv1beta1.
		NewQueryClient(persistence.Client).
		Proposal(context.Background(), &govv1beta1.QueryProposalRequest{ProposalId: proposalID})
	s.Require().NoError(err)
	s.Require().Equal(govv1beta1.StatusPassed, res.Proposal.Status, "upgrade proposal did not pass before upgrade height: %d", upgradeHeight)

	s.T().Log("verifying upgrade happened")
	planRes, err := upgradetypes.
		NewQueryClient(persistence.Client).
		AppliedPlan(context.Background(), &upgradetypes.QueryAppliedPlanRequest{Name: upgradeName})
	s.Require().NoError(err)
	s.Require().Equal(upgradeHeight, planRes.Height)
	s.T().Log("upgrade successful")
}
