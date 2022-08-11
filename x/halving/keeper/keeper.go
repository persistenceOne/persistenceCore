/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package keeper

import (
	"fmt"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceCore/v3/x/halving/types"
)

// Keeper of the halving store
type Keeper struct {
	storeKey   sdk.StoreKey
	paramSpace paramstypes.Subspace
	mintkeeper types.MintKeeper
}

// NewKeeper creates a new halving Keeper instance
func NewKeeper(
	key sdk.StoreKey, paramSpace paramstypes.Subspace,
	mintkeeper types.MintKeeper,
) Keeper {

	return Keeper{
		storeKey:   key,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
		mintkeeper: mintkeeper,
	}
}

//______________________________________________________________________

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

//______________________________________________________________________

// GetParams returns the total set of parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

//______________________________________________________________________

// GetMintingParams returns the total set of halving parameters.
func (k Keeper) GetMintingParams(ctx sdk.Context) (params minttypes.Params) {
	return k.mintkeeper.GetParams(ctx)
}

// SetMintingParams sets the total set of halving parameters.
func (k Keeper) SetMintingParams(ctx sdk.Context, params minttypes.Params) {
	k.mintkeeper.SetParams(ctx, params)
}
