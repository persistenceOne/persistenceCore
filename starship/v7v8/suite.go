package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
	coretypes "github.com/tendermint/tendermint/rpc/core/types"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	starship "github.com/cosmology-tech/starship/clients/go/client"
)

var configFile = "./config.yaml"

type TestSuite struct {
	suite.Suite

	config       *starship.Config
	chainClients starship.ChainClients
}

func (s *TestSuite) SetupTest() {
	s.T().Log("setting up e2e integration test suite...")

	// read config file from yaml
	yamlFile, err := os.ReadFile(configFile)
	s.Require().NoError(err)
	config := &starship.Config{}
	err = yaml.Unmarshal(yamlFile, config)
	s.Require().NoError(err)
	s.config = config

	// create chain clients
	chainClients, err := starship.NewChainClients(zap.L(), config)
	s.Require().NoError(err)
	s.chainClients = chainClients
}

func (s *TestSuite) MakeRequest(req *http.Request, expCode int) io.Reader {
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err, "trying to make request", zap.Any("request", req))

	s.Require().Equal(expCode, resp.StatusCode, "response code did not match")

	return resp.Body
}

// WaitForTx will wait for the tx to complete, fail if not able to find tx
func (s *TestSuite) WaitForTx(chain *starship.ChainClient, txHex string) {
	var tx *coretypes.ResultTx
	var err error
	s.Require().Eventuallyf(
		func() bool {
			tx, err = chain.Client.QueryTx(context.Background(), txHex, false)
			if err != nil {
				return false
			}
			if tx.TxResult.Code == 0 {
				return true
			}
			return false
		},
		300*time.Second,
		time.Second,
		"waited for too long, still txn not successfull",
	)
	s.Assert().NotNil(tx)
}

// WaitForHeight will wait till the chain reaches the block height
func (s *TestSuite) WaitForHeight(chain *starship.ChainClient, height int64) {
	s.Require().Eventuallyf(
		func() bool {
			curHeight, err := chain.GetHeight()
			s.Assert().NoError(err)
			if curHeight >= height {
				return true
			}
			return false
		},
		300*time.Second,
		5*time.Second,
		"waited for too long, still height did not reach desired block height",
	)
}

func (s *TestSuite) TransferTokens(chain *starship.ChainClient, addr string, amount int, denom string) {
	coin, err := sdk.ParseCoinNormalized(fmt.Sprintf("%d%s", amount, denom))
	s.Require().NoError(err)

	// Build transaction message
	req := &banktypes.MsgSend{
		FromAddress: chain.Address,
		ToAddress:   addr,
		Amount:      sdk.Coins{coin},
	}

	res, err := chain.Client.SendMsg(context.Background(), req, "Transfer tokens for e2e tests")
	s.Require().NoError(err)

	s.WaitForTx(chain, res.TxHash)
}
