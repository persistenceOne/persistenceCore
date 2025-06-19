package v11_12_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	liquidstaketypes "github.com/persistenceOne/pstake-native/v3/x/liquidstake/types"

	"github.com/persistenceOne/persistenceCore/v12/app/keepers"
	"github.com/persistenceOne/persistenceCore/v12/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running module migrations...")

		DelegateLiquidStakeRewards(ctx, args.Keepers)

		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}

func DelegateLiquidStakeRewards(ctx sdk.Context, k *keepers.AppKeepers) {
	proxyAccBalance := k.LiquidStakeKeeper.GetProxyAccBalance(ctx, liquidstaketypes.LiquidStakeProxyAcc)
	amountToDelegate := proxyAccBalance.Amount

	whitelistedValidators := k.LiquidStakeKeeper.GetParams(ctx).WhitelistedValidators
	whitelistedValsMap := liquidstaketypes.GetWhitelistedValsMap(whitelistedValidators)
	activeLiquidVals := k.LiquidStakeKeeper.GetActiveLiquidValidators(ctx, whitelistedValsMap)

	// currently auto compounding fee rate is zero, therefore fee logic is not added here.

	err := k.LiquidStakeKeeper.LiquidDelegate(ctx, liquidstaketypes.LiquidStakeProxyAcc, activeLiquidVals, amountToDelegate, whitelistedValsMap)

	// if liquid delegate fails then try to delegate to any one active validator
	if err != nil {
		ctx.Logger().Info("failed to liquid delegate, trying to delegate to any one validator...", "error", err.Error())

		for _, lv := range activeLiquidVals {
			val, _ := k.StakingKeeper.GetValidator(ctx, lv.GetOperator())
			err2 := k.LiquidStakeKeeper.DelegateWithCap(ctx, liquidstaketypes.LiquidStakeProxyAcc, val, amountToDelegate)
			if err2 != nil {
				ctx.Logger().Info("failed to delegate", "validator", val.GetOperator(), "error", err2.Error())
				// continue with next val
			} else {
				// successfully delegated, exit
				return
			}
		}

		ctx.Logger().Error("failed to delegate to any of active val set")
	}
}
