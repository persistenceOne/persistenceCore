package main

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

		s.Assert().Equal(chainClient.GetChainID(), status.NodeInfo.Network)
	}
}

func (s *TestSuite) TestChainTokenTransfer() {
	chain1, err := s.chainClients.GetChainClient("core-1")
	s.Require().NoError(err)

	keyName := "test-transfer"
	address, err := chain1.CreateRandWallet(keyName)
	s.Require().NoError(err)

	denom, err := chain1.GetChainDenom()
	s.Require().NoError(err)

	s.TransferTokens(chain1, address, 2345000, denom)

	// Verify the address recived the token
	balance, err := chain1.Client.QueryBalanceWithDenomTraces(context.Background(), sdk.MustAccAddressFromBech32(address), nil)
	s.Require().NoError(err)

	// Assert correct transfers
	s.Assert().Len(balance, 1)
	s.Assert().Equal(balance.Denoms(), []string{denom})
	s.Assert().Equal(balance[0].Amount, sdk.NewInt(2345000))
}
