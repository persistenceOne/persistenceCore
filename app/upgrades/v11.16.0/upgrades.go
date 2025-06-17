package v11_16_0

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/keeper"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"

	"github.com/persistenceOne/persistenceCore/v12/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		ctx.Logger().Info("running upgrade instructions...")
		RemoveUnbondedBalance(ctx, args.Keepers.LiquidStakeIBCKeeper,
			"cosmoshub-4", 636,
			sdk.NewInt64Coin("stk/uatom", 1), sdk.NewInt64Coin("uatom", 3244853),
			"persistence1dedp3sl7tu79s50ksfss42tddzh6w3xqzee4nt")
		RemoveUnbondedBalance(ctx, args.Keepers.LiquidStakeIBCKeeper,
			"chihuahua-1", 704,
			sdk.NewInt64Coin("stk/uhuahua", 1), sdk.NewInt64Coin("uhuahua", 4445702485942),
			"persistence174tktspp3r6x3ew9806tw8s9sx5ver87p6qdq9")

		ctx.Logger().Info("completed upgrade instructions...")
		ctx.Logger().Info("running module migrations...")
		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}

// as per https://github.com/persistenceOne/pstake-native/issues/853
func RemoveUnbondedBalance(ctx sdk.Context, liquidStakeIBCKeeper *keeper.Keeper,
	chainID string, epoch int64, stkAmt sdk.Coin, unbondAmt sdk.Coin, refillerAddr string) {
	ctx.Logger().Info("unbonding for", "chain-id", chainID, "epoch", epoch,
		"unbondAmount", unbondAmt.String(), "refillerAddr", refillerAddr)
	ctx.Logger().Info("set user unbonding...")
	liquidStakeIBCKeeper.SetUserUnbonding(ctx, &liquidstakeibctypes.UserUnbonding{
		ChainId:      chainID,
		EpochNumber:  epoch,
		Address:      refillerAddr,
		StkAmount:    stkAmt,
		UnbondAmount: unbondAmt,
	})

	ctx.Logger().Info("set epoch unbonding...")
	liquidStakeIBCKeeper.SetUnbonding(ctx, &liquidstakeibctypes.Unbonding{
		ChainId:       chainID,
		EpochNumber:   epoch,
		MatureTime:    ctx.BlockTime(),
		BurnAmount:    stkAmt,
		UnbondAmount:  unbondAmt,
		IbcSequenceId: "",
		State:         2,
	})
}
