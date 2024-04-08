package v11_8_0

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/persistenceOne/persistenceCore/v11/app/keepers"
	"github.com/persistenceOne/persistenceCore/v11/app/upgrades"
	liquidstaketypes "github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
	ratesynctypes "github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
)

// This was buggy in mainnet upgrade -> potantial fix/ diff should have been -> https://github.com/persistenceOne/persistenceCore/pull/308
// PR #308 is not merged to preserve upgrade history as is.
func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade handler")

		// make sure this region runs only during CI and mainnet v11 upgrade
		if chainID := ctx.ChainID(); chainID == "core-1" || chainID == "ictest-core-1" {
			if err := runLiquidstakeUpgradeMigration(ctx, args.Keepers); err != nil {
				panic(err)
			}
			if err := runRatesyncUpgradeMigration(ctx, args.Keepers); err != nil {
				panic(err)
			}

			// TODO: more migrations
		}

		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}

// runLiquidstakeUpgradeMigration contains stkXPRT related migrations
func runLiquidstakeUpgradeMigration(
	ctx sdk.Context,
	keepers *keepers.AppKeepers,
) error {
	// run x/liquidstake params upgrade to set the defaults in the state
	lsParams := liquidstaketypes.DefaultParams()
	lsParams.ModulePaused = false
	lsParams.AutocompoundFeeRate = sdk.ZeroDec()
	lsParams.UnstakeFeeRate = sdk.ZeroDec()
	lsParams.FeeAccountAddress = "persistence1ealyadcds02yvsn78he4wntt7tpdqhlhg7y2s6"
	lsParams.CwLockedPoolAddress = "persistence1md28ykl78fgxtuvj8gvntjlqdamqzx33dr80cf0tkqzx5cv0j8aqtkdjp3"

	// TODO: this to be authz
	lsParams.WhitelistAdminAddress = "persistence1ealyadcds02yvsn78he4wntt7tpdqhlhg7y2s6"

	if err := keepers.LiquidStakeKeeper.SetParams(ctx, lsParams); err != nil {
		panic(err)
	}

	// initialize module accounts for x/liquidstake
	ak := keepers.AccountKeeper
	moduleAccsToInitialize := []string{
		liquidstaketypes.ModuleName,

		// not including the other two since they are special
	}

	for _, modAccName := range moduleAccsToInitialize {
		// Get module account and relevant permissions from the accountKeeper.
		addr, perms := ak.GetModuleAddressAndPermissions(modAccName)
		if addr == nil {
			panic(fmt.Sprintf(
				"Did not find %v in `ak.GetModuleAddressAndPermissions`. This is not expected. Skipping.",
				modAccName,
			))
		}

		// Try to get the account in state.
		acc := ak.GetAccount(ctx, addr)
		if acc != nil {
			// Account has been initialized.
			macc, isModuleAccount := acc.(authtypes.ModuleAccountI)
			if isModuleAccount {
				// Module account was correctly initialized. Skipping
				ctx.Logger().Info(fmt.Sprintf(
					"module account %+v was correctly initialized. No-op",
					macc,
				))
				continue
			}
		}

		newModuleAccount := authtypes.NewEmptyModuleAccount(modAccName, perms...)
		maccI := (ak.NewAccount(ctx, newModuleAccount)).(authtypes.ModuleAccountI) // this set the account number
		ak.SetModuleAccount(ctx, maccI)
		ctx.Logger().Info(fmt.Sprintf(
			"Successfully initialized module account in state: %+v",
			newModuleAccount,
		))
	}

	return nil
}

// runRatesyncUpgradeMigration contains ratesync for cvalue on host-chains related migrations
func runRatesyncUpgradeMigration(
	ctx sdk.Context,
	keepers *keepers.AppKeepers,
) error {
	params := ratesynctypes.DefaultParams()
	params.Admin = "persistence1ealyadcds02yvsn78he4wntt7tpdqhlhg7y2s6"
	keepers.RateSyncKeeper.SetParams(ctx, params)

	return nil
}
