package share

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/constants"
	abciTypes "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(keeper Keeper) sdkTypes.Querier {
	return func(context sdkTypes.Context, path []string, requestQuery abciTypes.RequestQuery) ([]byte, sdkTypes.Error) {
		switch path[0] {
		case constants.ShareQuery:
			return keeper.getShareQuerier().Query(context, requestQuery)

		default:
			return nil, sdkTypes.ErrUnknownRequest("unknown bank query endpoint")
		}
	}
}
