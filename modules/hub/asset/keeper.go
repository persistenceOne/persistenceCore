package asset

import (
	"github.com/commitHub/commitBlockchain/modules/hub/asset/transactions/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
)

type Keeper interface {
	mint.Keeper
}

type BaseKeeper struct {
	mint.BaseKeeper
}

func NewKeeper(paramSpace params.Subspace) Keeper {
	return BaseKeeper{}
}

var _ Keeper = (*BaseKeeper)(nil)
