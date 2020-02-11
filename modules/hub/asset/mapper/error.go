package mapper

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/constants"
)

const (
	DefaultCodespace sdkTypes.CodespaceType = constants.ModuleName

	AssetNotFound sdkTypes.CodeType = 201
)

func AssetNotFoundError(errorMessage string) sdkTypes.Error {
	return sdkTypes.NewError(DefaultCodespace, AssetNotFound, errorMessage)
}
