package asset

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/mapper"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/burn"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/lock"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/send"
)

type Keeper interface {
	getBurnKeeper() burn.Keeper
	getLockKeeper() lock.Keeper
	getMintKeeper() mint.Keeper
	getSendKeeper() send.Keeper
}

type baseKeeper struct {
	burnKeeper burn.Keeper
	lockKeeper lock.Keeper
	mintKeeper mint.Keeper
	sendKeeper send.Keeper
}

func NewKeeper(codec *codec.Codec, storeKey sdkTypes.StoreKey, paramSpace params.Subspace) Keeper {
	mapper := mapper.NewMapper(codec, storeKey)
	return baseKeeper{
		burnKeeper: burn.NewKeeper(mapper),
		lockKeeper: lock.NewKeeper(mapper),
		mintKeeper: mint.NewKeeper(mapper),
		sendKeeper: send.NewKeeper(mapper),
	}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) getBurnKeeper() burn.Keeper { return baseKeeper.burnKeeper }
func (baseKeeper baseKeeper) getLockKeeper() lock.Keeper { return baseKeeper.lockKeeper }
func (baseKeeper baseKeeper) getMintKeeper() mint.Keeper { return baseKeeper.mintKeeper }
func (baseKeeper baseKeeper) getSendKeeper() send.Keeper { return baseKeeper.sendKeeper }
