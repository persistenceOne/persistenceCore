/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package halving

import (
	"fmt"
	"strconv"

	errors "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"

	"github.com/persistenceOne/persistenceCore/v17/x/halving/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) {
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}

	if params.BlockHeight != 0 && uint64(ctx.BlockHeight())%params.BlockHeight == 0 {
		mintParams, err := k.GetMintingParams(ctx)
		if err != nil {
			panic(err)
		}
		newMaxInflation := mintParams.InflationMax.QuoTruncate(sdkmath.LegacyNewDecFromInt(Factor))
		newMinInflation := mintParams.InflationMin.QuoTruncate(sdkmath.LegacyNewDecFromInt(Factor))

		if newMaxInflation.Sub(newMinInflation).LT(sdkmath.LegacyZeroDec()) {
			panic(fmt.Sprintf("max inflation (%s) must be greater than or equal to min inflation (%s)", newMaxInflation.String(), newMinInflation.String()))
		}

		updatedParams := minttypes.NewParams(mintParams.MintDenom, mintParams.InflationRateChange, newMaxInflation, newMinInflation, mintParams.GoalBonded, mintParams.BlocksPerYear)

		if err := k.SetMintingParams(ctx, updatedParams); err != nil {
			panic(errors.Wrap(err, "unable to set minting params at halving EndBlocker"))
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeHalving,
				sdk.NewAttribute(types.AttributeKeyBlockHeight, strconv.FormatInt(ctx.BlockHeight(), 10)),
				sdk.NewAttribute(types.AttributeKeyNewInflationMax, updatedParams.InflationMax.String()),
				sdk.NewAttribute(types.AttributeKeyNewInflationMin, updatedParams.InflationMin.String()),
				sdk.NewAttribute(types.AttributeKeyNewInflationRateChange, updatedParams.InflationRateChange.String()),
			),
		)
	}
}
