package interchaintest

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	chantypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v7/relayer"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/persistenceOne/persistenceCore/v11/interchaintest/helpers"
)

type PacketMetadata struct {
	Forward *ForwardMetadata `json:"forward"`
}

type ForwardMetadata struct {
	Receiver       string        `json:"receiver"`
	Port           string        `json:"port"`
	Channel        string        `json:"channel"`
	Timeout        time.Duration `json:"timeout"`
	Retries        *uint8        `json:"retries,omitempty"`
	Next           *string       `json:"next,omitempty"`
	RefundSequence *uint64       `json:"refund_sequence,omitempty"`
}

// TestPacketForwardMiddlewareRouter ensures the PFM module is set up properly and works as expected.
func TestPacketForwardMiddlewareRouter(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	var (
		ctx                                        = context.Background()
		client, network                            = interchaintest.DockerSetup(t)
		testReported                               = testreporter.NewNopReporter()
		relayerExecReporter                        = testReported.RelayerExecReporter(t)
		chainID_A, chainID_B, chainID_C, chainID_D = "chain-a", "chain-b", "chain-c", "chain-d"
		chainA, chainB, chainC, chainD             *cosmos.CosmosChain
	)

	// base config which all networks will use as defaults.
	baseCfg := ibc.ChainConfig{
		Type:    "cosmos",
		Name:    "persistence",
		ChainID: "", // change this for each
		Images: []ibc.DockerImage{
			PersistenceCoreImage,
		},
		Bin:                    "persistenceCore",
		Bech32Prefix:           "persistence",
		Denom:                  helpers.PersistenceBondDenom,
		CoinType:               fmt.Sprintf("%d", helpers.PersistenceCoinType),
		GasPrices:              fmt.Sprintf("0%s", helpers.PersistenceBondDenom),
		GasAdjustment:          2.0,
		TrustingPeriod:         "112h",
		NoHostMount:            false,
		ConfigFileOverrides:    nil,
		EncodingConfig:         persistenceEncoding(),
		UsingNewGenesisCommand: true,
		ModifyGenesis:          cosmos.ModifyGenesis(defaultGenesisOverridesKV),
	}

	// Set specific chain ids for each so they are their own unique networks
	baseCfg.ChainID = chainID_A
	configA := baseCfg

	baseCfg.ChainID = chainID_B
	configB := baseCfg

	baseCfg.ChainID = chainID_C
	configC := baseCfg

	baseCfg.ChainID = chainID_D
	configD := baseCfg

	// Create chain factory with multiple Persistence individual networks.
	numVals := 1
	numFullNodes := 0

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "persistence",
			ChainConfig:   configA,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "persistence",
			ChainConfig:   configB,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "persistence",
			ChainConfig:   configC,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "persistence",
			ChainConfig:   configD,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
	})

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chainA, chainB, chainC, chainD = chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain), chains[2].(*cosmos.CosmosChain), chains[3].(*cosmos.CosmosChain)

	relayer := interchaintest.NewBuiltinRelayerFactory(
		ibc.CosmosRly,
		zaptest.NewLogger(t),
		interchaintestrelayer.CustomDockerImage(IBCRelayerImage, IBCRelayerVersion, "100:1000"),
		interchaintestrelayer.StartupFlags("--processor", "events", "--block-history", "100"),
	).Build(t, client, network)

	const pathAB = "ab"
	const pathBC = "bc"
	const pathCD = "cd"

	ic := interchaintest.NewInterchain().
		AddChain(chainA).
		AddChain(chainB).
		AddChain(chainC).
		AddChain(chainD).
		AddRelayer(relayer, "relayer").
		AddLink(interchaintest.InterchainLink{
			Chain1:  chainA,
			Chain2:  chainB,
			Relayer: relayer,
			Path:    pathAB,
		}).
		AddLink(interchaintest.InterchainLink{
			Chain1:  chainB,
			Chain2:  chainC,
			Relayer: relayer,
			Path:    pathBC,
		}).
		AddLink(interchaintest.InterchainLink{
			Chain1:  chainC,
			Chain2:  chainD,
			Relayer: relayer,
			Path:    pathCD,
		})

	require.NoError(t, ic.Build(ctx, relayerExecReporter, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),

		SkipPathCreation: false,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	const userFunds = int64(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, chainA, chainB, chainC, chainD)

	abChan, err := ibc.GetTransferChannel(ctx, relayer, relayerExecReporter, chainID_A, chainID_B)
	require.NoError(t, err)

	baChan := abChan.Counterparty

	cbChan, err := ibc.GetTransferChannel(ctx, relayer, relayerExecReporter, chainID_C, chainID_B)
	require.NoError(t, err)

	bcChan := cbChan.Counterparty

	dcChan, err := ibc.GetTransferChannel(ctx, relayer, relayerExecReporter, chainID_D, chainID_C)
	require.NoError(t, err)

	cdChan := dcChan.Counterparty

	// Start the relayer on all paths
	err = relayer.StartRelayer(ctx, relayerExecReporter, pathAB, pathBC, pathCD)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			err := relayer.StopRelayer(ctx, relayerExecReporter)
			if err != nil {
				t.Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)

	// Get original account balances
	userA, userB, userC, userD := users[0], users[1], users[2], users[3]

	const transferAmount int64 = 100000

	// Compose the prefixed denoms and ibc denom for asserting balances
	firstHopDenom := transfertypes.GetPrefixedDenom(baChan.PortID, baChan.ChannelID, chainA.Config().Denom)
	secondHopDenom := transfertypes.GetPrefixedDenom(cbChan.PortID, cbChan.ChannelID, firstHopDenom)
	thirdHopDenom := transfertypes.GetPrefixedDenom(dcChan.PortID, dcChan.ChannelID, secondHopDenom)

	firstHopDenomTrace := transfertypes.ParseDenomTrace(firstHopDenom)
	secondHopDenomTrace := transfertypes.ParseDenomTrace(secondHopDenom)
	thirdHopDenomTrace := transfertypes.ParseDenomTrace(thirdHopDenom)

	firstHopIBCDenom := firstHopDenomTrace.IBCDenom()
	secondHopIBCDenom := secondHopDenomTrace.IBCDenom()
	thirdHopIBCDenom := thirdHopDenomTrace.IBCDenom()

	firstHopEscrowAccount := sdk.MustBech32ifyAddressBytes(chainA.Config().Bech32Prefix, transfertypes.GetEscrowAddress(abChan.PortID, abChan.ChannelID))
	secondHopEscrowAccount := sdk.MustBech32ifyAddressBytes(chainB.Config().Bech32Prefix, transfertypes.GetEscrowAddress(bcChan.PortID, bcChan.ChannelID))
	thirdHopEscrowAccount := sdk.MustBech32ifyAddressBytes(chainC.Config().Bech32Prefix, transfertypes.GetEscrowAddress(cdChan.PortID, abChan.ChannelID))

	t.Run("multi-hop a->b->c->d", func(t *testing.T) {
		// Send packet from Chain A->Chain B->Chain C->Chain D

		transfer := ibc.WalletAmount{
			Address: userB.FormattedAddress(),
			Denom:   chainA.Config().Denom,
			Amount:  math.NewInt(transferAmount),
		}

		secondHopMetadata := &PacketMetadata{
			Forward: &ForwardMetadata{
				Receiver: userD.FormattedAddress(),
				Channel:  cdChan.ChannelID,
				Port:     cdChan.PortID,
			},
		}
		nextBz, err := json.Marshal(secondHopMetadata)
		require.NoError(t, err)
		next := string(nextBz)

		firstHopMetadata := &PacketMetadata{
			Forward: &ForwardMetadata{
				Receiver: userC.FormattedAddress(),
				Channel:  bcChan.ChannelID,
				Port:     bcChan.PortID,
				Next:     &next,
			},
		}

		memo, err := json.Marshal(firstHopMetadata)
		require.NoError(t, err)

		chainAHeight, err := chainA.Height(ctx)
		require.NoError(t, err)

		transferTx, err := chainA.SendIBCTransfer(ctx, abChan.ChannelID, userA.KeyName(), transfer, ibc.TransferOptions{Memo: string(memo)})
		require.NoError(t, err)
		_, err = testutil.PollForAck(ctx, chainA, chainAHeight, chainAHeight+30, transferTx.Packet)
		require.NoError(t, err)
		err = testutil.WaitForBlocks(ctx, 2, chainA)
		require.NoError(t, err)

		chainABalance, err := chainA.GetBalance(ctx, userA.FormattedAddress(), chainA.Config().Denom)
		require.NoError(t, err)

		chainBBalance, err := chainB.GetBalance(ctx, userB.FormattedAddress(), firstHopIBCDenom)
		require.NoError(t, err)

		chainCBalance, err := chainC.GetBalance(ctx, userC.FormattedAddress(), secondHopIBCDenom)
		require.NoError(t, err)

		chainDBalance, err := chainD.GetBalance(ctx, userD.FormattedAddress(), thirdHopIBCDenom)
		require.NoError(t, err)

		require.Equal(t, math.NewInt(userFunds-transferAmount), chainABalance)
		require.Equal(t, math.NewInt(0), chainBBalance)
		require.Equal(t, math.NewInt(0), chainCBalance)
		require.Equal(t, math.NewInt(transferAmount), chainDBalance)

		firstHopEscrowBalance, err := chainA.GetBalance(ctx, firstHopEscrowAccount, chainA.Config().Denom)
		require.NoError(t, err)

		secondHopEscrowBalance, err := chainB.GetBalance(ctx, secondHopEscrowAccount, firstHopIBCDenom)
		require.NoError(t, err)

		thirdHopEscrowBalance, err := chainC.GetBalance(ctx, thirdHopEscrowAccount, secondHopIBCDenom)
		require.NoError(t, err)

		require.Equal(t, math.NewInt(transferAmount), firstHopEscrowBalance)
		require.Equal(t, math.NewInt(transferAmount), secondHopEscrowBalance)
		require.Equal(t, math.NewInt(transferAmount), thirdHopEscrowBalance)
	})
}

