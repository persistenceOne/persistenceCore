package asset

import (
	"github.com/persistenceOne/persistenceSDK/modules/asset/queries/asset"
	"github.com/persistenceOne/persistenceSDK/modules/asset/transactions/burn"
	"github.com/persistenceOne/persistenceSDK/modules/asset/transactions/lock"
	"github.com/persistenceOne/persistenceSDK/modules/asset/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/modules/asset/transactions/send"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
)

func GetCLIRootTransactionCommand(codec *codec.Codec) *cobra.Command {
	rootTransactionCommand := &cobra.Command{
		Use:                        TransactionRoute,
		Short:                      "Asset root transaction command.",
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
		Short:                      "Asset root query command.",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	rootQueryCommand.AddCommand(client.GetCommands(
		asset.QueryCommand(codec),
	)...)
	return rootQueryCommand
}
