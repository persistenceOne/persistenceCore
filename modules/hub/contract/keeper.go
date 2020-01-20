package contract

import (
	"github.com/persistenceOne/persistenceSDK/modules/hub/contract/transactions/sign"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type Keeper interface {
	sign.Keeper
}

type BaseKeeper struct {
	sign.BaseKeeper
}

func NewKeeper(paramSpace params.Subspace) Keeper {
	return BaseKeeper{}
}

var _ Keeper = (*BaseKeeper)(nil)
