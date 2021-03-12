/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

// RegisterRoutes registers halving module REST handlers on the provided router.
func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
}
