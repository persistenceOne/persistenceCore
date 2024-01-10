package v10_3_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"

	"github.com/persistenceOne/persistenceCore/v10/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade handler")

		// stuck unbonding epoch numbers
		RemovableUnbondings := map[string]map[int64]any{
			"cosmoshub-4":       {312: nil},
			"osmosis-1":         {429: nil, 432: nil},
			"theta-testnet-001": {104: nil, 128: nil, 140: nil, 148: nil},
			"osmo-test-5":       {193: nil, 194: nil, 201: nil},
		}

		// get the stuck unbondings from the store
		unbondings := args.Keepers.LiquidStakeIBCKeeper.FilterUnbondings(
			ctx,
			func(u liquidstakeibctypes.Unbonding) bool {
				_, chain := RemovableUnbondings[u.ChainId]
				if chain {
					_, epoch := RemovableUnbondings[u.ChainId][u.EpochNumber]
					if epoch {
						return true
					}
				}
				return false
			},
		)

		// mark the stuck unbondings as failed, so they can be processed
		for _, unbonding := range unbondings {
			unbonding.State = liquidstakeibctypes.Unbonding_UNBONDING_FAILED
			args.Keepers.LiquidStakeIBCKeeper.SetUnbonding(ctx, unbonding)
		}

		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}
