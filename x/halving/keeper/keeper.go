/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/persistenceOne/persistenceCore/v16/x/halving/types"
)

type Keeper struct {
	storeKey   storetypes.KVStoreService
	paramSpace paramsTypes.Subspace
	mintKeeper mintkeeper.Keeper

	Schema      collections.Schema
	ParamsStore collections.Item[types.Params]
	authority   string
}

func NewKeeper(cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService, paramSpace paramsTypes.Subspace,
	mintKeeper mintkeeper.Keeper, ak authkeeper.AccountKeeper,
	authority string,
) Keeper {
	// ensure that authority is a valid AccAddress
	if _, err := ak.AddressCodec().StringToBytes(authority); err != nil {
		panic("authority is not a valid acc address")
	}
	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeKey:    storeService,
		paramSpace:  paramSpace.WithKeyTable(types.ParamKeyTable()),
		ParamsStore: collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		mintKeeper:  mintKeeper,
		authority:   authority,
	}
	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema
	return k
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
	return k.ParamsStore.Get(ctx)
}

func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	return k.ParamsStore.Set(ctx, params)
}

func (k Keeper) GetMintingParams(ctx sdk.Context) (mintTypes.Params, error) {
	return k.mintKeeper.Params.Get(ctx)
}

func (k Keeper) SetMintingParams(ctx sdk.Context, params mintTypes.Params) error {
	return k.mintKeeper.Params.Set(ctx, params)
}

// UpdateParams updates the halving module parameters
func (k Keeper) UpdateParams(ctx sdk.Context, authority string, params types.Params) error {
	if authority != k.GetAuthority() {
		return fmt.Errorf("unauthorized: authority %s is not the module authority", authority)
	}

	err := k.SetParams(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

// GetAuthority returns the module authority address
func (k Keeper) GetAuthority() string {
	return k.authority
}
