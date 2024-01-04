package main

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
)

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestChainsStatus() {
	s.T().Log("running test for /status endpoint for each chain")

	for _, chainClient := range s.chainClients {
		status, err := chainClient.GetStatus()
		s.Assert().NoError(err)

		s.Assert().Equal(chainClient.ChainID, status.NodeInfo.Network)
	}
}

func (s *TestSuite) TestRegression() {
	s.RunTokenTransferTests()
	s.RunIBCTokenTransferTests()
}

func (s *TestSuite) TestUpgrade() {
	if testing.Short() {
		s.T().Skip("Skipping chain upgrade tests for short test")
	}

	/* Pre upgrade tests */
	s.RunLSCosmosTests()

	/* Upgrade */
	s.Upgrade()

	/* Post upgrade tests */
	s.VerifyParams()
	s.VerifyValidatorCommissionRates()
	s.RunLiquidstakeibcTests()
}

func (s *TestSuite) VerifyValidatorCommissionRates() {
	minRate := sdk.NewDecWithPrec(5, 2)
	minMaxRate := sdk.NewDecWithPrec(1, 1)
	client := s.GetChainClient("test-core-2").Client
	vals, err := stakingtypes.NewQueryClient(client).Validators(context.Background(), &stakingtypes.QueryValidatorsRequest{})
	s.Require().NoError(err)
	for _, val := range vals.Validators {
		s.Require().True(val.Commission.Rate.GTE(minRate))
		s.Require().True(val.Commission.MaxRate.GTE(minMaxRate))
	}
}
