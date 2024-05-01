package v11_11_0_rc0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/persistenceOne/persistenceCore/v11/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running module migrations...")

		err := args.Keepers.StakingKeeper.RefreshTotalLiquidStaked(ctx)
		if err != nil {
			ctx.Logger().Error("LSM failed to refresh total liquid staked", "error", err)
		}

		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}
