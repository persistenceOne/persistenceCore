package lock

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/share/mapper"
)

type Keeper interface {
	transact(sdkTypes.Context, Message) error
}

type baseKeeper struct {
	mapper mapper.Mapper
}

func NewKeeper(mapper mapper.Mapper) Keeper {
	return baseKeeper{mapper: mapper}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) transact(context sdkTypes.Context, message Message) error {
	share, err := baseKeeper.mapper.Read(context, mapper.NewShareAddress(message.Address))
	if err != nil {
		return err
	}
	share.SetLock(message.Lock)
	baseKeeper.mapper.Update(context, share)
	return nil
}
