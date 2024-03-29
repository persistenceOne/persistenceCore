package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	starship "github.com/cosmology-tech/starship/clients/go/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibcchanneltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
)

func (s *TestSuite) RunIBCTokenTransferTests() {
	persistence := s.GetChainClient("test-core-2")
	gaia := s.GetChainClient("test-gaia-1")
	uxprt := persistence.MustGetChainDenom()
	amt := 10000000    // 10 XPRTs
	amtBack := 1000000 // 1 XPRT

	c := s.GetTransferChannel(gaia, persistence.ChainID)
	port, channel := c.PortId, c.ChannelId

	// ibc atom denom would be empty if there are no traces of this uxprt yet
	ibcXPRT := s.GetIBCDenom(gaia, port, channel, uxprt)
	balBefore := s.GetBalance(gaia, gaia.Address, ibcXPRT)

	s.T().Logf("transferring %d%s to gaia address", amt, uxprt)
	seq := s.IBCTransferTokens(persistence, gaia.Address, port, channel, uxprt, amt)
	s.WaitForIBCPacketAck(gaia, port, channel, seq)

	ibcXPRT = s.GetIBCDenom(gaia, port, channel, uxprt)
	balAfter := s.GetBalance(gaia, gaia.Address, ibcXPRT)

	s.T().Log("check balance after ibc transfer")
	s.Require().Equal(balBefore.Amount.Add(sdk.NewInt(int64(amt))), balAfter.Amount)

	s.T().Logf("transferring %d%s back to persistence address (in ibc denoms)", amtBack, uxprt)
	port, channel = c.Counterparty.PortId, c.Counterparty.ChannelId
	balBefore = s.GetBalance(persistence, persistence.Address, uxprt)
	seq = s.IBCTransferTokens(gaia, persistence.Address, port, channel, ibcXPRT, amtBack)
	s.WaitForIBCPacketAck(persistence, port, channel, seq)

	balAfter = s.GetBalance(persistence, persistence.Address, uxprt)
	s.T().Log("check balance after transferring back ibc token")
	s.Require().Equal(balBefore.AddAmount(sdk.NewInt(int64(amtBack))), balAfter)
}

func (s *TestSuite) GetIBCDenom(chain *starship.ChainClient, port, channel, denom string) string {
	res, err := ibctransfertypes.
		NewQueryClient(chain.Client).
		DenomHash(context.Background(), &ibctransfertypes.QueryDenomHashRequest{
			Trace: fmt.Sprintf("%s/%s/%s", port, channel, denom),
		})
	if err != nil && strings.Contains(err.Error(), ibctransfertypes.ErrTraceNotFound.Error()) {
		return ""
	}
	s.Require().NoError(err)
	return fmt.Sprintf("ibc/%s", res.Hash)
}

func (s *TestSuite) IBCTransferTokens(chain *starship.ChainClient, receiver, port, channel, denom string, amount int) string {
	coin, err := sdk.ParseCoinNormalized(fmt.Sprintf("%d%s", amount, denom))
	s.Require().NoError(err)

	lh := s.GetIBCClientLatestHeight(chain)
	msg := &ibctransfertypes.MsgTransfer{
		SourcePort:    port,
		SourceChannel: channel,
		Token:         coin,
		Sender:        chain.Address,
		Receiver:      receiver,
		TimeoutHeight: ibcclienttypes.Height{
			RevisionNumber: lh.GetRevisionNumber(),
			RevisionHeight: lh.GetRevisionHeight() + 100,
		},
	}
	res := s.SendMsgAndWait(chain, msg, "IBC transfer tokens for e2e tests")
	return s.FindEventAttr(res, "send_packet", "packet_sequence")
}

func (s *TestSuite) GetIBCClientLatestHeight(chain *starship.ChainClient) ibcexported.Height {
	res, err := ibcclienttypes.
		NewQueryClient(chain.Client).
		ClientStates(context.Background(), &ibcclienttypes.QueryClientStatesRequest{})
	s.Require().NoError(err)

	var cs ibcexported.ClientState
	err = chain.Client.Codec.InterfaceRegistry.UnpackAny(res.ClientStates[0].ClientState, &cs)
	s.Require().NoError(err)
	return cs.GetLatestHeight()
}

func (s *TestSuite) GetTransferChannel(chain *starship.ChainClient, counterparty string) *ibcchanneltypes.IdentifiedChannel {
	res, err := ibcchanneltypes.
		NewQueryClient(chain.Client).
		Channels(context.Background(), &ibcchanneltypes.QueryChannelsRequest{})
	s.Require().NoError(err)

	for _, c := range res.GetChannels() {
		if c.PortId == ibctransfertypes.PortID && s.GetChannelClientChainID(chain, c.PortId, c.ChannelId) == counterparty {
			return c
		}
	}
	s.FailNow("transfer channel not found")
	return nil
}

func (s *TestSuite) GetChannelClientChainID(chain *starship.ChainClient, port, channel string) string {
	res, err := ibcchanneltypes.
		NewQueryClient(chain.Client).
		ChannelClientState(context.Background(), &ibcchanneltypes.QueryChannelClientStateRequest{
			PortId:    port,
			ChannelId: channel,
		})
	s.Require().NoError(err)

	var cs ibcexported.ClientState
	err = chain.Client.Codec.InterfaceRegistry.UnpackAny(res.IdentifiedClientState.ClientState, &cs)
	s.Require().NoError(err)

	ts, ok := cs.(*ibctm.ClientState)
	s.Require().True(ok)
	return ts.ChainId
}

func (s *TestSuite) WaitForIBCPacketAck(chain *starship.ChainClient, port, channel, seq string) {
	s.T().Logf("wait for packet ack; port: %s, channel: %s, seq: %s", port, channel, seq)
	seqInt, err := strconv.Atoi(seq)
	s.Require().NoError(err)
	s.Require().Eventuallyf(
		func() bool {
			res, err := ibcchanneltypes.
				NewQueryClient(chain.Client).
				PacketReceipt(context.Background(), &ibcchanneltypes.QueryPacketReceiptRequest{
					PortId:    port,
					ChannelId: channel,
					Sequence:  uint64(seqInt),
				})
			s.Require().NoError(err)
			return res.Received
		},
		300*time.Second,
		time.Second,
		fmt.Sprintf("waited for too long, still packet not received; port: %s, channel: %s, seq: %s", port, channel, seq),
	)
}
