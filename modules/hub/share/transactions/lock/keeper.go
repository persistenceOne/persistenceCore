package lock

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type Keeper interface {
	transact(Message) sdkTypes.Error
}

type BaseKeeper struct {
	codespace sdkTypes.CodespaceType
}

func NewBaseKeeper() BaseKeeper {
	return BaseKeeper{}
}

var _ Keeper = (*BaseKeeper)(nil)

func (baseKeeper BaseKeeper) transact(message Message) sdkTypes.Error {
	return nil
}
