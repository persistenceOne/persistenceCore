package send

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/constants"
)

func HandleMessage(ctx sdkTypes.Context, keeper Keeper, message Message) sdkTypes.Result {

	if error := keeper.transact(message); error != nil {
		return error.Result()
	}

	ctx.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, constants.AttributeValueCategory),
		),
	)

	return sdkTypes.Result{Events: ctx.EventManager().Events()}
}
