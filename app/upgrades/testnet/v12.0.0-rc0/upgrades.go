package v12_0_0_rc0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/persistenceOne/persistenceCore/v11/app/upgrades"
)

// This was buggy in mainnet upgrade -> potantial fix/ diff should have been -> https://github.com/persistenceOne/persistenceCore/pull/308
// PR #308 is not merged to preserve upgrade history as is.
func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade handler")

		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}
