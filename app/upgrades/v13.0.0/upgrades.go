package v13_0_0

import (
	"context"
	ibctmtypes "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/persistenceOne/persistenceCore/v13/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.Logger().Info("Running upgrade handler")
		sdkCtx.Logger().Info("running module migrations...")
		vm, err := args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
		if err != nil {
			return vm, err
		}
		sdkCtx.Logger().Info("Setting IBC Client AllowedClients")
		params := args.Keepers.IBCKeeper.ClientKeeper.GetParams(sdkCtx)
		params.AllowedClients = []string{ibctmtypes.ModuleName}
		args.Keepers.IBCKeeper.ClientKeeper.SetParams(sdkCtx, params)
		return vm, nil
	}
}
