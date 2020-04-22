package escrow

import (
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/persistenceSDK/modules/escrow/constants"

	"github.com/persistenceOne/persistenceSDK/modules/escrow/transactions/execute"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(keeper Keeper) sdkTypes.Handler {
	return func(context sdkTypes.Context, msg sdkTypes.Msg) (*sdkTypes.Result, error) {
		context = context.WithEventManager(sdkTypes.NewEventManager())

		switch message := msg.(type) {
		case execute.Message:
			return execute.HandleMessage(context, keeper.getExecuteKeeper(), message)

		default:
			return nil, errors.Wrapf(constants.UnknownMessageCode, "%T", msg)
		}
	}
}
