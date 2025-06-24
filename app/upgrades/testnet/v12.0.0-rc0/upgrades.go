package v12_0_0_rc0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/persistenceOne/persistenceCore/v12/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running module migrations...")
		err := FixLSMData(ctx, args)
		if err != nil {
			return vm, err
		}
		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}

func FixLSMData(ctx sdk.Context, args upgrades.UpgradeHandlerArgs) error {
	k := args.Keepers.StakingKeeper
	err := k.RefreshTotalLiquidStaked(ctx)
	if err != nil {
		return err
	}
	for _, validator := range k.GetAllValidators(ctx) {
		validator.ValidatorBondShares = sdk.ZeroDec()
		k.SetValidator(ctx, validator)
	}

	// Sum up the total liquid tokens and increment each validator's liquid shares
	for _, delegation := range k.GetAllDelegations(ctx) {
		delegatorAddress, err := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
		if err != nil {
			return err
		}

		// disable valbond on delegator level.
		if k.DelegatorIsLiquidStaker(delegatorAddress) {
			delegation.ValidatorBond = false
			k.SetDelegation(ctx, delegation)
		}
	}

	return nil
}
