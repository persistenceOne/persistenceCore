package lock

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/asset/mapper"
)

type Keeper interface {
	transact(sdkTypes.Context, Message) sdkTypes.Error
}

type baseKeeper struct {
	mapper mapper.Mapper
}

func NewKeeper(mapper mapper.Mapper) Keeper {
	return baseKeeper{mapper: mapper}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) transact(context sdkTypes.Context, message Message) sdkTypes.Error {
	asset, err := baseKeeper.mapper.Read(context, mapper.NewAssetAddress(message.Address))
	if err != nil {
		return err
	}
	asset.SetLock(message.Lock)
	baseKeeper.mapper.Update(context, asset)
	return nil
}
