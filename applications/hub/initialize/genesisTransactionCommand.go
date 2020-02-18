package initialize

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/cosmos/cosmos-sdk/x/genutil/types"
)

func GenesisTransactionCommand(
	serverContext *server.Context,
	codec *codec.Codec,
	moduleBasicManager module.BasicManager,
	stakingMessageBuildingHelpers cli.StakingMsgBuildingHelpers,
	genesisAccountsIterator types.GenesisAccountsIterator,
	defaultNodeHome string,
	defaultClientHome string,
) *cobra.Command {
	return cli.GenTxCmd(
		serverContext,
		codec,
		moduleBasicManager,
		stakingMessageBuildingHelpers,
		genesisAccountsIterator,
		defaultNodeHome,
		defaultClientHome,
	)
}
