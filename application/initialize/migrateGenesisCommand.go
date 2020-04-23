package initialize

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
)

func MigrateGenesisCommand(
	serverContext *server.Context,
	codec *codec.Codec,
) *cobra.Command {
	return cli.MigrateGenesisCmd(
		serverContext,
		codec,
	)
}
