package asset

import (
	"fmt"

	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/mint"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdkTypes.Handler {
	return func(context sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		context = context.WithEventManager(sdkTypes.NewEventManager())

		switch message := msg.(type) {
		case mint.Message:
			return mint.HandleMessage(context, keeper, message)

		default:
			return sdkTypes.ErrUnknownRequest(fmt.Sprintf("Unknown asset message type: %T", msg)).Result()
		}
	}
}
