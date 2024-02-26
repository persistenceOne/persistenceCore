package v11_7_0

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/persistenceOne/persistenceCore/v11/app/keepers"
	"github.com/persistenceOne/persistenceCore/v11/app/upgrades"
	liquidstaketypes "github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade handler")

		// make sure this region runs only during CI and testnet v11.7.0 upgrade
		if chainID := ctx.ChainID(); chainID == "test-core-1" || chainID == "ictest-core-1" {
			if err := runLiquidstakeUpgradeMigration(ctx, args.Keepers); err != nil {
				panic(err)
			}
		} else {
			panic("chainID not expected: " + chainID)
		}

		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}

// runLiquidstakeUpgradeMigration contains stkXPRT related migrations
func runLiquidstakeUpgradeMigration(
	ctx sdk.Context,
	keepers *keepers.AppKeepers,
) error {
	oldProxyAcc := authtypes.NewModuleAddress(liquidstaketypes.ModuleName + "-LiquidStakeProxyAcc")

	delegations := []stakingtypes.Delegation{}
	keepers.StakingKeeper.IterateDelegatorDelegations(ctx, oldProxyAcc, func(delegation stakingtypes.Delegation) bool {
		delegations = append(delegations, delegation)
		return false
	})

	ctx.Logger().Info(fmt.Sprintf("Found %d existing delegations from %s", len(delegations), oldProxyAcc.String()))

	recovered := math.NewInt(0)
	for _, del := range delegations {
		valAddr, err := sdk.ValAddressFromBech32(del.ValidatorAddress)
		if err != nil {
			panic("failed to parse validator address from delegation")
		}

		amt, err := keepers.StakingKeeper.Unbond(ctx, oldProxyAcc, valAddr, del.Shares)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("unbonding failed for %s->%s: %+v", del.DelegatorAddress, del.ValidatorAddress, err))
			continue
		}

		recovered = recovered.Add(amt)
	}

	ctx.Logger().Info(fmt.Sprintf("Recovered tokens after unbonding: %s", recovered.String()))

	// TODO: next testnet upgrade send all tokens from oldProxyAcc to liquidstaketypes.LiquidStakeProxyAcc

	if ctx.ChainID() == "ictest-core-1" {
		panic("runLiquidstakeUpgradeMigration done")
	}

	return nil
}
