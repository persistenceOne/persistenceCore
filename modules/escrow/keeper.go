package escrow

import (
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/persistenceOne/persistenceSDK/modules/escrow/transactions/execute"
)

type Keeper interface {
	getExecuteKeeper() execute.Keeper
}

type baseKeeper struct {
	executeKeeper execute.Keeper
}

func NewKeeper(paramSpace params.Subspace) Keeper {
	return baseKeeper{
		executeKeeper: execute.NewKeeper(),
	}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) getExecuteKeeper() execute.Keeper { return baseKeeper.executeKeeper }
