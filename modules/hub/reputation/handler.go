package reputation

import (
	"fmt"

	"github.com/persistenceOne/persistenceSDK/modules/hub/reputation/transactions/feedback"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdkTypes.Handler {
	return func(context sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		context = context.WithEventManager(sdkTypes.NewEventManager())

		switch message := msg.(type) {
		case feedback.Message:
			return feedback.HandleMessage(context, keeper.getFeedbackKeeper(), message)

		default:
			return sdkTypes.ErrUnknownRequest(fmt.Sprintf("Unknown reputation message type: %T", msg)).Result()
		}
	}
}
