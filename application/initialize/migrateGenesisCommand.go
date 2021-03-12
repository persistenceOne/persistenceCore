/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package initialize

import (
	"github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
)

func MigrateGenesisCommand() *cobra.Command { return cli.MigrateGenesisCmd() }
