package v14_0_0_rc0

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	liquidkeeper "github.com/cosmos/gaia/v24/x/liquid/keeper"
	liquidtypes "github.com/cosmos/gaia/v24/x/liquid/types"
	"github.com/persistenceOne/persistenceCore/v14/app/keepers"
	"github.com/persistenceOne/persistenceCore/v14/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.Logger().Info("Running upgrade handler")
		err := FixLSMData(sdkCtx, args)
		if err != nil {
			return vm, err
		}
		sdkCtx.Logger().Info("running module migrations...")
		vm, err = args.ModuleManager.RunMigrations(sdkCtx, args.Configurator, vm)
		if err != nil {
			return vm, err
		}

		err = MigrateLSMState(sdkCtx, args.Keepers)
		if err != nil {
			return vm, errorsmod.Wrapf(err, "migrating LSM state to x/liquid")
		}

		sdkCtx.Logger().Info("Upgrade complete")
		return vm, nil
	}
}

func FixLSMData(ctx context.Context, args upgrades.UpgradeHandlerArgs) error {
	k := args.Keepers.StakingKeeper
	err := k.RefreshTotalLiquidStaked(ctx)
	if err != nil {
		return err
	}
	vals, err := k.GetAllValidators(ctx)
	if err != nil {
		return err
	}
	for _, validator := range vals {
		validator.ValidatorBondShares = math.LegacyZeroDec()
		err := k.SetValidator(ctx, validator)
		if err != nil {
			return err
		}
	}

	// Sum up the total liquid tokens and increment each validator's liquid shares
	dels, err := k.GetAllDelegations(ctx)
	for _, delegation := range dels {
		// disable valbond on delegator level.
		if delegation.ValidatorBond {
			delegation.ValidatorBond = false
			err = k.SetDelegation(ctx, delegation)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// taken from https://github.com/cosmos/gaia/tree/v24.0.0/app/upgrades/v24#L44C1-L163C2
func MigrateLSMState(ctx sdk.Context, keepers *keepers.AppKeepers) error {
	sk := keepers.StakingKeeper
	lsmk := keepers.LiquidKeeper

	err := migrateParams(ctx, sk, lsmk)
	if err != nil {
		return fmt.Errorf("error migrating params: %w", err)
	}

	err = migrateTokenizeShares(ctx, sk, lsmk)
	if err != nil {
		return fmt.Errorf("error migrating tokenize records: %w", err)
	}

	migrateLastTokenizeShareRecordID(ctx, sk, lsmk)
	migrateTokenizeShareLocks(ctx, sk, lsmk)

	return nil
}

func migrateParams(ctx sdk.Context, sk *stakingkeeper.Keeper, lsmk *liquidkeeper.Keeper) error {
	stakingParams, err := sk.GetParams(ctx)
	if err != nil {
		return err
	}

	liquidParams, err := lsmk.GetParams(ctx)
	if err != nil {
		return err
	}

	liquidParams.GlobalLiquidStakingCap = stakingParams.GlobalLiquidStakingCap
	liquidParams.ValidatorLiquidStakingCap = stakingParams.ValidatorLiquidStakingCap

	return lsmk.SetParams(ctx, liquidParams)
}

func migrateTokenizeShares(ctx sdk.Context, sk *stakingkeeper.Keeper, lsmk *liquidkeeper.Keeper) error {
	totalLiquidStaked := math.ZeroInt()
	liquidValidators := make(map[string]liquidtypes.LiquidValidator)

	tokenizeShareRecords := sk.GetAllTokenizeShareRecords(ctx)
	for _, record := range tokenizeShareRecords {
		lsmRecord := liquidtypes.TokenizeShareRecord{
			Id:            record.Id,
			Owner:         record.Owner,
			ModuleAccount: record.ModuleAccount,
			Validator:     record.Validator,
		}
		if err := lsmk.AddTokenizeShareRecord(ctx, lsmRecord); err != nil {
			return err
		}

		valAddress, err := sdk.ValAddressFromBech32(record.Validator)
		if err != nil {
			return fmt.Errorf("invalid validator address: %w", err)
		}

		validator, err := sk.GetValidator(ctx, valAddress)
		if err != nil {
			return fmt.Errorf("invalid validator address: %w", err)
		}

		delegation, err := sk.GetDelegation(ctx, record.GetModuleAddress(), valAddress)
		if err != nil {
			return fmt.Errorf("unable to get delegation: %w", err)
		}

		liquidVal, found := liquidValidators[record.Validator]
		if !found {
			liquidValidators[record.Validator] = liquidtypes.NewLiquidValidator(validator.OperatorAddress)
			liquidVal = liquidValidators[record.Validator]
		}

		liquidStatedTokensInDelegation := validator.TokensFromShares(delegation.Shares).TruncateInt()
		liquidVal.LiquidShares = liquidVal.LiquidShares.Add(delegation.Shares)
		liquidValidators[record.Validator] = liquidVal
		totalLiquidStaked = totalLiquidStaked.Add(liquidStatedTokensInDelegation)
	}

	lsmk.SetTotalLiquidStakedTokens(ctx, totalLiquidStaked)

	// Also add zero-ed out liquid vals
	allVals, err := sk.GetAllValidators(ctx)
	if err != nil {
		return fmt.Errorf("unable to get all validators: %w", err)
	}
	for _, val := range allVals {
		if _, ok := liquidValidators[val.OperatorAddress]; !ok {
			liquidValidators[val.OperatorAddress] = liquidtypes.NewLiquidValidator(val.OperatorAddress)
		}
	}

	for _, liquidVal := range liquidValidators {
		if err := lsmk.SetLiquidValidator(ctx, liquidVal); err != nil {
			return fmt.Errorf("error migrating liquid validator: %w", err)
		}
	}

	return nil
}

func migrateLastTokenizeShareRecordID(ctx sdk.Context, sk *stakingkeeper.Keeper, lsmk *liquidkeeper.Keeper) {
	lastTokenizeShareRecordID := sk.GetLastTokenizeShareRecordID(ctx)
	lsmk.SetLastTokenizeShareRecordID(ctx, lastTokenizeShareRecordID)
}

func migrateTokenizeShareLocks(ctx sdk.Context, sk *stakingkeeper.Keeper, lsmk *liquidkeeper.Keeper) {
	tokenizeShareLocks := sk.GetAllTokenizeSharesLocks(ctx)
	converted := make([]liquidtypes.TokenizeShareLock, len(tokenizeShareLocks))
	for i, tokenizeShareLock := range tokenizeShareLocks {
		converted[i] = liquidtypes.TokenizeShareLock{
			Address:        tokenizeShareLock.Address,
			Status:         tokenizeShareLock.Status,
			CompletionTime: tokenizeShareLock.CompletionTime,
		}
	}
	lsmk.SetTokenizeShareLocks(ctx, converted)
}
