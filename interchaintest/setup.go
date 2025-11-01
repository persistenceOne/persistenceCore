package interchaintest

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	liquidtypes "github.com/cosmos/gaia/v25/x/liquid/types"
	"github.com/cosmos/interchaintest/v10"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	"github.com/cosmos/interchaintest/v10/ibc"
	"github.com/cosmos/interchaintest/v10/testreporter"
	liquidstaketypes "github.com/persistenceOne/pstake-native/v6/x/liquidstake/types"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/persistenceOne/persistenceCore/v16/interchaintest/helpers"
)

var (
	PersistenceE2ERepo = "persistenceone/persistencecore"
	IBCRelayerImage    = "ghcr.io/cosmos/relayer"
	IBCRelayerVersion  = "main"

	HermesRelayerImage   = "ghcr.io/informalsystems/hermes"
	HermesRelayerVersion = "main"

	PersistenceCoreImage = ibc.DockerImage{
		Repository: "persistence",
		Version:    "local",
		UIDGID:     "1025:1025",
	}

	defaultGenesisOverridesKV = []cosmos.GenesisKV{
		{
			Key:   "app_state.gov.params.voting_period",
			Value: "15s",
		},
		{
			Key:   "app_state.gov.params.max_deposit_period",
			Value: "10s",
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.denom",
			Value: helpers.PersistenceBondDenom,
		},
	}

	votingGenesisOverridesKV = []cosmos.GenesisKV{
		{
			Key:   "app_state.gov.params.voting_period",
			Value: "600s",
		},
		{
			Key:   "app_state.gov.params.max_deposit_period",
			Value: "5s",
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.denom",
			Value: helpers.PersistenceBondDenom,
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.amount",
			Value: "10",
		},
	}

	fastVotingGenesisOverridesKV = []cosmos.GenesisKV{
		{
			Key:   "app_state.gov.params.voting_period",
			Value: "10s",
		},
		{
			Key:   "app_state.gov.params.max_deposit_period",
			Value: "5s",
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.denom",
			Value: helpers.PersistenceBondDenom,
		},
		{
			Key:   "app_state.gov.params.min_deposit.0.amount",
			Value: "10",
		},
	}

	genesisWalletAmount       = math.NewInt(10_000_000)
	cosmosGenesisWalletAmount = math.NewInt(1_000_000_000_000)
)

// PersistenceEncoding registers the persistenceCore specific module codecs so that the associated types and msgs
// will be supported when writing to the blocksdb sqlite database.
func PersistenceEncoding() *testutil.TestEncodingConfig {
	cfg := cosmos.DefaultEncoding()

	// register custom types
	wasmtypes.RegisterInterfaces(cfg.InterfaceRegistry)
	liquidstaketypes.RegisterInterfaces(cfg.InterfaceRegistry)
	liquidtypes.RegisterInterfaces(cfg.InterfaceRegistry)

	return &cfg
}

// persistenceChainConfig returns dynamic config for persistence chains, allowing to inject genesis overrides
func persistenceChainConfig(
	genesisOverrides ...cosmos.GenesisKV,
) ibc.ChainConfig {
	if len(genesisOverrides) == 0 {
		genesisOverrides = defaultGenesisOverridesKV
	}

	config := ibc.ChainConfig{
		Type:                "cosmos",
		Name:                "persistence",
		ChainID:             "ictest-core-1",
		Bin:                 "persistenceCore",
		Bech32Prefix:        "persistence",
		Denom:               helpers.PersistenceBondDenom,
		CoinType:            fmt.Sprintf("%d", helpers.PersistenceCoinType),
		GasPrices:           fmt.Sprintf("0%s", helpers.PersistenceBondDenom),
		GasAdjustment:       1.5,
		TrustingPeriod:      "112h",
		NoHostMount:         false,
		ConfigFileOverrides: nil,
		EncodingConfig:      PersistenceEncoding(),
		ModifyGenesis:       cosmos.ModifyGenesis(genesisOverrides),

		Images: []ibc.DockerImage{
			PersistenceCoreImage,
		},
	}

	return config
}

// func InitChains(t *testing.T, numVals, numFull int, genesisOverrides ...cosmos.GenesisKV) []ibc.Chain {
// 	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
// 		{
// 			Name:          "persistence",
// 			ChainName:     "persistence",
// 			Version:       appVersion,
// 			ChainConfig:   persistenceConfig,
// 			NumValidators: &numVals,
// 			NumFullNodes:  &numFull,
// 		},
// 	})

// 	// Get chains from the chain factory
// 	chains, err := cf.Chains(t.Name())
// 	require.NoError(t, err)

// 	// chain := chains[0].(*cosmos.CosmosChain)
// 	return chains
// }

// func InitInterchain(t *testing.T, ctx context.Context, chains []ibc.Chain) *interchaintest.Interchain {
// 	ic := interchaintest.NewInterchain()
// 	for _, chain := range chains {
// 		ic = ic.AddChain(chain)
// 	}

// 	client, network := interchaintest.DockerSetup(t)
// 	err := ic.Build(
// 		ctx,
// 		testreporter.NewNopReporter().RelayerExecReporter(t),
// 		interchaintest.InterchainBuildOptions{
// 			TestName:         t.Name(),
// 			Client:           client,
// 			NetworkID:        network,
// 			SkipPathCreation: true,
// 			// This can be used to write to the block database which will index all block data e.g. txs, msgs, events, etc.
// 			// BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),

// 		},
// 	)
// 	require.NoError(t, err)

// 	return ic
// }

func CreateChain(
	t *testing.T,
	ctx context.Context,
	numVals, numFull int,
	genesisOverrides ...cosmos.GenesisKV,
) (*interchaintest.Interchain, *cosmos.CosmosChain) {
	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:          "persistence",
			ChainName:     "persistence",
			Version:       PersistenceCoreImage.Version,
			ChainConfig:   persistenceChainConfig(genesisOverrides...),
			NumValidators: &numVals,
			NumFullNodes:  &numFull,
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	ic := interchaintest.NewInterchain().AddChain(chains[0])
	client, network := interchaintest.DockerSetup(t)

	err = ic.Build(
		ctx,
		testreporter.NewNopReporter().RelayerExecReporter(t),
		interchaintest.InterchainBuildOptions{
			TestName:         t.Name(),
			Client:           client,
			NetworkID:        network,
			SkipPathCreation: true,
		},
	)
	require.NoError(t, err)

	return ic, chains[0].(*cosmos.CosmosChain)
}

func firstUserName(prefix string) string {
	return prefix + "-user1"
}

func secondUserName(prefix string) string {
	return prefix + "-user2"
}
