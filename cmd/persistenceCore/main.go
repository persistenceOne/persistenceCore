/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"os"

	servercmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/persistenceOne/persistenceCore/v15/app"
	"github.com/persistenceOne/persistenceCore/v15/cmd/persistenceCore/cmd"
)

func main() {

	rootCmd, _ := cmd.NewRootCmd()

	if err := servercmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
