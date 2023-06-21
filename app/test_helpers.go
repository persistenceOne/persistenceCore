package app

import (
	"fmt"
	"os"

	"github.com/CosmWasm/wasmd/x/wasm"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	pruningtypes "github.com/cosmos/cosmos-sdk/store/pruning/types"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
)

// NewTestNetworkFixture returns a new persistenceCore AppConstructor for network simulation tests.
func NewTestNetworkFixture() network.TestFixture {
	dir, err := os.MkdirTemp("", "persistenceCore")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	defer os.RemoveAll(dir)

	app := NewApplication(log.NewNopLogger(), dbm.NewMemDB(), nil, true, GetEnabledProposals(), simtestutil.NewAppOptionsWithFlagHome(dir), []wasm.Option{})

	appCtr := func(val network.ValidatorI) servertypes.Application {
		return NewApplication(
			val.GetCtx().Logger, dbm.NewMemDB(), nil, true, GetEnabledProposals(),
			simtestutil.NewAppOptionsWithFlagHome(val.GetCtx().Config.RootDir), []wasm.Option{},
			bam.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			bam.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			bam.SetChainID(val.GetCtx().Viper.GetString(flags.FlagChainID)),
		)
	}

	return network.TestFixture{
		AppConstructor: appCtr,
		GenesisState:   app.DefaultGenesis(),
		EncodingConfig: testutil.TestEncodingConfig{
			InterfaceRegistry: app.InterfaceRegistry(),
			Codec:             app.AppCodec(),
			TxConfig:          app.TxConfig(),
			Amino:             app.LegacyAmino(),
		},
	}
}
