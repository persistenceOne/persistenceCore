package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	buildertypes "github.com/skip-mev/pob/x/builder/types"
)

func SetBuilderDefaultConfig() {
	buildertypes.DefaultEscrowAccountAddress = authtypes.NewModuleAddress(buildertypes.ModuleName)
	buildertypes.DefaultReserveFee = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1))
	buildertypes.DefaultMinBidIncrement = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1))
}
