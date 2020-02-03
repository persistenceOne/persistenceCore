package asset

import (
	"fmt"

	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/burn"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/lock"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/send"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdkTypes.Handler {
	return func(context sdkTypes.Context, msg sdkTypes.Msg) sdkTypes.Result {
		context = context.WithEventManager(sdkTypes.NewEventManager())

		switch message := msg.(type) {
		case burn.Message:
			return burn.HandleMessage(context, keeper.getBurnKeeper(), message)
		case lock.Message:
			return lock.HandleMessage(context, keeper.getLockKeeper(), message)
		case mint.Message:
			return mint.HandleMessage(context, keeper.getMintKeeper(), message)
		case send.Message:
			return send.HandleMessage(context, keeper.getSendKeeper(), message)

		default:
			return sdkTypes.ErrUnknownRequest(fmt.Sprintf("Unknown asset message type: %T", msg)).Result()
		}
	}
}
