/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/server"
	serverCmd "github.com/cosmos/cosmos-sdk/server/cmd"
	serverTypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/persistenceOne/persistenceCore/application"
	"github.com/persistenceOne/persistenceCore/application/initialize"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tendermintClient "github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tendermintDB "github.com/tendermint/tm-db"
	"io"
	"os"
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

	encodingConfig := application.MakeEncodingConfig()
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
		Use:   "persistenceCore",
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
		application.DefaultNodeHome,
	))
	rootCommand.AddCommand(tendermintClient.NewCompletionCmd(rootCommand, true))
	rootCommand.AddCommand(debug.Cmd())
	rootCommand.AddCommand(version.NewVersionCommand())
	rootCommand.AddCommand(
		rpc.StatusCommand(),
		queryCommand(),
		txCommand(),
		keys.Commands(application.DefaultNodeHome),
	)
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
		applicationOptions serverTypes.AppOptions,
	) serverTypes.Application {
		var cache sdkTypes.MultiStorePersistentCache

		if viper.GetBool(server.FlagInterBlockCache) {
			cache = store.NewCommitKVStoreCacheManager()
		}

		skipUpgradeHeights := make(map[int64]bool)
		for _, h := range viper.GetIntSlice(server.FlagUnsafeSkipUpgrades) {
			skipUpgradeHeights[int64(h)] = true
		}
		pruningOpts, err := server.GetPruningOptionsFromFlags(applicationOptions)
		if err != nil {
			panic(err)
		}
		return application.NewApplication().Initialize(
			application.Name,
			encodingConfig,
			application.ModuleAccountPermissions,
			logger,
			db,
			traceStore,
			true,
			invalidCheckPeriod,
			skipUpgradeHeights,
			viper.GetString(flags.FlagHome),
			applicationOptions,
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
		applicationOptions serverTypes.AppOptions,
	) (serverTypes.ExportedApp, error) {

		if height != -1 {
			genesisApplication := application.NewApplication().Initialize(
				application.Name,
				encodingConfig,
				application.ModuleAccountPermissions,
				logger,
				db,
				traceStore,
				false,
				uint(1),
				map[int64]bool{},
				"",
				applicationOptions,
			)
			err := genesisApplication.LoadHeight(height)
			if err != nil {
				return serverTypes.ExportedApp{}, err
			}
			return genesisApplication.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
		}
		//else
		genesisApplication := application.NewApplication().Initialize(
			application.Name,
			encodingConfig,
			application.ModuleAccountPermissions,
			logger,
			db,
			traceStore,
			true,
			uint(1),
			map[int64]bool{},
			"",
			applicationOptions,
		)
		return genesisApplication.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)

	}

	initFlags := func(startCmd *cobra.Command) {
		crisis.AddModuleInitFlags(startCmd)
	}

	server.AddCommands(
		rootCommand,
		application.DefaultNodeHome,
		appCreator,
		appExporter,
		initFlags,
	)

	if err := serverCmd.Execute(rootCommand, application.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}

func queryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Querying subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetAccountCmd(),
		rpc.ValidatorCommand(),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(),
		authcmd.QueryTxCmd(),
	)

	application.ModuleBasics.AddQueryCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}

func txCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "tx",
		Short:                      "Transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		authcmd.GetSignCommand(),
		authcmd.GetSignBatchCommand(),
		authcmd.GetMultiSignCommand(),
		authcmd.GetValidateSignaturesCommand(),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(),
		authcmd.GetEncodeCommand(),
		authcmd.GetDecodeCommand(),
		flags.LineBreak,
	)

	application.ModuleBasics.AddTxCommands(cmd)
	cmd.PersistentFlags().String(flags.FlagChainID, "", "The network chain ID")

	return cmd
}
