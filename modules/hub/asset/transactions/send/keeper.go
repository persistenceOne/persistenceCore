package send

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/mapper"
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
	asset, error := baseKeeper.mapper.GetAsset(context, mapper.NewAssetAddress(message.Address))
	if error != nil {
		return error
	}
	baseKeeper.mapper.SetAsset(context, mapper.NewAsset(asset.GetAddress(), message.To))
	return nil
}
