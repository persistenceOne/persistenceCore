package liquidstake

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/persistenceCore/v17/x/liquidstake/keeper"
	"github.com/persistenceOne/persistenceCore/v17/x/liquidstake/types"
)

func BeginBlock(ctx context.Context, k keeper.Keeper) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params, err := k.GetParams(sdkCtx)
	if err != nil {
		return err
	}
	if !params.ModulePaused {
		// return value of UpdateLiquidValidatorSet is useful only in testing
		_ = k.UpdateLiquidValidatorSet(sdkCtx, false)
	}
	return nil
}
