/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	servercmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/monopauli/persistenceCore/app"
	"github.com/monopauli/persistenceCore/cmd/persistenceCore/cmd"
)

func main() {

	rootCmd, _ := cmd.NewRootCmd()

	if err := servercmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
