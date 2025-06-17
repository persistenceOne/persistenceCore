package v11_14_0

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
		ctx.Logger().Info("running module migrations...")

		RemoveStargazeUnbondedBalance(ctx, args.Keepers.LiquidStakeIBCKeeper)

		return args.ModuleManager.RunMigrations(ctx, args.Configurator, vm)
	}
}

func RemoveStargazeUnbondedBalance(ctx sdk.Context, liquidStakeIBCKeeper *keeper.Keeper) {
	ctx.Logger().Info("starting to move tokens...")

	// as per https://github.com/persistenceOne/pstake-native/issues/853
	chainID := "stargaze-1"
	epoch := int64(582)
	stkAmt := sdk.NewCoin("stk/ustars", sdk.NewInt(60621412694))
	unbondAmt := sdk.NewCoin("ustars", sdk.NewInt(62810179898))
	refillerAddr := "persistence1fp6qhht94pmfdq9h94dvw0tnmnlf2vutnlu7pt"

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

	ctx.Logger().Info("done remove stargaze unbonded balance...")

}
