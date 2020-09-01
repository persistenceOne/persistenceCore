/*
 Copyright [2019] - [2020], PERSISTENCE TECHNOLOGIES PTE. LTD. and the assetMantle contributors
 SPDX-License-Identifier: Apache-2.0
*/

package initialize

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/types/module"
)

func GenesisTransactionCommand(
	serverContext *server.Context,
	codec *codec.Codec,
	moduleBasicManager module.BasicManager,
	stakingMessageBuildingHelpers cli.StakingMsgBuildingHelpers,
	genesisBalancesIterator auth.GenesisAccountIterator,
	defaultNodeHome string,
	defaultClientHome string,
) *cobra.Command {
	return cli.GenTxCmd(
		serverContext,
		codec,
		moduleBasicManager,
		stakingMessageBuildingHelpers,
		genesisBalancesIterator,
		defaultNodeHome,
		defaultClientHome,
	)
}
