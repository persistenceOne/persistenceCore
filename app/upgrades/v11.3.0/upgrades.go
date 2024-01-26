package v11_3_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/persistenceOne/persistenceCore/v11/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade handler")

		for _, hc := range args.Keepers.LiquidStakeIBCKeeper.GetAllHostChains(ctx) {
			switch hc.ChainId {
			case "cosmoshub-4":
				upperLimit, _ := sdk.NewDecFromStr("1.01")
				lowerLimit, _ := sdk.NewDecFromStr("0.80")
				hc.Params.UpperCValueLimit = upperLimit
				hc.Params.LowerCValueLimit = lowerLimit
				args.Keepers.LiquidStakeIBCKeeper.SetHostChain(ctx, hc)

			case "osmosis-1":
				upperLimit, _ := sdk.NewDecFromStr("1.01")
				lowerLimit, _ := sdk.NewDecFromStr("0.97")
				hc.Params.UpperCValueLimit = upperLimit
				hc.Params.LowerCValueLimit = lowerLimit
				args.Keepers.LiquidStakeIBCKeeper.SetHostChain(ctx, hc)

			case "theta-testnet-001":
				upperLimit, _ := sdk.NewDecFromStr("1.01")
				lowerLimit, _ := sdk.NewDecFromStr("0.9")
				hc.Params.UpperCValueLimit = upperLimit
				hc.Params.LowerCValueLimit = lowerLimit
				args.Keepers.LiquidStakeIBCKeeper.SetHostChain(ctx, hc)

			case "osmo-test-5":
				upperLimit, _ := sdk.NewDecFromStr("1.01")
				lowerLimit, _ := sdk.NewDecFromStr("0.95")
				hc.Params.UpperCValueLimit = upperLimit
				hc.Params.LowerCValueLimit = lowerLimit
				args.Keepers.LiquidStakeIBCKeeper.SetHostChain(ctx, hc)

			case "dydx-test-4":
				upperLimit, _ := sdk.NewDecFromStr("1.01")
				lowerLimit, _ := sdk.NewDecFromStr("0.95")
				hc.Params.UpperCValueLimit = upperLimit
				hc.Params.LowerCValueLimit = lowerLimit
				args.Keepers.LiquidStakeIBCKeeper.SetHostChain(ctx, hc)

			case "gaia-1":
				upperLimit, _ := sdk.NewDecFromStr("1.01")
				lowerLimit, _ := sdk.NewDecFromStr("0.95")
				hc.Params.UpperCValueLimit = upperLimit
				hc.Params.LowerCValueLimit = lowerLimit
				args.Keepers.LiquidStakeIBCKeeper.SetHostChain(ctx, hc)

			}
		}

		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}
