package share

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/mapper"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/queries/share"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/burn"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/lock"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/send"
)

type Keeper interface {
	getBurnKeeper() burn.Keeper
	getLockKeeper() lock.Keeper
	getMintKeeper() mint.Keeper
	getSendKeeper() send.Keeper
	getShareQuerier() share.Querier
}

type baseKeeper struct {
	burnKeeper   burn.Keeper
	lockKeeper   lock.Keeper
	mintKeeper   mint.Keeper
	sendKeeper   send.Keeper
	shareQuerier share.Querier
}

func NewKeeper(codec *codec.Codec, storeKey sdkTypes.StoreKey, paramSpace params.Subspace) Keeper {
	Mapper := mapper.NewMapper(codec, storeKey)
	return baseKeeper{
		burnKeeper:   burn.NewKeeper(Mapper),
		lockKeeper:   lock.NewKeeper(Mapper),
		mintKeeper:   mint.NewKeeper(Mapper),
		sendKeeper:   send.NewKeeper(Mapper),
		shareQuerier: share.NewQuerier(Mapper),
	}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) getBurnKeeper() burn.Keeper     { return baseKeeper.burnKeeper }
func (baseKeeper baseKeeper) getLockKeeper() lock.Keeper     { return baseKeeper.lockKeeper }
func (baseKeeper baseKeeper) getMintKeeper() mint.Keeper     { return baseKeeper.mintKeeper }
func (baseKeeper baseKeeper) getSendKeeper() send.Keeper     { return baseKeeper.sendKeeper }
func (baseKeeper baseKeeper) getShareQuerier() share.Querier { return baseKeeper.shareQuerier }
