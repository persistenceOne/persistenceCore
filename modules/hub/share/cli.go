package share

import (
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/burn"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/lock"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/send"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

func GetCLIRootTransactionCommand(codec *codec.Codec) *cobra.Command {
	rootTransactionCommand := &cobra.Command{
		Use:                        TransactionRoute,
		Short:                      "Share root transaction command.",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	rootTransactionCommand.AddCommand(client.PostCommands(
		burn.TransactionCommand(codec),
		lock.TransactionCommand(codec),
		mint.TransactionCommand(codec),
		send.TransactionCommand(codec),
	)...)
	return rootTransactionCommand
}

func GetCLIRootQueryCommand(codec *codec.Codec) *cobra.Command {
	rootQueryCommand := &cobra.Command{
		Use:                        QuerierRoute,
		Short:                      "Share root query command.",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	rootQueryCommand.AddCommand()
	return rootQueryCommand
}
