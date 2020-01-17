package mint

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type Keeper interface {
	transact(Message) sdkTypes.Error
}

type BaseKeeper struct {
}

func NewKeeper() Keeper {
	return BaseKeeper{}
}

var _ Keeper = (*BaseKeeper)(nil)

func (baseKeeper BaseKeeper) transact(message Message) sdkTypes.Error {
	return nil
}
