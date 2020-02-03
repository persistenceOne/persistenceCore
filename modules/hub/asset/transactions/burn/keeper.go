package burn

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	mapper "github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions"
)

type Keeper interface {
	transact(Message) sdkTypes.Error
}

type baseKeeper struct {
	mapper mapper.Mapper
}

func NewKeeper(mapper mapper.Mapper) Keeper {
	return baseKeeper{mapper: mapper}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) transact(message Message) sdkTypes.Error {
	asset := baseKeeper.mapper.GetAsset(message.From)
	baseKeeper.mapper.SetAsset(asset)
	return nil
}
