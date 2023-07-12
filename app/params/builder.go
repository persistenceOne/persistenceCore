package params

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	buildertypes "github.com/skip-mev/pob/x/builder/types"
)

func SetBuilderDefaultConfig() {
	sdk.SetAddrCacheEnabled(false)
	defer sdk.SetAddrCacheEnabled(true)

	buildertypes.DefaultEscrowAccountAddress = authtypes.NewModuleAddress(buildertypes.ModuleName).String()
	buildertypes.DefaultReserveFee = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1))
	buildertypes.DefaultMinBidIncrement = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1))
}
