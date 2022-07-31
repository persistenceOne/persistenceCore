/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	serverCmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/persistenceOne/persistenceCore/v3/app"
	"github.com/persistenceOne/persistenceCore/v3/cmd/persistenceCore/cmd"
)

func main() {

	rootCmd, _ := cmd.NewRootCmd()

	if err := serverCmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		switch e := err.(type) {
		case server.ErrorCode:
			os.Exit(e.Code)

		default:
			os.Exit(1)
		}
	}
}
