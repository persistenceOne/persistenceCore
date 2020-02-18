package initialize

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"
)

func CollectGenesisTransactionsCommand(
	serverContext *server.Context,
	codec *codec.Codec,
	genesisAccountsIterator types.GenesisAccountsIterator,
	defaultNodeHome string,
) *cobra.Command {
	return cli.CollectGenTxsCmd(
		serverContext,
		codec,
		genesisAccountsIterator,
		defaultNodeHome,
	)
}
