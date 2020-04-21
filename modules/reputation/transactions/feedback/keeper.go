package feedback

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type Keeper interface {
	transact(Message) error
}

type baseKeeper struct {
}

func NewKeeper() Keeper {
	return baseKeeper{}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) transact(message Message) error {
	return nil
}
