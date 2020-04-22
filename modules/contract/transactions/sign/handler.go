package sign

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/contract/constants"
)

func HandleMessage(context sdkTypes.Context, keeper Keeper, message Message) (*sdkTypes.Result, error) {

	if error := keeper.transact(context, message); error != nil {
		return nil, error
	}

	context.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, constants.AttributeValueCategory),
		),
	)

	return &sdkTypes.Result{Events: context.EventManager().Events()}, nil
}
