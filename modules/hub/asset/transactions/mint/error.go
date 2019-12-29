package mint

import (
	"github.com/commitHub/commitBlockchain/modules/hub/asset/constants"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdkTypes.CodespaceType = constants.ModuleName

	IncorrectMessageCode sdkTypes.CodeType = 101
)

func IncorrectMessageError(errorMessage string) sdkTypes.Error {
	return sdkTypes.NewError(DefaultCodespace, IncorrectMessageCode, errorMessage)
}
