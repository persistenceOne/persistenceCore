package share

import (
	"fmt"

	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/mint"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdkTypes.Handler {
	return func(context sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		context = context.WithEventManager(sdkTypes.NewEventManager())

		switch message := msg.(type) {
		case mint.Message:
			return mint.HandleMessage(context, keeper, message)

		default:
			return sdkTypes.ErrUnknownRequest(fmt.Sprintf("Unknown share message type: %T", msg)).Result()
		}
	}
}
