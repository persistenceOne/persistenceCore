/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"encoding/json"
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/persistenceOne/persistenceCore/application"
	applicationParams "github.com/persistenceOne/persistenceCore/application/params"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	tendermintABCITypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tendermintTypes "github.com/tendermint/tendermint/types"
	tendermintDB "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/persistenceOne/persistenceCore/application/initialize"
)

const flagInvalidCheckPeriod = "invalid-check-period"

var invalidCheckPeriod uint

func main() {

	configuration := sdkTypes.GetConfig()
	configuration.SetBech32PrefixForAccount(application.Bech32PrefixAccAddr, application.Bech32PrefixAccPub)
	configuration.SetBech32PrefixForValidator(application.Bech32PrefixValAddr, application.Bech32PrefixValPub)
	configuration.SetBech32PrefixForConsensusNode(application.Bech32PrefixConsAddr, application.Bech32PrefixConsPub)
	configuration.SetCoinType(application.CoinType)
	configuration.SetFullFundraiserPath(application.FullFundraiserPath)
	configuration.Seal()

	encodingConfig := applicationParams.MakeEncodingConfig()
	initClientCtx := client.Context{}.
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TransactionConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authTypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(application.DefaultNodeHome)

	cobra.EnableCommandSorting = false

	rootCommand := &cobra.Command{
		Use:   "persistenceNode",
		Short: "Persistence Hub Node Daemon (server)",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if err := client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			return server.InterceptConfigsPreRunHandler(cmd)
		},
	}

	rootCommand.AddCommand(initialize.Command(
		application.ModuleBasics,
		application.DefaultNodeHome,
	))
	rootCommand.AddCommand(initialize.CollectGenesisTransactionsCommand(
		bankTypes.GenesisBalancesIterator{},
		application.DefaultNodeHome,
	))
	rootCommand.AddCommand(initialize.MigrateGenesisCommand())
	rootCommand.AddCommand(initialize.GenesisTransactionCommand(
		application.ModuleBasics,
		encodingConfig.TransactionConfig,
		bankTypes.GenesisBalancesIterator{},
		application.DefaultNodeHome,
	))
	rootCommand.AddCommand(initialize.ValidateGenesisCommand(
		application.ModuleBasics,
	))
	rootCommand.AddCommand(initialize.AddGenesisAccountCommand(
		encodingConfig.Marshaler,
		application.DefaultNodeHome,
	))
	rootCommand.AddCommand(flags.NewCompletionCmd(rootCommand, true))
	rootCommand.AddCommand(debug.Cmd())
	rootCommand.AddCommand(version.NewVersionCommand())
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
	) tendermintABCITypes.Application {
		var cache sdkTypes.MultiStorePersistentCache

		if viper.GetBool(server.FlagInterBlockCache) {
			cache = store.NewCommitKVStoreCacheManager()
		}

		skipUpgradeHeights := make(map[int64]bool)
		for _, h := range viper.GetIntSlice(server.FlagUnsafeSkipUpgrades) {
			skipUpgradeHeights[int64(h)] = true
		}
		pruningOpts, err := server.GetPruningOptionsFromFlags()
		if err != nil {
			panic(err)
		}
		return application.NewApplication().Initialize(
			application.Name,
			application.Codec,
			wasm.EnableAllProposals,
			application.ModuleAccountPermissions,
			application.TokenReceiveAllowedModules,
			logger,
			db,
			traceStore,
			true,
			invalidCheckPeriod,
			skipUpgradeHeights,
			viper.GetString(flags.FlagHome),
			baseapp.SetPruning(pruningOpts),
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
			genesisApplication := application.NewApplication().Initialize(
				application.Name,
				application.Codec,
				wasm.EnableAllProposals,
				application.ModuleAccountPermissions,
				application.TokenReceiveAllowedModules,
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
		genesisApplication := application.NewApplication().Initialize(
			application.Name,
			application.Codec,
			wasm.EnableAllProposals,
			application.ModuleAccountPermissions,
			application.TokenReceiveAllowedModules,
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
		application.Codec,
		rootCommand,
		appCreator,
		appExporter,
	)

	executor := cli.PrepareBaseCmd(rootCommand, "CA", application.DefaultNodeHome)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}
