package v15_0_0

import (
	"context"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/persistenceOne/persistenceCore/v15/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.Logger().Info("Running upgrade handler")

		vm, err := args.ModuleManager.RunMigrations(sdkCtx, args.Configurator, vm)
		if err != nil {
			return vm, err
		}

		sdkCtx.Logger().Info("Upgrade complete")
		return vm, nil
	}
}
