package contract

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/persistenceSDK/modules/contract/constants"

	"github.com/persistenceOne/persistenceSDK/modules/contract/transactions/bid"
	"github.com/persistenceOne/persistenceSDK/modules/contract/transactions/sign"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdkTypes.Handler {
	return func(context sdkTypes.Context, msg sdkTypes.Msg) (*sdkTypes.Result, error) {
		context = context.WithEventManager(sdkTypes.NewEventManager())

		switch message := msg.(type) {
		case bid.Message:
			return bid.HandleMessage(context, keeper.getBidKeeper(), message)
		case sign.Message:
			return sign.HandleMessage(context, keeper.getSignKeeper(), message)

		default:
			return nil, errors.Wrapf(constants.UnknownMessageCode, "%T", msg)
		}
	}
}
