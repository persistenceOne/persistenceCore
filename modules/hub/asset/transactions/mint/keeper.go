package mint

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type Keeper interface {
	Mint() sdkTypes.Error
}

type BaseKeeper struct {
}

func NewKeeper() Keeper {
	return BaseKeeper{}
}

var _ Keeper = (*BaseKeeper)(nil)

func (baseKeeper BaseKeeper) Mint() sdkTypes.Error {
	return nil
}
