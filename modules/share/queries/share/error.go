package share

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/share/constants"
)

func incorrectQueryError(errorMessage string) sdkTypes.Error {
	return sdkTypes.NewError(constants.DefaultCodespace, constants.IncorrectQueryCode, errorMessage)
}
