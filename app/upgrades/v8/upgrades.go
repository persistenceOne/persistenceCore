package v8

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icaMigrations "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/migrations/v6"
	lscosmostypes "github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"

	"github.com/persistenceOne/persistenceCore/v7/app/keepers"
	"github.com/persistenceOne/persistenceCore/v7/app/upgrades"
)

func setInitialMinCommissionRate(ctx sdk.Context, keepers *keepers.AppKeepers) {
	stakingParams := keepers.StakingKeeper.GetParams(ctx)
	stakingParams.MinCommissionRate = sdk.NewDecWithPrec(5, 2)
	keepers.StakingKeeper.SetParams(ctx, stakingParams)
}

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		err := icaMigrations.MigrateICS27ChannelCapability(
			ctx,
			args.Codec,
			args.CapabilityStoreKey,
			args.CapabilityKeeper,
			lscosmostypes.ModuleName,
		)

		if err != nil {
			return nil, err
		}

		setInitialMinCommissionRate(ctx, args.Keepers)

		// force set validator (if mcr < 5)

		ctx.Logger().Info("start to run module migrations...")
		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}
