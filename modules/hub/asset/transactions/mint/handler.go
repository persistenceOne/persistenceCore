package mint

import (
	"github.com/commitHub/commitBlockchain/modules/hub/asset/constants"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func HandleMessage(ctx sdkTypes.Context, keeper Keeper, message Message) sdkTypes.Result {

	if error := keeper.Mint(); error != nil {
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
