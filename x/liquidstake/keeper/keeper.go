package keeper

import (
	"fmt"

	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	liquidkeeper "github.com/cosmos/gaia/v25/x/liquid/keeper"

	"github.com/persistenceOne/persistenceCore/v17/x/liquidstake/types"
)

// Keeper of the liquidstake store
type Keeper struct {
	cdc          codec.BinaryCodec
	storeService storetypes.KVStoreService

	accountKeeper  types.AccountKeeper
	bankKeeper     types.BankKeeper
	stakingKeeper  types.StakingKeeper
	mintKeeper     mintkeeper.Keeper
	distrKeeper    types.DistrKeeper
	slashingKeeper types.SlashingKeeper
	liquidKeeper   liquidkeeper.Keeper

	router    *baseapp.MsgServiceRouter
	authority string
}

// NewKeeper returns a liquidstake keeper.
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	mintKeeper mintkeeper.Keeper,
	distrKeeper types.DistrKeeper,
	slashingKeeper types.SlashingKeeper,
	liquidKeeper liquidkeeper.Keeper,
	router *baseapp.MsgServiceRouter,
	authority string,
) Keeper {
	// ensure liquidstake module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:            cdc,
		storeService:   storeService,
		accountKeeper:  accountKeeper,
		bankKeeper:     bankKeeper,
		stakingKeeper:  stakingKeeper,
		mintKeeper:     mintKeeper,
		distrKeeper:    distrKeeper,
		slashingKeeper: slashingKeeper,
		liquidKeeper:   liquidKeeper,
		router:         router,
		authority:      authority,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// SetParams sets the auth module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) error {
	if err := params.Validate(); err != nil {
		return err
	}

	store := k.storeService.OpenKVStore(ctx)
	bz := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, bz)

	return nil
}

// GetParams gets the auth module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (types.Params, error) {
	store := k.storeService.OpenKVStore(ctx)
	var params types.Params

	bz, err := store.Get(types.ParamsKey)
	if bz == nil || err != nil {
		return params, err
	}
	k.cdc.MustUnmarshal(bz, &params)
	return params, nil
}

// GetCodec return codec.Codec object used by the keeper
func (k Keeper) GetCodec() codec.BinaryCodec { return k.cdc }

// Router returns the keeper's msg router
func (k Keeper) Router() *baseapp.MsgServiceRouter {
	return k.router
}
