package main

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"io"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/persistenceOne/persistenceSDK/applications/hub"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	tendermintABSITypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tendermintTypes "github.com/tendermint/tendermint/types"
	tendermintDB "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/persistenceOne/persistenceSDK/applications/hub/initialize"
)

const flagInvalidCheckPeriod = "invalid-check-period"

var invalidCheckPeriod uint

func main() {

	serverContext := server.NewDefaultContext()

	codec := hub.MakeCodec()

	configuration := sdkTypes.GetConfig()
	configuration.SetBech32PrefixForAccount(sdkTypes.Bech32PrefixAccAddr, sdkTypes.Bech32PrefixAccPub)
	configuration.SetBech32PrefixForValidator(sdkTypes.Bech32PrefixValAddr, sdkTypes.Bech32PrefixValPub)
	configuration.SetBech32PrefixForConsensusNode(sdkTypes.Bech32PrefixConsAddr, sdkTypes.Bech32PrefixConsPub)
	configuration.Seal()

	cobra.EnableCommandSorting = false

	rootCommand := &cobra.Command{
		Use:               "hubNode",
		Short:             "Persistence Hub Node Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(serverContext),
	}

	rootCommand.AddCommand(initialize.InitializeCommand(
		serverContext,
		codec,
		hub.ModuleBasics,
		hub.DefaultNodeHome,
	))
	rootCommand.AddCommand(initialize.CollectGenesisTransactionsCommand(
		serverContext,
		codec,
		auth.GenesisAccountIterator{},
		hub.DefaultNodeHome,
	))
	rootCommand.AddCommand(initialize.MigrateGenesisCommand(
		serverContext,
		codec,
	))
	rootCommand.AddCommand(initialize.GenesisTransactionCommand(
		serverContext,
		codec,
		hub.ModuleBasics,
		staking.AppModuleBasic{},
		auth.GenesisAccountIterator{},
		hub.DefaultNodeHome,
		hub.DefaultClientHome,
	))
	rootCommand.AddCommand(initialize.ValidateGenesisCommand(
		serverContext,
		codec,
		hub.ModuleBasics,
	))
	rootCommand.AddCommand(initialize.AddGenesisAccountCommand(
		serverContext,
		codec,
		hub.DefaultNodeHome,
		hub.DefaultClientHome,
	))
	rootCommand.AddCommand(flags.NewCompletionCmd(rootCommand, true))
	rootCommand.AddCommand(initialize.TestnetCommand(
		serverContext,
		codec,
		hub.ModuleBasics,
		auth.GenesisAccountIterator{},
	))
	rootCommand.AddCommand(initialize.ReplayTransactionsCommand())
	rootCommand.AddCommand(debug.Cmd(codec))

	rootCommand.PersistentFlags().UintVar(
		&invalidCheckPeriod,
		flagInvalidCheckPeriod,
		0,
		"Assert registered invariants every N blocks",
	)

	appCreator := func(
		logger log.Logger,
		db tendermintDB.DB,
		traceStore io.Writer,
	) tendermintABSITypes.Application {
		var cache sdkTypes.MultiStorePersistentCache

		if viper.GetBool(server.FlagInterBlockCache) {
			cache = store.NewCommitKVStoreCacheManager()
		}

		skipUpgradeHeights := make(map[int64]bool)
		for _, h := range viper.GetIntSlice(server.FlagUnsafeSkipUpgrades) {
			skipUpgradeHeights[int64(h)] = true
		}
		return hub.NewPersistenceHubApplication(
			logger,
			db,
			traceStore,
			true,
			invalidCheckPeriod,
			skipUpgradeHeights,
			viper.GetString(flags.FlagHome),
			baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
			baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
			baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
			baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),
			baseapp.SetInterBlockCache(cache),
		)
	}

	appExporter := func(
		logger log.Logger,
		db tendermintDB.DB,
		traceStore io.Writer,
		height int64,
		forZeroHeight bool,
		jailWhiteList []string,
	) (json.RawMessage, []tendermintTypes.GenesisValidator, error) {

		if height != -1 {
			genesisApplication := hub.NewPersistenceHubApplication(
				logger,
				db,
				traceStore,
				false,
				uint(1),
				map[int64]bool{},
				"",
			)
			err := genesisApplication.LoadHeight(height)
			if err != nil {
				return nil, nil, err
			}
			return genesisApplication.ExportApplicationStateAndValidators(forZeroHeight, jailWhiteList)
		}
		//else
		genesisApplication := hub.NewPersistenceHubApplication(
			logger,
			db,
			traceStore,
			true,
			uint(1),
			map[int64]bool{},
			"",
		)
		return genesisApplication.ExportApplicationStateAndValidators(forZeroHeight, jailWhiteList)

	}

	server.AddCommands(
		serverContext,
		codec,
		rootCommand,
		appCreator,
		appExporter,
	)

	executor := cli.PrepareBaseCmd(rootCommand, "CA", hub.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
