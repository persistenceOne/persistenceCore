/*
 Copyright [2019] - [2020], PERSISTENCE TECHNOLOGIES PTE. LTD. and the assetMantle contributors
 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"fmt"
	//"github.com/CosmWasm/wasmd/app"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/persistenceOne/assetMantle/application"
	"github.com/persistenceOne/persistenceSDK/schema/helpers/base"
	keysAdd "github.com/persistenceOne/persistenceSDK/utilities/rest/keys/add"
	"github.com/persistenceOne/persistenceSDK/utilities/rest/queuing"
	"github.com/persistenceOne/persistenceSDK/utilities/rest/queuing/rest"
	"github.com/persistenceOne/persistenceSDK/utilities/rest/signTx"
	"os"
	"path"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authCLI "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authREST "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankCLI "github.com/cosmos/cosmos-sdk/x/bank/client/cli"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
)

func main() {
	cobra.EnableCommandSorting = false

	config := sdkTypes.GetConfig()
	config.SetBech32PrefixForAccount(sdkTypes.Bech32PrefixAccAddr, sdkTypes.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdkTypes.Bech32PrefixValAddr, sdkTypes.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdkTypes.Bech32PrefixConsAddr, sdkTypes.Bech32PrefixConsPub)
	config.Seal()

	rootCommand := &cobra.Command{
		Use:   "client",
		Short: "Command line interface for interacting with node",
	}

	rootCommand.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCommand.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initializeConfiguration(rootCommand)
	}

	rootCommand.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(application.DefaultClientHome),
		queryCommand(application.Codec),
		transactionCommand(application.Codec),
		flags.LineBreak,
		ServeCmd(application.Codec),
		flags.LineBreak,
		keys.Commands(),
		flags.LineBreak,
		version.Cmd,
		flags.NewCompletionCmd(rootCommand, true),
	)

	executor := cli.PrepareMainCmd(rootCommand, "HC", application.DefaultClientHome)

	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func registerRoutes(restServer *lcd.RestServer) {
	client.RegisterRoutes(restServer.CliCtx, restServer.Mux)
	authREST.RegisterTxRoutes(restServer.CliCtx, restServer.Mux)
	application.ModuleBasics.RegisterRESTRoutes(restServer.CliCtx, restServer.Mux)
	keysAdd.RegisterRESTRoutes(restServer.CliCtx, restServer.Mux)
	signTx.RegisterRESTRoutes(restServer.CliCtx, restServer.Mux)
}

func queryCommand(codec *amino.Codec) *cobra.Command {
	queryCommand := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Root command for querying.",
	}

	queryCommand.AddCommand(
		authCLI.GetAccountCmd(codec),
		flags.LineBreak,
		rpc.ValidatorCommand(codec),
		rpc.BlockCommand(),
		authCLI.QueryTxsByEventsCmd(codec),
		authCLI.QueryTxCmd(codec),
		flags.LineBreak,
	)

	application.ModuleBasics.AddQueryCommands(queryCommand, codec)

	return queryCommand
}

// ServeCommand will start the application REST service as a blocking process. It
// takes a codec to create a RestServer object and a function to register all
// necessary routes.
func ServeCmd(codec *amino.Codec) *cobra.Command {
	flagKafka := "kafka"
	kafkaPorts := "kafkaPort"
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			generateOnly := viper.GetBool(flags.FlagGenerateOnly)
			viper.Set(flags.FlagGenerateOnly, false)
			rs := lcd.NewRestServer(codec)
			viper.Set(flags.FlagGenerateOnly, generateOnly)
			kafkaBool := viper.GetBool(flagKafka)
			var kafkaState queuing.KafkaState
			corsBool := viper.GetBool(flags.FlagUnsafeCORS)
			if kafkaBool == true {
				kafkaPort := viper.GetString(kafkaPorts)
				kafkaPort = strings.Trim(kafkaPort, "\" ")
				kafkaPorts := strings.Split(kafkaPort, " ")
				kafkaState = queuing.NewKafkaState(kafkaPorts)
				base.KafkaBool = kafkaBool
				base.KafkaState = kafkaState
				rs.Mux.HandleFunc("/response/{ticketID}", queuing.QueryDB(codec, kafkaState.KafkaDB)).Methods("GET")
			}
			registerRoutes(rs)
			if kafkaBool == true {
				go func() {
					for {
						rest.KafkaConsumerMessages(rs.CliCtx, kafkaState)
						time.Sleep(queuing.SleepRoutine)
					}
				}()
			}
			// Start the rest server and return error if one exists
			err = rs.Start(
				viper.GetString(flags.FlagListenAddr),
				viper.GetInt(flags.FlagMaxOpenConnections),
				uint(viper.GetInt(flags.FlagRPCReadTimeout)),
				uint(viper.GetInt(flags.FlagRPCWriteTimeout)),
				corsBool,
			)
			return err
		},
	}
	cmd.Flags().Bool(flagKafka, false, "Whether have kafka running")
	cmd.Flags().String(kafkaPorts, "localhost:9092", "Space separated addresses in quotes of the kafka listening node: example: --kafkaPort \"addr1 addr2\" ")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|kwallet|pass|test)")
	cmd.Flags().Bool(flags.FlagGenerateOnly, false, "Build an unsigned transaction and write it as response to rest (when enabled, the local Keybase is not accessible and the node operates offline)")
	cmd.Flags().StringP(flags.FlagBroadcastMode, "b", flags.BroadcastSync, "Transaction broadcasting mode (sync|async|block)")

	return flags.RegisterRestServerFlags(cmd)
}

func transactionCommand(codec *amino.Codec) *cobra.Command {
	transactionCommand := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	transactionCommand.AddCommand(
		bankCLI.SendTxCmd(codec),
		flags.LineBreak,
		authCLI.GetSignCommand(codec),
		authCLI.GetMultiSignCommand(codec),
		flags.LineBreak,
		authCLI.GetBroadcastCommand(codec),
		authCLI.GetEncodeCommand(codec),
		flags.LineBreak,
	)

	application.ModuleBasics.AddTxCommands(transactionCommand, codec)

	var commandListToRemove []*cobra.Command

	for _, cmd := range transactionCommand.Commands() {
		if cmd.Use == auth.ModuleName || cmd.Use == bank.ModuleName {
			commandListToRemove = append(commandListToRemove, cmd)
		}
	}

	transactionCommand.RemoveCommand(commandListToRemove...)

	return transactionCommand
}

func initializeConfiguration(command *cobra.Command) error {
	home, err := command.PersistentFlags().GetString(cli.HomeFlag)
	if err != nil {
		return err
	}

	configurationFile := path.Join(home, "configuration", "configuration.toml")
	if _, err := os.Stat(configurationFile); err == nil {
		viper.SetConfigFile(configurationFile)

		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	}
	if err := viper.BindPFlag(flags.FlagChainID, command.PersistentFlags().Lookup(flags.FlagChainID)); err != nil {
		return err
	}
	if err := viper.BindPFlag(cli.EncodingFlag, command.PersistentFlags().Lookup(cli.EncodingFlag)); err != nil {
		return err
	}
	return viper.BindPFlag(cli.OutputFlag, command.PersistentFlags().Lookup(cli.OutputFlag))
}
