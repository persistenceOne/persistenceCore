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
	persistence := s.GetChainClient("test-core-1")
	gaia := s.GetChainClient("test-gaia-1")
	uatom := gaia.MustGetChainDenom()
	amt := 10000000    // 10 ATOMs
	amtBack := 1000000 // 1 ATOM

	randAddr, err := persistence.CreateRandWallet("test-ibc-transfer")
	s.Require().NoError(err)

	c := s.GetTransferChannel(persistence, gaia.ChainID)
	port, channel := c.PortId, c.ChannelId

	// ibc atom denom would be empty if there are no traces of this uatom yet
	ibcAtom := s.GetIBCDenom(persistence, port, channel, uatom)
	balBefore := s.GetBalance(persistence, persistence.Address, ibcAtom)

	s.T().Logf("transferring %d%s to persistence address", amt, uatom)
	seq := s.IBCTransferTokens(gaia, gaia.Address, randAddr, port, channel, uatom, amt)
	s.WaitForIBCPacketAck(persistence, port, channel, seq)

	ibcAtom = s.GetIBCDenom(persistence, port, channel, uatom)
	balAfter := s.GetBalance(persistence, persistence.Address, ibcAtom)

	s.T().Log("check balance after ibc transfer")
	s.Require().Equal(balBefore.Amount.Add(sdk.NewInt(int64(amt))), balAfter.Amount)

	s.T().Logf("transferring back %d%s back to gaia address (in ibc denoms)", amtBack, uatom)
	port, channel = c.Counterparty.PortId, c.Counterparty.ChannelId
	gaiaBalBefore := s.GetBalance(gaia, gaia.Address, uatom)
	seq = s.IBCTransferTokens(persistence, randAddr, gaia.Address, port, channel, ibcAtom, amtBack)
	s.WaitForIBCPacketAck(persistence, port, channel, seq)

	gaiaBalAfter := s.GetBalance(gaia, gaia.Address, uatom)
	s.T().Log("check balance after transferring back ibc token")
	s.Require().Equal(gaiaBalBefore.AddAmount(sdk.NewInt(int64(amtBack))), gaiaBalAfter)
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

func (s *TestSuite) IBCTransferTokens(chain *starship.ChainClient, sender, receiver, port, channel, denom string, amount int) string {
	coin, err := sdk.ParseCoinNormalized(fmt.Sprintf("%d%s", amount, denom))
	s.Require().NoError(err)

	lh := s.GetIBCClientLatestHeight(chain)
	msg := &ibctransfertypes.MsgTransfer{
		SourcePort:    port,
		SourceChannel: channel,
		Token:         coin,
		Sender:        sender,
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
	s.T().Logf("wait for packet act; port: %s, channel: %s, seq: %s", port, channel, seq)
	seqInt, err := strconv.Atoi(seq)
	s.Require().NoError(err)
	s.Require().Eventuallyf(
		func() bool {
			_, err := ibcchanneltypes.
				NewQueryClient(chain.Client).
				PacketAcknowledgement(context.Background(), &ibcchanneltypes.QueryPacketAcknowledgementRequest{
					PortId:    port,
					ChannelId: channel,
					Sequence:  uint64(seqInt),
				})
			return err == nil
		},
		300*time.Second,
		time.Second,
		fmt.Sprintf("waited for too long, still no packet ack; port: %s, channel: %s, seq: %s", port, channel, seq),
	)
}
