package initialize

import (
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
)

func InitializeCommand(
	serverContext *server.Context,
	cdc *codec.Codec,
	moduleBasicManager module.BasicManager,
	defaultNodeHome string,
) *cobra.Command {
	return cli.InitCmd(
		serverContext,
		cdc,
		moduleBasicManager,
		defaultNodeHome,
	)
}
