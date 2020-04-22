package contract

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/persistenceOne/persistenceSDK/modules/contract/mapper"
	"github.com/persistenceOne/persistenceSDK/modules/contract/transactions/bid"
	"github.com/persistenceOne/persistenceSDK/modules/contract/transactions/sign"
)

type Keeper interface {
	getBidKeeper() bid.Keeper
	getSignKeeper() sign.Keeper
}

type baseKeeper struct {
	bidKeeper  bid.Keeper
	signKeeper sign.Keeper
}

func NewKeeper(codec *codec.Codec, storeKey sdkTypes.StoreKey, paramSpace params.Subspace) Keeper {
	Mapper := mapper.NewMapper(codec, storeKey)
	return baseKeeper{
		bidKeeper:  bid.NewKeeper(Mapper),
		signKeeper: sign.NewKeeper(Mapper),
	}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) getBidKeeper() bid.Keeper   { return baseKeeper.bidKeeper }
func (baseKeeper baseKeeper) getSignKeeper() sign.Keeper { return baseKeeper.signKeeper }
