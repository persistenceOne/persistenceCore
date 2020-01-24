package contract

import (
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/persistenceOne/persistenceSDK/modules/hub/contract/transactions/bid"
	"github.com/persistenceOne/persistenceSDK/modules/hub/contract/transactions/sign"
)

type Keeper interface {
	getBidKeeper() bid.Keeper
	getSignKeeper() sign.Keeper
}

type baseKeeper struct {
	bidKeeper  bid.Keeper
	signKeeper sign.Keeper
}

func NewKeeper(paramSpace params.Subspace) Keeper {
	return baseKeeper{
		bidKeeper:  bid.NewKeeper(),
		signKeeper: sign.NewKeeper(),
	}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) getBidKeeper() bid.Keeper   { return baseKeeper.bidKeeper }
func (baseKeeper baseKeeper) getSignKeeper() sign.Keeper { return baseKeeper.signKeeper }
