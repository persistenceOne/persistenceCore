package share

import (
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/burn"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/lock"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/send"
)

type Keeper interface {
	getBurnKeeper() burn.Keeper
	getLockKeeper() lock.Keeper
	getMintKeeper() mint.Keeper
	getSendKeeper() send.Keeper
}

type baseKeeper struct {
	burnKeeper burn.Keeper
	lockKeeper lock.Keeper
	mintKeeper mint.Keeper
	sendKeeper send.Keeper
}

func NewKeeper(paramSpace params.Subspace) Keeper {
	return baseKeeper{
		burnKeeper: burn.NewKeeper(),
		lockKeeper: lock.NewKeeper(),
		mintKeeper: mint.NewKeeper(),
		sendKeeper: send.NewKeeper(),
	}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) getBurnKeeper() burn.Keeper { return baseKeeper.burnKeeper }
func (baseKeeper baseKeeper) getLockKeeper() lock.Keeper { return baseKeeper.lockKeeper }
func (baseKeeper baseKeeper) getMintKeeper() mint.Keeper { return baseKeeper.mintKeeper }
func (baseKeeper baseKeeper) getSendKeeper() send.Keeper { return baseKeeper.sendKeeper }
