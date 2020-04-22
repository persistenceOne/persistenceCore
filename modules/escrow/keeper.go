package escrow

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/persistenceOne/persistenceSDK/modules/escrow/mapper"
	"github.com/persistenceOne/persistenceSDK/modules/escrow/transactions/execute"
)

type Keeper interface {
	getExecuteKeeper() execute.Keeper
}

type baseKeeper struct {
	executeKeeper execute.Keeper
}

func NewKeeper(codec *codec.Codec, storeKey sdkTypes.StoreKey, paramSpace params.Subspace) Keeper {
	Mapper := mapper.NewMapper(codec, storeKey)
	return baseKeeper{
		executeKeeper: execute.NewKeeper(Mapper),
	}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) getExecuteKeeper() execute.Keeper { return baseKeeper.executeKeeper }
