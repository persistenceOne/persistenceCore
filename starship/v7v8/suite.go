package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	coretypes "github.com/cometbft/cometbft/rpc/core/types"
	starship "github.com/cosmology-tech/starship/clients/go/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	lscosmos "github.com/persistenceOne/pstake-native/v2/x/lscosmos"
	lscosmostypes "github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
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
	chainClients, err := starship.NewChainClients(
		zap.L(),
		config,
		// extra codecs to register
		ibctm.AppModuleBasic{},
		lscosmos.AppModuleBasic{},
	)
	s.Require().NoError(err)
	s.chainClients = chainClients
}

func (s *TestSuite) GetChainClient(chainID string) *starship.ChainClient {
	chain, err := s.chainClients.GetChainClient(chainID)
	s.Require().NoError(err)
	return chain
}

func (s *TestSuite) SendMsgAndWait(chain *starship.ChainClient, msg sdk.Msg, memo string) *coretypes.ResultTx {
	res, err := chain.Client.SendMsg(context.Background(), msg, memo)
	s.Require().NoError(err)
	return s.WaitForTx(chain, res.TxHash)
}

// WaitForTx will wait for the tx to complete, fail if not able to find tx
func (s *TestSuite) WaitForTx(chain *starship.ChainClient, txHex string) *coretypes.ResultTx {
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
	return tx
}

// WaitForHeight will wait till the chain reaches the block height
func (s *TestSuite) WaitForHeight(chain *starship.ChainClient, height int64) {
	s.Require().Eventuallyf(
		func() bool {
			curHeight, err := chain.GetHeight()
			s.Assert().NoError(err)
			return curHeight >= height
		},
		300*time.Second,
		time.Second,
		"waited for too long, still height did not reach desired block height",
	)
}

func (s *TestSuite) WaitForNextBlock(chain *starship.ChainClient) {
	currHeight, err := chain.GetHeight()
	s.Require().NoError(err)
	s.WaitForHeight(chain, currHeight+1)
}

func (s *TestSuite) WaitForProposalToPass(chain *starship.ChainClient, proposalID uint64) {
	s.Require().Eventuallyf(
		func() bool {
			res, err := govv1beta1.
				NewQueryClient(chain.Client).
				Proposal(context.Background(), &govv1beta1.QueryProposalRequest{ProposalId: proposalID})
			s.Require().NoError(err)
			return res != nil && res.Proposal.Status == govv1beta1.StatusPassed
		},
		300*time.Second,
		time.Second,
		"waited for too long, proposal is still not passed",
	)
}

func (s *TestSuite) SubmitAndPassProposal(chain *starship.ChainClient, content govv1beta1.Content, memo string) {
	denom := chain.MustGetChainDenom()
	msg := &govv1beta1.MsgSubmitProposal{
		InitialDeposit: sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(10000000))),
		Proposer:       chain.Address,
	}
	err := msg.SetContent(&lscosmostypes.PstakeFeeAddressChangeProposal{
		Title:            "pstake fee address change proposal",
		Description:      "pstake fee address change proposal",
		PstakeFeeAddress: chain.Address,
	})
	s.T().Logf("submitting proposal: %s", memo)
	res := s.SendMsgAndWait(chain, msg, memo)

	id := s.FindEventAttr(res, "submit_proposal", "proposal_id")
	proposalID, err := strconv.Atoi(id)
	s.Require().NoError(err)

	s.T().Logf("submitting vote on proposal: %d | memo: %s", proposalID, memo)
	vote := &govv1beta1.MsgVote{ProposalId: uint64(proposalID), Voter: chain.Address, Option: govv1beta1.OptionYes}
	s.SendMsgAndWait(chain, vote, fmt.Sprintf("vote: %s", memo))
	s.WaitForProposalToPass(chain, uint64(proposalID))
}

func (s *TestSuite) FindEventAttr(res *coretypes.ResultTx, event, attr string) string {
	for _, txEvent := range res.TxResult.Events {
		if txEvent.Type == event {
			for _, txAttr := range txEvent.Attributes {
				key, err := base64.StdEncoding.DecodeString(txAttr.Key)
				s.Require().NoError(err)
				if string(key) == attr {
					val, err := base64.StdEncoding.DecodeString(txAttr.Value)
					s.Require().NoError(err)
					return string(val)
				}
			}
		}
	}
	s.FailNow("event attr not found in tx events")
	return ""
}
