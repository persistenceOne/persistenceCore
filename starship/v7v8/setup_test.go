package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (s *TestSuite) TestChainsStatus() {
	s.T().Log("runing test for /status endpoint for each chain")

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
	// s.Upgrade()

	/* Post upgrade tests */
	// s.PostUpgradeTest1()
	// s.PostUpgradeTest2()
}
