package initialize

import (
	"github.com/cosmos/cosmos-sdk/x/genaccounts/client/cli"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
)

func AddGenesisAccountCommand(
	serverContext *server.Context,
	codec *codec.Codec,
	defaultNodeHome string,
	defaultClientHome string,
) *cobra.Command {
	return cli.AddGenesisAccountCmd(
		serverContext,
		codec,
		defaultNodeHome,
		defaultClientHome,
	)
}
