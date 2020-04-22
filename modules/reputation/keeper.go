package reputation

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/persistenceOne/persistenceSDK/modules/reputation/mapper"
	"github.com/persistenceOne/persistenceSDK/modules/reputation/transactions/feedback"
)

type Keeper interface {
	getFeedbackKeeper() feedback.Keeper
}

type baseKeeper struct {
	feedbackKeeper feedback.Keeper
}

func NewKeeper(codec *codec.Codec, storeKey sdkTypes.StoreKey, paramSpace params.Subspace) Keeper {
	Mapper := mapper.NewMapper(codec, storeKey)
	return baseKeeper{
		feedbackKeeper: feedback.NewKeeper(Mapper),
	}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) getFeedbackKeeper() feedback.Keeper { return baseKeeper.feedbackKeeper }
