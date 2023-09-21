package interchaintest

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v7/relayer"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TestPersistenceGaiaIBCTransfer spins up a Persistence and Gaia network, initializes an IBC connection between them,
// and sends an ICS20 token transfer from Persistence->Gaia and then back from Gaia->Persistence.
func TestPersistenceGaiaIBCTransfer(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// Create chain factory with Persistence and Gaia
	numVals := 1
	numFullNodes := 1

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "persistence",
			ChainConfig:   persistenceChainConfig(),
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "gaia",
			Version:       "v9.1.0",
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
	})

	const (
		path = "ibc-path"
	)

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	client, network := interchaintest.DockerSetup(t)

	persistenceChain, gaiaChain := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	relayerType, relayerName := ibc.CosmosRly, "relay"

	// Get a relayer instance
	rf := interchaintest.NewBuiltinRelayerFactory(
		relayerType,
		zaptest.NewLogger(t),
		interchaintestrelayer.CustomDockerImage(IBCRelayerImage, IBCRelayerVersion, "100:1000"),
		interchaintestrelayer.StartupFlags("--processor", "events", "--block-history", "100"),
	)

	r := rf.Build(t, client, network)

	ic := interchaintest.NewInterchain().
		AddChain(persistenceChain).
		AddChain(gaiaChain).
		AddRelayer(r, relayerName).
		AddLink(interchaintest.InterchainLink{
			Chain1:  persistenceChain,
			Chain2:  gaiaChain,
			Relayer: r,
			Path:    path,
		})

	ctx := context.Background()

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation:  false,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Create some user accounts on both chains
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), genesisWalletAmount, persistenceChain, gaiaChain)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, persistenceChain, gaiaChain)
	require.NoError(t, err)

	// Get our Bech32 encoded user addresses
	persistenceUser, gaiaUser := users[0], users[1]

	persistenceUserAddr := persistenceUser.FormattedAddress()
	gaiaUserAddr := gaiaUser.FormattedAddress()

	// Get original account balances
	persistenceOrigBal, err := persistenceChain.GetBalance(ctx, persistenceUserAddr, persistenceChain.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(genesisWalletAmount), persistenceOrigBal)

	gaiaOrigBal, err := gaiaChain.GetBalance(ctx, gaiaUserAddr, gaiaChain.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(genesisWalletAmount), gaiaOrigBal)

	// Compose an IBC transfer and send from Persistence -> Gaia
	var transferAmount = math.NewInt(1_000)
	transfer := ibc.WalletAmount{
		Address: gaiaUserAddr,
		Denom:   persistenceChain.Config().Denom,
		Amount:  transferAmount,
	}

	channel, err := ibc.GetTransferChannel(ctx, r, eRep, persistenceChain.Config().ChainID, gaiaChain.Config().ChainID)
	require.NoError(t, err)

	persistenceHeight, err := persistenceChain.Height(ctx)
	require.NoError(t, err)

	transferTx, err := persistenceChain.SendIBCTransfer(ctx, channel.ChannelID, persistenceUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)

	err = r.StartRelayer(ctx, eRep, path)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				t.Logf("an error occured while stopping the relayer: %s", err)
			}
		},
	)

	// Poll for the ack to know the transfer was successful
	_, err = testutil.PollForAck(ctx, persistenceChain, persistenceHeight, persistenceHeight+50, transferTx.Packet)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 10, persistenceChain)
	require.NoError(t, err)

	// Get the IBC denom for uxprt on Gaia
	persistenceTokenDenom := transfertypes.GetPrefixedDenom(channel.Counterparty.PortID, channel.Counterparty.ChannelID, persistenceChain.Config().Denom)
	persistenceIBCDenom := transfertypes.ParseDenomTrace(persistenceTokenDenom).IBCDenom()

	// Assert that the funds are no longer present in user acc on Persistence and are in the user acc on Gaia
	persistenceUpdateBal, err := persistenceChain.GetBalance(ctx, persistenceUserAddr, persistenceChain.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, persistenceOrigBal.Sub(transferAmount), persistenceUpdateBal)

	gaiaUpdateBal, err := gaiaChain.GetBalance(ctx, gaiaUserAddr, persistenceIBCDenom)
	require.NoError(t, err)
	require.Equal(t, transferAmount, gaiaUpdateBal)

	// Compose an IBC transfer and send from Gaia -> Persistence
	transfer = ibc.WalletAmount{
		Address: persistenceUserAddr,
		Denom:   persistenceIBCDenom,
		Amount:  transferAmount,
	}

	gaiaHeight, err := gaiaChain.Height(ctx)
	require.NoError(t, err)

	transferTx, err = gaiaChain.SendIBCTransfer(ctx, channel.Counterparty.ChannelID, gaiaUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)

	// Poll for the ack to know the transfer was successful
	_, err = testutil.PollForAck(ctx, gaiaChain, gaiaHeight, gaiaHeight+25, transferTx.Packet)
	require.NoError(t, err)

	// Assert that the funds are now back on Persistence and not on Gaia
	persistenceUpdateBal, err = persistenceChain.GetBalance(ctx, persistenceUserAddr, persistenceChain.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, persistenceOrigBal, persistenceUpdateBal)

	gaiaUpdateBal, err = gaiaChain.GetBalance(ctx, gaiaUserAddr, persistenceIBCDenom)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(0), gaiaUpdateBal)
}
