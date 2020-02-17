package asset

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/constants"
)

const (
	DefaultCodespace sdkTypes.CodespaceType = constants.ModuleName

	IncorrectQueryCode sdkTypes.CodeType = 301
)

func IncorrectQueryError(errorMessage string) sdkTypes.Error {
	return sdkTypes.NewError(DefaultCodespace, IncorrectQueryCode, errorMessage)
}
