package constants

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

const (
	validator1 = "cosmosvaloper1d20g0gcwhrwv8f2626dx0nkhauu0rsqkz4axrg"
)

var (
	Denom         = "stake"
	Validator1, _ = sdkTypes.ValAddressFromBech32(validator1)
	Address       = sdkTypes.AccAddress{}
)
