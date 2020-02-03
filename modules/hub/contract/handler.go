package contract

import (
	"fmt"

	"github.com/persistenceOne/persistenceSDK/modules/hub/contract/transactions/bid"
	"github.com/persistenceOne/persistenceSDK/modules/hub/contract/transactions/sign"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdkTypes.Handler {
	return func(context sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		context = context.WithEventManager(sdkTypes.NewEventManager())

		switch message := msg.(type) {
		case bid.Message:
			return bid.HandleMessage(context, keeper.getBidKeeper(), message)
		case sign.Message:
			return sign.HandleMessage(context, keeper.getSignKeeper(), message)

		default:
			return sdkTypes.ErrUnknownRequest(fmt.Sprintf("Unknown contract message type: %T", msg)).Result()
		}
	}
}
