package v8_fix_invariant

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	// disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"

	"github.com/persistenceOne/persistenceCore/v8/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade handler")

		valAddr, err := sdk.ValAddressFromBech32(ValidatorAddress)
		if err != nil {
			return nil, err
		}

		val, found := args.Keepers.StakingKeeper.GetValidator(ctx, valAddr)
		if !found {
			return nil, fmt.Errorf("validator not found")
		}

		dels := args.Keepers.StakingKeeper.GetAllDelegations(ctx)
		for _, del := range dels {
			if del.ValidatorAddress == ValidatorAddress || del.DelegatorAddress == DelegatorAddress {
				args.Keepers.StakingKeeper.RemoveDelegation(ctx, del)
				validator, _ := args.Keepers.StakingKeeper.GetValidator(ctx, del.GetValidatorAddr())
				args.Keepers.StakingKeeper.RemoveValidatorTokensAndShares(ctx, validator, del.Shares)
			}
		}

		args.Keepers.StakingKeeper.RemoveValidatorTokens(ctx, val, val.Tokens)

		args.Keepers.BankKeeper.BurnCoins(ctx, stakingtypes.NotBondedPoolName, sdk.NewCoins(sdk.NewCoin(args.Keepers.StakingKeeper.GetParams(ctx).BondDenom, val.Tokens)))

		args.Keepers.StakingKeeper.RemoveValidator(ctx, valAddr)
		
		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}