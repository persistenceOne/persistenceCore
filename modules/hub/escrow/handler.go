package escrow

import (
	"fmt"

	"github.com/persistenceOne/persistenceSDK/modules/hub/escrow/transactions/execute"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdkTypes.Handler {
	return func(context sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		context = context.WithEventManager(sdkTypes.NewEventManager())

		switch message := msg.(type) {
		case execute.Message:
			return execute.HandleMessage(context, keeper, message)

		default:
			return sdkTypes.ErrUnknownRequest(fmt.Sprintf("Unknown escrow message type: %T", msg)).Result()
		}
	}
}
