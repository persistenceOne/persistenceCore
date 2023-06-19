package v8

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icaMigrations "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/migrations/v6"
	oracletypes "github.com/persistenceOne/persistence-sdk/v2/x/oracle/types"
	lscosmostypes "github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"

	"github.com/persistenceOne/persistenceCore/v8/app/keepers"
	"github.com/persistenceOne/persistenceCore/v8/app/upgrades"
)

func setInitialMinCommissionRate(ctx sdk.Context, keepers *keepers.AppKeepers) {
	stakingParams := keepers.StakingKeeper.GetParams(ctx)
	stakingParams.MinCommissionRate = sdk.NewDecWithPrec(5, 2)
	keepers.StakingKeeper.SetParams(ctx, stakingParams)

	// Force update validator commission rate if it is lower than the minimum rate
	validators := keepers.StakingKeeper.GetAllValidators(ctx)
	for _, v := range validators {
		if v.Commission.Rate.LT(stakingParams.MinCommissionRate) {
			if v.Commission.MaxRate.LT(stakingParams.MinCommissionRate) {
				v.Commission.MaxRate = stakingParams.MinCommissionRate
			}

			v.Commission.Rate = stakingParams.MinCommissionRate
			v.Commission.UpdateTime = ctx.BlockHeader().Time

			// call the before-modification hook since we're about to update the commission
			if err := keepers.StakingKeeper.BeforeValidatorModified(ctx, v.GetOperator()); err != nil {
				ctx.Logger().Info(fmt.Sprintf("BeforeValidatorModified failed with: %s", err.Error()))
			}
			keepers.StakingKeeper.SetValidator(ctx, v)
		}
	}
}

func setOraclePairListEmpty(ctx sdk.Context, keepers *keepers.AppKeepers) {
	oracleParams := keepers.OracleKeeper.GetParams(ctx)
	oracleParams.AcceptList = oracletypes.DenomList{}
	keepers.OracleKeeper.SetParams(ctx, oracleParams)
}

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade handler")

		ctx.Logger().Info("migrating ics27 channel capability")
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

		ctx.Logger().Info("running module migrations")
		newVm, err := args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
		if err != nil {
			return newVm, err
		}

		ctx.Logger().Info("setting min commission rate to 5%")
		setInitialMinCommissionRate(ctx, args.Keepers)

		ctx.Logger().Info("setting acceptList to empty in oracle params")
		setOraclePairListEmpty(ctx, args.Keepers)

		return newVm, err
	}
}
