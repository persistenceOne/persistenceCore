package v11_5_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/persistenceOne/persistenceCore/v12/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade handler")

		// run x/liquidstake params upgrade to set the defaults in the state
		lsParams := args.Keepers.LiquidStakeKeeper.GetParams(ctx)
		lsParams.ModulePaused = false
		lsParams.WhitelistAdminAddress = "persistence1my083l0vkug8p4yvajpa3dnpf2pwmzmym35lne"
		if err := args.Keepers.LiquidStakeKeeper.SetParams(ctx, lsParams); err != nil {
			panic(err)
		}

		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}
