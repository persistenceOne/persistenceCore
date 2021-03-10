/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package initialize

import (
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
)

func CollectGenesisTransactionsCommand(
	genesisBalancesIterator bankTypes.GenesisBalancesIterator,
	defaultNodeHome string,
) *cobra.Command {
	return cli.CollectGenTxsCmd(
		genesisBalancesIterator,
		defaultNodeHome,
	)
}
