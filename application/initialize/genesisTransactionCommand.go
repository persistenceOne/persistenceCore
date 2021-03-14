/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package initialize

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/types/module"
)

func GenesisTransactionCommand(
	moduleBasicManager module.BasicManager,
	stakingMessageBuildingHelpers client.TxEncodingConfig,
	genesisBalancesIterator types.GenesisBalancesIterator,
	defaultNodeHome string,
) *cobra.Command {
	return cli.GenTxCmd(
		moduleBasicManager,
		stakingMessageBuildingHelpers,
		genesisBalancesIterator,
		defaultNodeHome,
	)
}
