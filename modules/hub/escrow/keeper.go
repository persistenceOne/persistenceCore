package escrow

import (
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/persistenceOne/persistenceSDK/modules/hub/escrow/transactions/execute"
)

type Keeper interface {
	execute.Keeper
}

type BaseKeeper struct {
	execute.BaseKeeper
}

func NewKeeper(paramSpace params.Subspace) Keeper {
	return BaseKeeper{}
}

var _ Keeper = (*BaseKeeper)(nil)
