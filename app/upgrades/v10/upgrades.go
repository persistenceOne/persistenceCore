package v10

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/persistenceOne/persistenceCore/v9/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade handler")

		consensusParams, err := args.Keepers.ConsensusParamsKeeper.Get(ctx)
		if err != nil {
			panic(err)
		}
		consensusParams.Block.MaxBytes = 5242880
		args.Keepers.ConsensusParamsKeeper.Set(ctx, consensusParams)

		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}
