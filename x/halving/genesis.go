/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package halving

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/persistenceCore/v17/x/halving/keeper"
	"github.com/persistenceOne/persistenceCore/v17/x/halving/types"
)

// InitGenesis new halving genesis
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, data types.GenesisState) {
	err := keeper.SetParams(ctx, data.Params)
	if err != nil {
		panic(err)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	params, err := keeper.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	return types.NewGenesisState(params)
}
