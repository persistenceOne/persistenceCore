package mapper

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/asset/constants"
)

func assetNotFoundError(errorMessage string) sdkTypes.Error {
	return sdkTypes.NewError(constants.DefaultCodespace, constants.AssetNotFoundCode, errorMessage)
}
