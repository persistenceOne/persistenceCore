package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"os"
	"path"

	"github.com/persistenceOne/persistenceSDK/applications/hub"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/cli"
)

func main() {
	cobra.EnableCommandSorting = false
	codec := hub.MakeCodec()

	config := sdkTypes.GetConfig()
	config.SetBech32PrefixForAccount(sdkTypes.Bech32PrefixAccAddr, sdkTypes.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdkTypes.Bech32PrefixValAddr, sdkTypes.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdkTypes.Bech32PrefixConsAddr, sdkTypes.Bech32PrefixConsPub)
	config.Seal()

	rootCommand := &cobra.Command{
		Use:   "hubClient",
		Short: "Command line interface for interacting with hubNode",
	}

	rootCommand.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCommand.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return initalizeConfiguration(rootCommand)
	}

	rootCommand.AddCommand(
		rpc.StatusCommand(),
		client.ConfigCmd(hub.DefaultClientHome),
		queryCommand(codec),
		transactionCommand(codec),
		flags.LineBreak,
		lcd.ServeCommand(codec, registerRoutes),
		flags.LineBreak,
		keys.Commands(),
		flags.LineBreak,
		version.Cmd,
		flags.NewCompletionCmd(rootCommand, true),
	)

	executor := cli.PrepareMainCmd(rootCommand, "HC", hub.DefaultClientHome)

	err := executor.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}

func queryCommand(codec *amino.Codec) *cobra.Command {
	queryCommand := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Root command for querying.",
	}

	queryCommand.AddCommand(
		authcmd.GetAccountCmd(codec),
		flags.LineBreak,
		rpc.ValidatorCommand(codec),
		rpc.BlockCommand(),
		authcmd.QueryTxsByEventsCmd(codec),
		authcmd.QueryTxCmd(codec),
		flags.LineBreak,
	)

	hub.ModuleBasics.AddQueryCommands(queryCommand, codec)

	return queryCommand
}

func transactionCommand(codec *amino.Codec) *cobra.Command {
	transactionCommand := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	transactionCommand.AddCommand(
		bankcmd.SendTxCmd(codec),
		flags.LineBreak,
		authcmd.GetSignCommand(codec),
		authcmd.GetMultiSignCommand(codec),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(codec),
		authcmd.GetEncodeCommand(codec),
		flags.LineBreak,
	)

	hub.ModuleBasics.AddTxCommands(transactionCommand, codec)

	var cmdsToRemove []*cobra.Command

	for _, cmd := range transactionCommand.Commands() {
		if cmd.Use == auth.ModuleName || cmd.Use == bank.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	transactionCommand.RemoveCommand(cmdsToRemove...)

	return transactionCommand
}

func registerRoutes(rs *lcd.RestServer) {
	client.RegisterRoutes(rs.CliCtx, rs.Mux)
	authrest.RegisterTxRoutes(rs.CliCtx, rs.Mux)
	hub.ModuleBasics.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
}

func initalizeConfiguration(command *cobra.Command) error {
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
