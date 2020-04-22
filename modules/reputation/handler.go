package reputation

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/persistenceSDK/modules/reputation/constants"

	"github.com/persistenceOne/persistenceSDK/modules/reputation/transactions/feedback"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdkTypes.Handler {
	return func(context sdkTypes.Context, msg sdkTypes.Msg) (*sdkTypes.Result, error) {
		context = context.WithEventManager(sdkTypes.NewEventManager())

		switch message := msg.(type) {
		case feedback.Message:
			return feedback.HandleMessage(context, keeper.getFeedbackKeeper(), message)

		default:
			return nil, errors.Wrapf(constants.UnknownMessageCode, "%T", msg)
		}
	}
}