// TestTimeoutOnForward ensures the PFM has proper accounting on the escrow accounts over a timeout event.
func TestTimeoutOnForward(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	var (
		ctx                                    = context.Background()
		client, network                        = interchaintest.DockerSetup(t)
		rep                                    = testreporter.NewNopReporter()
		eRep                                   = rep.RelayerExecReporter(t)
		chainIdA, chainIdB, chainIdC, chainIdD = "chain-a", "chain-b", "chain-c", "chain-d"
	)

	helpers.SetConfig()

	// base config which all networks will use as defaults.
	baseCfg := ibc.ChainConfig{
		Type:    "cosmos",
		Name:    "persistence",
		ChainID: "", // change this for each
		Images: []ibc.DockerImage{
			PersistenceCoreImage,
		},
		Bin:                    "persistenceCore",
		Bech32Prefix:           "persistence",
		Denom:                  helpers.PersistenceBondDenom,
		CoinType:               fmt.Sprintf("%d", helpers.PersistenceCoinType),
		GasPrices:              fmt.Sprintf("0%s", helpers.PersistenceBondDenom),
		GasAdjustment:          2.0,
		TrustingPeriod:         "112h",
		NoHostMount:            false,
		ConfigFileOverrides:    nil,
		EncodingConfig:         persistenceEncoding(),
		UsingNewGenesisCommand: true,
		ModifyGenesis:          cosmos.ModifyGenesis(defaultGenesisOverridesKV),
	}

	// Set specific chain ids for each so they are their own unique networks
	baseCfg.ChainID = chainIdA
	configA := baseCfg

	baseCfg.ChainID = chainIdB
	configB := baseCfg

	baseCfg.ChainID = chainIdC
	configC := baseCfg

	baseCfg.ChainID = chainIdD
	configD := baseCfg

	// Create chain factory with multiple Persistence individual networks.
	numVals := 1
	numFullNodes := 0

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "persistence",
			ChainConfig:   configA,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "persistence",
			ChainConfig:   configB,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "persistence",
			ChainConfig:   configC,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "persistence",
			ChainConfig:   configD,
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chainA, chainB, chainC, chainD := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain), chains[2].(*cosmos.CosmosChain), chains[3].(*cosmos.CosmosChain)

	r := interchaintest.NewBuiltinRelayerFactory(ibc.CosmosRly, zaptest.NewLogger(t)).Build(t, client, network)

	const pathAB = "ab"
	const pathBC = "bc"
	const pathCD = "cd"

	ic := interchaintest.NewInterchain().
		AddChain(chainA).
		AddChain(chainB).
		AddChain(chainC).
		AddChain(chainD).
		AddRelayer(r, "relayer").
		AddLink(interchaintest.InterchainLink{
			Chain1:  chainA,
			Chain2:  chainB,
			Relayer: r,
			Path:    pathAB,
		}).
		AddLink(interchaintest.InterchainLink{
			Chain1:  chainB,
			Chain2:  chainC,
			Relayer: r,
			Path:    pathBC,
		}).
		AddLink(interchaintest.InterchainLink{
			Chain1:  chainC,
			Chain2:  chainD,
			Relayer: r,
			Path:    pathCD,
		})

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: false,
	}))

	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Start the relayer on only the path between chainA<>chainB so that the initial transfer succeeds
	err = r.StartRelayer(ctx, eRep, pathAB)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				t.Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)

	// Fund user accounts with initial balances and get the transfer channel information between each set of chains
	initBal := math.NewInt(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), initBal.Int64(), chainA, chainB, chainC, chainD)

	abChan, err := ibc.GetTransferChannel(ctx, r, eRep, chainIdA, chainIdB)
	require.NoError(t, err)

	baChan := abChan.Counterparty

	cbChan, err := ibc.GetTransferChannel(ctx, r, eRep, chainIdC, chainIdB)
	require.NoError(t, err)

	bcChan := cbChan.Counterparty

	dcChan, err := ibc.GetTransferChannel(ctx, r, eRep, chainIdD, chainIdC)
	require.NoError(t, err)

	cdChan := dcChan.Counterparty

	userA, userB, userC, userD := users[0], users[1], users[2], users[3]

	// Compose the prefixed denoms and ibc denom for asserting balances
	firstHopDenom := transfertypes.GetPrefixedDenom(baChan.PortID, baChan.ChannelID, chainA.Config().Denom)
	secondHopDenom := transfertypes.GetPrefixedDenom(cbChan.PortID, cbChan.ChannelID, firstHopDenom)
	thirdHopDenom := transfertypes.GetPrefixedDenom(dcChan.PortID, dcChan.ChannelID, secondHopDenom)

	firstHopDenomTrace := transfertypes.ParseDenomTrace(firstHopDenom)
	secondHopDenomTrace := transfertypes.ParseDenomTrace(secondHopDenom)
	thirdHopDenomTrace := transfertypes.ParseDenomTrace(thirdHopDenom)

	firstHopIBCDenom := firstHopDenomTrace.IBCDenom()
	secondHopIBCDenom := secondHopDenomTrace.IBCDenom()
	thirdHopIBCDenom := thirdHopDenomTrace.IBCDenom()

	firstHopEscrowAccount := transfertypes.GetEscrowAddress(abChan.PortID, abChan.ChannelID).String()
	secondHopEscrowAccount := transfertypes.GetEscrowAddress(bcChan.PortID, bcChan.ChannelID).String()
	thirdHopEscrowAccount := transfertypes.GetEscrowAddress(cdChan.PortID, abChan.ChannelID).String()

	zeroBal := math.ZeroInt()
	transferAmount := math.NewInt(100_000)

	// Attempt to send packet from Chain A->Chain B->Chain C->Chain D
	transfer := ibc.WalletAmount{
		Address: userB.FormattedAddress(),
		Denom:   chainA.Config().Denom,
		Amount:  transferAmount,
	}

	retries := uint8(0)
	secondHopMetadata := &PacketMetadata{
		Forward: &ForwardMetadata{
			Receiver: userD.FormattedAddress(),
			Channel:  cdChan.ChannelID,
			Port:     cdChan.PortID,
			Retries:  &retries,
		},
	}
	nextBz, err := json.Marshal(secondHopMetadata)
	require.NoError(t, err)
	next := string(nextBz)

	firstHopMetadata := &PacketMetadata{
		Forward: &ForwardMetadata{
			Receiver: userC.FormattedAddress(),
			Channel:  bcChan.ChannelID,
			Port:     bcChan.PortID,
			Next:     &next,
			Retries:  &retries,
			Timeout:  time.Second * 10, // Set low timeout for forward from chainB<>chainC
		},
	}

	memo, err := json.Marshal(firstHopMetadata)
	require.NoError(t, err)

	opts := ibc.TransferOptions{
		Memo: string(memo),
	}

	chainBHeight, err := chainB.Height(ctx)
	require.NoError(t, err)

	transferTx, err := chainA.SendIBCTransfer(ctx, abChan.ChannelID, userA.KeyName(), transfer, opts)
	require.NoError(t, err)

	// Poll for MsgRecvPacket on chainB
	_, err = cosmos.PollForMessage[*chantypes.MsgRecvPacket](ctx, chainB, cosmos.DefaultEncoding().InterfaceRegistry, chainBHeight, chainBHeight+20, nil)
	require.NoError(t, err)

	// Stop the relayer and wait for the timeout to happen on chainC
	err = r.StopRelayer(ctx, eRep)
	require.NoError(t, err)

	time.Sleep(time.Second * 11)

	// Restart the relayer
	err = r.StartRelayer(ctx, eRep, pathAB, pathBC, pathCD)
	require.NoError(t, err)

	chainAHeight, err := chainA.Height(ctx)
	require.NoError(t, err)

	chainBHeight, err = chainB.Height(ctx)
	require.NoError(t, err)

	// Poll for the MsgTimeout on chainB and the MsgAck on chainA
	_, err = cosmos.PollForMessage[*chantypes.MsgTimeout](ctx, chainB, chainB.Config().EncodingConfig.InterfaceRegistry, chainBHeight, chainBHeight+20, nil)
	require.NoError(t, err)

	_, err = testutil.PollForAck(ctx, chainA, chainAHeight, chainAHeight+30, transferTx.Packet)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 1, chainA)
	require.NoError(t, err)

	// Assert balances to ensure that the funds are still on the original sending chain
	chainABalance, err := chainA.GetBalance(ctx, userA.FormattedAddress(), chainA.Config().Denom)
	require.NoError(t, err)

	chainBBalance, err := chainB.GetBalance(ctx, userB.FormattedAddress(), firstHopIBCDenom)
	require.NoError(t, err)

	chainCBalance, err := chainC.GetBalance(ctx, userC.FormattedAddress(), secondHopIBCDenom)
	require.NoError(t, err)

	chainDBalance, err := chainD.GetBalance(ctx, userD.FormattedAddress(), thirdHopIBCDenom)
	require.NoError(t, err)

	require.True(t, chainABalance.Equal(initBal))
	require.True(t, chainBBalance.Equal(zeroBal))
	require.True(t, chainCBalance.Equal(zeroBal))
	require.True(t, chainDBalance.Equal(zeroBal))

	firstHopEscrowBalance, err := chainA.GetBalance(ctx, firstHopEscrowAccount, chainA.Config().Denom)
	require.NoError(t, err)

	secondHopEscrowBalance, err := chainB.GetBalance(ctx, secondHopEscrowAccount, firstHopIBCDenom)
	require.NoError(t, err)

	thirdHopEscrowBalance, err := chainC.GetBalance(ctx, thirdHopEscrowAccount, secondHopIBCDenom)
	require.NoError(t, err)

	require.True(t, firstHopEscrowBalance.Equal(zeroBal))
	require.True(t, secondHopEscrowBalance.Equal(zeroBal))
	require.True(t, thirdHopEscrowBalance.Equal(zeroBal))

	// Send IBC transfer from ChainA -> ChainB -> ChainC -> ChainD that will succeed
	secondHopMetadata = &PacketMetadata{
		Forward: &ForwardMetadata{
			Receiver: userD.FormattedAddress(),
			Channel:  cdChan.ChannelID,
			Port:     cdChan.PortID,
		},
	}
	nextBz, err = json.Marshal(secondHopMetadata)
	require.NoError(t, err)
	next = string(nextBz)

	firstHopMetadata = &PacketMetadata{
		Forward: &ForwardMetadata{
			Receiver: userC.FormattedAddress(),
			Channel:  bcChan.ChannelID,
			Port:     bcChan.PortID,
			Next:     &next,
		},
	}

	memo, err = json.Marshal(firstHopMetadata)
	require.NoError(t, err)

	opts = ibc.TransferOptions{
		Memo: string(memo),
	}

	chainAHeight, err = chainA.Height(ctx)
	require.NoError(t, err)

	transferTx, err = chainA.SendIBCTransfer(ctx, abChan.ChannelID, userA.KeyName(), transfer, opts)
	require.NoError(t, err)

	_, err = testutil.PollForAck(ctx, chainA, chainAHeight, chainAHeight+30, transferTx.Packet)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 5, chainA)
	require.NoError(t, err)

	// Assert balances are updated to reflect tokens now being on ChainD
	chainABalance, err = chainA.GetBalance(ctx, userA.FormattedAddress(), chainA.Config().Denom)
	require.NoError(t, err)

	chainBBalance, err = chainB.GetBalance(ctx, userB.FormattedAddress(), firstHopIBCDenom)
	require.NoError(t, err)

	chainCBalance, err = chainC.GetBalance(ctx, userC.FormattedAddress(), secondHopIBCDenom)
	require.NoError(t, err)

	chainDBalance, err = chainD.GetBalance(ctx, userD.FormattedAddress(), thirdHopIBCDenom)
	require.NoError(t, err)

	require.True(t, chainABalance.Equal(initBal.Sub(transferAmount)))
	require.True(t, chainBBalance.Equal(zeroBal))
	require.True(t, chainCBalance.Equal(zeroBal))
	require.True(t, chainDBalance.Equal(transferAmount))

	firstHopEscrowBalance, err = chainA.GetBalance(ctx, firstHopEscrowAccount, chainA.Config().Denom)
	require.NoError(t, err)

	secondHopEscrowBalance, err = chainB.GetBalance(ctx, secondHopEscrowAccount, firstHopIBCDenom)
	require.NoError(t, err)

	thirdHopEscrowBalance, err = chainC.GetBalance(ctx, thirdHopEscrowAccount, secondHopIBCDenom)
	require.NoError(t, err)

	require.True(t, firstHopEscrowBalance.Equal(transferAmount))
	require.True(t, secondHopEscrowBalance.Equal(transferAmount))
	require.True(t, thirdHopEscrowBalance.Equal(transferAmount))

	// Compose IBC tx that will attempt to go from ChainD -> ChainC -> ChainB -> ChainA but timeout between ChainB->ChainA
	transfer = ibc.WalletAmount{
		Address: userC.FormattedAddress(),
		Denom:   thirdHopDenom,
		Amount:  transferAmount,
	}

	secondHopMetadata = &PacketMetadata{
		Forward: &ForwardMetadata{
			Receiver: userA.FormattedAddress(),
			Channel:  baChan.ChannelID,
			Port:     baChan.PortID,
			Timeout:  1 * time.Second,
		},
	}
	nextBz, err = json.Marshal(secondHopMetadata)
	require.NoError(t, err)
	next = string(nextBz)

	firstHopMetadata = &PacketMetadata{
		Forward: &ForwardMetadata{
			Receiver: userB.FormattedAddress(),
			Channel:  cbChan.ChannelID,
			Port:     cbChan.PortID,
			Next:     &next,
		},
	}

	memo, err = json.Marshal(firstHopMetadata)
	require.NoError(t, err)

	chainDHeight, err := chainD.Height(ctx)
	require.NoError(t, err)

	transferTx, err = chainD.SendIBCTransfer(ctx, dcChan.ChannelID, userD.KeyName(), transfer, ibc.TransferOptions{Memo: string(memo)})
	require.NoError(t, err)

	_, err = testutil.PollForAck(ctx, chainD, chainDHeight, chainDHeight+25, transferTx.Packet)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 5, chainD)
	require.NoError(t, err)

	// Assert balances to ensure timeout happened and user funds are still present on ChainD
	chainABalance, err = chainA.GetBalance(ctx, userA.FormattedAddress(), chainA.Config().Denom)
	require.NoError(t, err)

	chainBBalance, err = chainB.GetBalance(ctx, userB.FormattedAddress(), firstHopIBCDenom)
	require.NoError(t, err)

	chainCBalance, err = chainC.GetBalance(ctx, userC.FormattedAddress(), secondHopIBCDenom)
	require.NoError(t, err)

	chainDBalance, err = chainD.GetBalance(ctx, userD.FormattedAddress(), thirdHopIBCDenom)
	require.NoError(t, err)

	require.True(t, chainABalance.Equal(initBal.Sub(transferAmount)))
	require.True(t, chainBBalance.Equal(zeroBal))
	require.True(t, chainCBalance.Equal(zeroBal))
	require.True(t, chainDBalance.Equal(transferAmount))

	firstHopEscrowBalance, err = chainA.GetBalance(ctx, firstHopEscrowAccount, chainA.Config().Denom)
	require.NoError(t, err)

	secondHopEscrowBalance, err = chainB.GetBalance(ctx, secondHopEscrowAccount, firstHopIBCDenom)
	require.NoError(t, err)

	thirdHopEscrowBalance, err = chainC.GetBalance(ctx, thirdHopEscrowAccount, secondHopIBCDenom)
	require.NoError(t, err)

	require.True(t, firstHopEscrowBalance.Equal(transferAmount))
	require.True(t, secondHopEscrowBalance.Equal(transferAmount))
	require.True(t, thirdHopEscrowBalance.Equal(transferAmount))
}
