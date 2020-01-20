package reputation

import (
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/persistenceOne/persistenceSDK/modules/hub/reputation/transactions/feedback"
)

type Keeper interface {
	feedback.Keeper
}

type BaseKeeper struct {
	feedback.BaseKeeper
}

func NewKeeper(paramSpace params.Subspace) Keeper {
	return BaseKeeper{}
}

var _ Keeper = (*BaseKeeper)(nil)
