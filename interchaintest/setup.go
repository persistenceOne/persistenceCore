package interchaintest

import (
	"context"
	"fmt"
	"testing"

	"github.com/docker/docker/client"
	"github.com/persistenceOne/persistenceCore/v8/interchaintest/helpers"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	testutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	ibclocalhost "github.com/cosmos/ibc-go/v7/modules/light-clients/09-localhost"
	interchaintest "github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
)

var (
	VotingPeriod     = "15s"
	MaxDepositPeriod = "10s"

	PersistenceE2ERepo = "persistenceone/persistencecore"

	IBCRelayerImage   = "ghcr.io/cosmos/relayer"
	IBCRelayerVersion = "main"

	appRepo, appVersion = GetDockerImageInfo()

	PersistenceCoreImage = ibc.DockerImage{
		Repository: appRepo,
		Version:    appVersion,
		UidGid:     "1025:1025",
	}

	defaultGenesisOverridesKV = append([]cosmos.GenesisKV{
		{
			Key:   "app_state.gov.params.voting_period",
			Value: VotingPeriod,
		},
		{
			Key:   "app_state.gov.params.max_deposit_period",
			Value: MaxDepositPeriod,
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.denom",
			Value: helpers.PersistenceBondDenom,
		},
		{
			Key:   "app_state.builder.params.reserve_fee.denom",
			Value: helpers.PersistenceBondDenom,
		},
		{
			Key:   "app_state.builder.params.min_bid_increment.denom",
			Value: helpers.PersistenceBondDenom,
		},
	})

	persistenceConfig = ibc.ChainConfig{
		Type:    "cosmos",
		Name:    "persistence",
		ChainID: "ictest-core-1",
		Images: []ibc.DockerImage{
			PersistenceCoreImage,
		},
		Bin:                    "persistenceCore",
		Bech32Prefix:           "persistence",
		Denom:                  helpers.PersistenceBondDenom,
		CoinType:               fmt.Sprintf("%d", helpers.PersistenceCoinType),
		GasPrices:              fmt.Sprintf("0%s", helpers.PersistenceBondDenom),
		GasAdjustment:          1.5,
		TrustingPeriod:         "112h",
		NoHostMount:            false,
		ConfigFileOverrides:    nil,
		EncodingConfig:         persistenceEncoding(),
		UsingNewGenesisCommand: true,
		ModifyGenesis:          cosmos.ModifyGenesis(defaultGenesisOverridesKV),
	}

	genesisWalletAmount = int64(10_000_000)
)

// persistenceEncoding registers the persistenceCore specific module codecs so that the associated types and msgs
// will be supported when writing to the blocksdb sqlite database.
func persistenceEncoding() *testutil.TestEncodingConfig {
	cfg := cosmos.DefaultEncoding()

	// register custom types
	ibclocalhost.RegisterInterfaces(cfg.InterfaceRegistry)
	wasmtypes.RegisterInterfaces(cfg.InterfaceRegistry)

	return &cfg
}

// Base chain, no relaying off this branch (or persistence:local if no branch is provided.)
func CreateThisBranchChain(t *testing.T, numVals, numFull int) []ibc.Chain {
	// Create chain factory with persistence on this current branch

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "persistence",
			ChainName:     "persistence",
			Version:       appVersion,
			ChainConfig:   persistenceConfig,
			NumValidators: &numVals,
			NumFullNodes:  &numFull,
		},
	})

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	// chain := chains[0].(*cosmos.CosmosChain)
	return chains
}

func BuildInitialChain(t *testing.T, chains []ibc.Chain) (*interchaintest.Interchain, context.Context, *client.Client, string) {
	// Create a new Interchain object which describes the chains, relayers, and IBC connections we want to use
	ic := interchaintest.NewInterchain()

	for _, chain := range chains {
		ic = ic.AddChain(chain)
	}

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	ctx := context.Background()
	client, network := interchaintest.DockerSetup(t)

	err := ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: true,
		// This can be used to write to the block database which will index all block data e.g. txs, msgs, events, etc.
		// BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),

	})
	require.NoError(t, err)

	return ic, ctx, client, network
}
