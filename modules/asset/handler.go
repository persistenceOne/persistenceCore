package asset

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/persistenceSDK/modules/asset/constants"

	"github.com/persistenceOne/persistenceSDK/modules/asset/transactions/burn"
	"github.com/persistenceOne/persistenceSDK/modules/asset/transactions/lock"
	"github.com/persistenceOne/persistenceSDK/modules/asset/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/modules/asset/transactions/send"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdkTypes.Handler {
	return func(context sdkTypes.Context, msg sdkTypes.Msg) (*sdkTypes.Result, error) {
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
			return nil, errors.Wrapf(constants.UnknownMessageCode, "%T", msg)
		}
	}
}
