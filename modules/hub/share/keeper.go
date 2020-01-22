package share

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

type Keeper interface {
	// burn.Keeper
	// lock.Keeper
	// mint.Keeper
	// send.Keeper
}

type BaseKeeper struct {
	// burnBaseKeeper burn.BaseKeeper
	// lockBaseKeeper lock.BaseKeeper
	// mintBaseKeeper mint.BaseKeeper
	// sendBaseKeeper send.BaseKeeper
}

func NewKeeper(paramSpace params.Subspace) Keeper {
	return BaseKeeper{}
	// return BaseKeeper{
	// 	burnBaseKeeper: burn.NewBaseKeeper(),
	// 	lockBaseKeeper: lock.NewBaseKeeper(),
	// 	mintBaseKeeper: mint.NewBaseKeeper(),
	// 	sendBaseKeeper: send.NewBaseKeeper(),
	// }
}

var _ Keeper = (*BaseKeeper)(nil)
