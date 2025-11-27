package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

const (
	// ModuleName is the name of the liquidstake module
	ModuleName = "liquidstake"

	// RouterKey is the message router key for the liquidstake module
	RouterKey = ModuleName

	// StoreKey is the default store key for the liquidstake module
	// To avoid collision with liquidstakeibc we make it xprtliquidstake
	StoreKey = "xprt" + ModuleName

	// QuerierRoute is the querier route for the liquidstake module
	QuerierRoute = ModuleName

	// Epoch identifiers
	AutocompoundEpoch = "hour"
	RebalanceEpoch    = "day"
)

var (
	ParamsKey = []byte{0x01}

	// LiquidValidatorsKey defines prefix for each key to a liquid validator
	LiquidValidatorsKey = []byte{0x02}
)

// GetLiquidValidatorKey creates the key for the liquid validator with address
// VALUE: liquidstake/LiquidValidator
func GetLiquidValidatorKey(operatorAddr sdk.ValAddress) []byte {
	tmp := append([]byte{}, LiquidValidatorsKey...)
	return append(tmp, address.MustLengthPrefix(operatorAddr)...)
}
