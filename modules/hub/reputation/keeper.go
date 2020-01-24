package reputation

import (
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/persistenceOne/persistenceSDK/modules/hub/reputation/transactions/feedback"
)

type Keeper interface {
	getFeedbackKeeper() feedback.Keeper
}

type baseKeeper struct {
	feedbackKeeper feedback.Keeper
}

func NewKeeper(paramSpace params.Subspace) Keeper {
	return baseKeeper{
		feedbackKeeper: feedback.NewKeeper(),
	}
}

var _ Keeper = (*baseKeeper)(nil)

func (baseKeeper baseKeeper) getFeedbackKeeper() feedback.Keeper { return baseKeeper.feedbackKeeper }
