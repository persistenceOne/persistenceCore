/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/mint"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/persistenceOne/persistenceCore/x/halving/internal/types"
)

// Keeper of the halving store
type Keeper struct {
	cdc        *codec.Codec
	storeKey   sdk.StoreKey
	paramSpace params.Subspace
	mintKeeper types.MintKeeper
}

// NewKeeper creates a new halving Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace,
	mintKeeper types.MintKeeper,
) Keeper {

	return Keeper{
		cdc:        cdc,
		storeKey:   key,
		paramSpace: paramSpace.WithKeyTable(types.ParamKeyTable()),
		mintKeeper: mintKeeper,
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
func (k Keeper) GetMintingParams(ctx sdk.Context) (params mint.Params) {
	return k.mintKeeper.GetParams(ctx)
}

// SetMintingParams sets the total set of halving parameters.
func (k Keeper) SetMintingParams(ctx sdk.Context, params mint.Params) {
	k.mintKeeper.SetParams(ctx, params)
}
