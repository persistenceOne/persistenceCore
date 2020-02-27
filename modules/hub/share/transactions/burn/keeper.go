package burn

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/mapper"
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
	baseKeeper.mapper.Delete(context, mapper.NewShareAddress(message.Address))
	return nil
}
