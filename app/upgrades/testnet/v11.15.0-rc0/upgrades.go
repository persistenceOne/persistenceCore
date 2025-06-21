package v11_15_0_rc0

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/persistenceOne/persistenceCore/v12/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running module migrations...")
		err := MigrateVestingAccounts(ctx, args)
		if err != nil {
			return vm, err
		}
		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}

func MigrateVestingAccounts(ctx sdk.Context, args upgrades.UpgradeHandlerArgs) error {
	accounts := args.Keepers.AccountKeeper.GetAllAccounts(ctx)
	for _, account := range accounts {
		switch account.(type) {
		case *vestingtypes.PeriodicVestingAccount:
			a, ok := account.(*vestingtypes.PeriodicVestingAccount)
			if !ok {
				return errors.Wrapf(sdkerrors.ErrInvalidType, "invalid account type: %T", account)
			}
			args.Keepers.AccountKeeper.SetAccount(ctx, a.BaseAccount)
		case *vestingtypes.ContinuousVestingAccount:
			a, ok := account.(*vestingtypes.ContinuousVestingAccount)
			if !ok {
				return errors.Wrapf(sdkerrors.ErrInvalidType, "invalid account type: %T", account)
			}
			args.Keepers.AccountKeeper.SetAccount(ctx, a.BaseAccount)
		case *vestingtypes.DelayedVestingAccount:
			a, ok := account.(*vestingtypes.DelayedVestingAccount)
			if !ok {
				return errors.Wrapf(sdkerrors.ErrInvalidType, "invalid account type: %T", account)
			}
			args.Keepers.AccountKeeper.SetAccount(ctx, a.BaseAccount)
		default:
		}
	}
	return nil
}
