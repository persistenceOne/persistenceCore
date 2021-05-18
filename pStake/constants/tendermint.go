package constants

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

const (
	validator1 = "cosmosvaloper1pkkayn066msg6kn33wnl5srhdt3tnu2v8fyhft"
)

var (
	Denom         = "uatom"
	Validator1, _ = sdkTypes.ValAddressFromBech32(validator1)
	Address       = sdkTypes.AccAddress{}
)
