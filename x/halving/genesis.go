/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package halving

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis new halving genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	err := keeper.SetParams(ctx, data.Params)
	if err != nil {
		panic(err)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) *GenesisState {
	params, err := keeper.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	return NewGenesisState(params)
}
