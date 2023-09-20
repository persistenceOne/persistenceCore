package v9_2_0

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/persistenceOne/persistenceCore/v9/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade handler")

		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}

func Fork(ctx sdk.Context, keepers *stakingkeeper.Keeper) {

	params := keepers.GetParams(ctx)
	maxValidators := params.MaxValidators

	alldels := keepers.GetAllDelegations(ctx)

	valShareMap := make(map[string]stakingtypes.Validator)
	for _, del := range alldels {
		validator, ok := valShareMap[del.ValidatorAddress]
		if !ok {
			validator, found := keepers.GetValidator(ctx, del.GetValidatorAddr())
			if !found {
				panic("Validator not found" + del.ValidatorAddress)
			}
			validator.Tokens = sdk.ZeroInt()
			validator.DelegatorShares = sdk.ZeroDec()
			validator.ValidatorBondShares = sdk.ZeroDec()
			//validator.LiquidShares = sdk.ZeroDec() //Will be refreshed
			valShareMap[del.ValidatorAddress] = validator
		}
		validator.DelegatorShares = validator.DelegatorShares.Add(del.Shares)
		if del.ValidatorBond {
			validator.ValidatorBondShares = validator.ValidatorBondShares.Add(del.Shares)
		}
		valShareMap[del.ValidatorAddress] = validator
		// Tokens we do directly in the next.
	}
	err := keepers.RefreshTotalLiquidStaked(ctx)
	if err != nil {
		panic(err)
	}
	allvals := keepers.GetAllValidators(ctx)
	for _, val := range allvals {
		calculatedVal, ok := valShareMap[val.OperatorAddress]
		if !ok {
			panic("validator not found" + val.OperatorAddress)
		}
		if !val.DelegatorShares.Equal(calculatedVal.DelegatorShares) {
			// SHOW ME
			ctx.Logger().Info(fmt.Sprintf("Validator %s is affected", val.OperatorAddress))
			tokens := val.TokensFromShares(calculatedVal.DelegatorShares)
			calculatedVal.Tokens = tokens.TruncateInt()
			valShareMap[calculatedVal.OperatorAddress] = calculatedVal

			keepers.SetValidator(ctx, calculatedVal)
			// Fix voting power
			fixPower(ctx, keepers, val, calculatedVal, maxValidators)
		} else {
			ctx.Logger().Info(fmt.Sprintf("Validator %s is ok", val.OperatorAddress))
		}
	}

}

func fixPower(ctx sdk.Context, k *stakingkeeper.Keeper, oldval, newval stakingtypes.Validator, maxValidators uint32) {

	iterator := k.ValidatorsPowerStoreIterator(ctx)

	keys := [][]byte{}
	values := []sdk.ValAddress{}
	for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {
		// everything that is iterated in this loop is becoming or already a
		// part of the bonded validator set
		valAddr := sdk.ValAddress(iterator.Value())
		if newval.GetOperator().Equals(valAddr) {
			keys = append(keys, iterator.Key())
			values = append(values, valAddr)
		}
	}
	iterator.Close()

	for i, key := range keys {
		validator, found := k.GetValidator(ctx, values[i])
		k.DeleteValidatorByPowerIndexUsingKey(ctx, key) //no op if key already deleted
		k.DeleteValidatorByPowerIndex(ctx, validator)   // no op if validator not found

		if !found {
			panic("validator not foudn while changing power" + values[i].String())
		}
		//create using new key
		k.SetValidatorByPowerIndex(ctx, validator) // re-add
	}

}
