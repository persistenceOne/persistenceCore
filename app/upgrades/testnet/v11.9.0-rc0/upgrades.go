package v11_9_0_rc0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"

	"github.com/persistenceOne/persistenceCore/v11/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running module migrations...")

		// register param subspace in order to migrate from this legacy params.
		subspace, ok := args.Keepers.ParamsKeeper.GetSubspace(packetforwardtypes.ModuleName)
		if ok && !subspace.HasKeyTable() {
			subspace.WithKeyTable(packetforwardtypes.ParamKeyTable())
		}

		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}
