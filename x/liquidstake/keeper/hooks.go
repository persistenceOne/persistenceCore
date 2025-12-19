package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	epochstypes "github.com/cosmos/cosmos-sdk/x/epochs/types"

	liquidstake "github.com/persistenceOne/persistenceCore/v17/x/liquidstake/types"
)

type EpochHooks struct {
	k Keeper
}

var _ epochstypes.EpochHooks = EpochHooks{}

func (k Keeper) EpochHooks() EpochHooks {
	return EpochHooks{k}
}

func (h EpochHooks) BeforeEpochStart(ctx context.Context, epochIdentifier string, epochNumber int64) error {
	return h.k.BeforeEpochStart(ctx, epochIdentifier, epochNumber)
}

func (h EpochHooks) AfterEpochEnd(_ context.Context, _ string, _ int64) error {
	// Nothing to do
	return nil
}

func (k Keeper) BeforeEpochStart(ctx context.Context, epochIdentifier string, _ int64) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params, err := k.GetParams(sdkCtx)
	if err != nil {
		return err
	}
	if !params.ModulePaused {
		// Update the liquid validator set at the start of each epoch
		if epochIdentifier == liquidstake.AutocompoundEpoch {
			k.AutocompoundStakingRewards(sdkCtx, liquidstake.GetWhitelistedValsMap(params.WhitelistedValidators))
		}

		// This has been commented as introducing redelegations for rebalancing affects stkAsset unstake flow
		// https://github.com/cosmos/gaia/security/advisories/GHSA-r47q-464x-wx5x.
		// TODO think of better approach for rebalancing
		//if epochIdentifier == liquidstake.RebalanceEpoch {
		//	// return value of UpdateLiquidValidatorSet is useful only in testing
		//	_ = k.UpdateLiquidValidatorSet(sdkCtx, true)
		//}
	}

	return nil
}
