package burn

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type Keeper interface {
	transact(Message) sdkTypes.Error
}

type baseKeeper struct {
}

func NewKeeper() baseKeeper {
	return baseKeeper{}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) transact(message Message) sdkTypes.Error {
	return nil
}
