package initialize

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
)

func CollectGenesisTransactionsCommand(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cli.CollectGenTxsCmd(ctx, cdc)
}
