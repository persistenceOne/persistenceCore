package keeper

import (
	"fmt"
	"strconv"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/types"
)

func (k Keeper) GetProxyAccBalance(ctx sdk.Context, proxyAcc sdk.AccAddress) (sdk.Coin, error) {
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return sdk.Coin{}, err
	}
	return sdk.NewCoin(bondDenom, k.bankKeeper.SpendableCoins(ctx, proxyAcc).AmountOf(bondDenom)), nil
}

// TryRedelegation attempts redelegation, which is applied only when successful through cached context because there is a constraint that fails if already receiving redelegation.
func (k Keeper) TryRedelegation(ctx sdk.Context, re types.Redelegation) (completionTime time.Time, err error) {
	dstVal := re.DstValidator.GetOperator()
	srcVal := re.SrcValidator.GetOperator()

	// check the source validator already has receiving transitive redelegation
	hasReceiving, err := k.stakingKeeper.HasReceivingRedelegation(ctx, re.Delegator, srcVal)
	if err != nil {
		return time.Time{}, err
	}
	if hasReceiving {
		return time.Time{}, stakingtypes.ErrTransitiveRedelegation
	}

	// calculate delShares from tokens with validation
	_, err = k.stakingKeeper.ValidateUnbondAmount(
		ctx, re.Delegator, srcVal, re.Amount,
	)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to validate unbond amount: %w", err)
	}

	// when last, full redelegation of shares from delegation
	amt := re.Amount
	if re.Last {
		amt = re.SrcValidator.GetLiquidTokens(ctx, k.stakingKeeper, false)
	}
	cachedCtx, writeCache := ctx.CacheContext()
	completionTime, err = k.RedelegateWithCap(cachedCtx, re.Delegator, srcVal, dstVal, amt)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to begin redelegation: %w", err)
	}
	writeCache()
	return completionTime, nil
}

// Rebalance argument liquidVals containing ValidatorStatusActive which is containing just added on whitelist(liquidToken 0) and ValidatorStatusInactive to delist
func (k Keeper) Rebalance(
	ctx sdk.Context,
	proxyAcc sdk.AccAddress,
	liquidVals types.LiquidValidators,
	whitelistedValsMap types.WhitelistedValsMap,
	rebalancingTrigger math.LegacyDec,
) (redelegations []types.Redelegation) {
	totalLiquidTokens, liquidTokenMap := liquidVals.TotalLiquidTokens(ctx, k.stakingKeeper, false)
	if !totalLiquidTokens.IsPositive() {
		return redelegations
	}

	weightMap, totalWeight := k.GetWeightMap(ctx, liquidVals, whitelistedValsMap)

	// no active liquid validators
	if !totalWeight.IsPositive() {
		return redelegations
	}

	// calculate rebalancing target map
	targetMap := map[string]math.Int{}
	totalTargetMap := math.ZeroInt()
	for _, val := range liquidVals {
		targetMap[val.OperatorAddress] = totalLiquidTokens.Mul(weightMap[val.OperatorAddress]).Quo(totalWeight)
		totalTargetMap = totalTargetMap.Add(targetMap[val.OperatorAddress])
	}
	crumb := totalLiquidTokens.Sub(totalTargetMap)
	if !totalTargetMap.IsPositive() {
		return redelegations
	}
	// crumb to first non zero liquid validator
	for _, val := range liquidVals {
		if targetMap[val.OperatorAddress].IsPositive() {
			targetMap[val.OperatorAddress] = targetMap[val.OperatorAddress].Add(crumb)
			break
		}
	}

	failCount := 0
	rebalancingThresholdAmt := rebalancingTrigger.Mul(math.LegacyNewDecFromInt(totalLiquidTokens)).TruncateInt()
	redelegations = make([]types.Redelegation, 0, liquidVals.Len())

	for i := 0; i < liquidVals.Len(); i++ {
		// get min, max of liquid token gap
		minVal, maxVal, amountNeeded, last := liquidVals.MinMaxGap(targetMap, liquidTokenMap)
		if amountNeeded.IsZero() || (i == 0 && !amountNeeded.GT(rebalancingThresholdAmt)) {
			break
		}

		// sync liquidTokenMap applied rebalancing
		liquidTokenMap[maxVal.OperatorAddress] = liquidTokenMap[maxVal.OperatorAddress].Sub(amountNeeded)
		liquidTokenMap[minVal.OperatorAddress] = liquidTokenMap[minVal.OperatorAddress].Add(amountNeeded)

		// try redelegation from max validator to min validator
		redelegation := types.Redelegation{
			Delegator:    proxyAcc,
			SrcValidator: maxVal,
			DstValidator: minVal,
			Amount:       amountNeeded,
			Last:         last,
		}

		_, err := k.TryRedelegation(ctx, redelegation)
		if err != nil {
			redelegation.Error = err
			failCount++

			k.Logger(ctx).Info(
				"redelegation failed",
				types.DelegatorKeyVal, proxyAcc.String(),
				types.SrcValidatorKeyVal, maxVal.OperatorAddress,
				types.DstValidatorKeyVal, minVal.OperatorAddress,
				types.AmountKeyVal, amountNeeded.String(),
				types.ErrorKeyVal, err.Error(),
			)
		}

		redelegations = append(redelegations, redelegation)
	}

	if len(redelegations) != 0 {
		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeBeginRebalancing,
				sdk.NewAttribute(types.AttributeKeyDelegator, types.LiquidStakeProxyAcc.String()),
				sdk.NewAttribute(types.AttributeKeyRedelegationCount, strconv.Itoa(len(redelegations))),
				sdk.NewAttribute(types.AttributeKeyRedelegationFailCount, strconv.Itoa(failCount)),
			),
		})
		k.Logger(ctx).Info(types.EventTypeBeginRebalancing,
			types.AttributeKeyDelegator, types.LiquidStakeProxyAcc.String(),
			types.AttributeKeyRedelegationCount, strconv.Itoa(len(redelegations)),
			types.AttributeKeyRedelegationFailCount, strconv.Itoa(failCount))
	}

	return redelegations
}

func (k Keeper) UpdateLiquidValidatorSet(ctx sdk.Context, redelegate bool) (redelegations []types.Redelegation) {
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	liquidValidators := k.GetAllLiquidValidators(ctx)
	liquidValsMap := liquidValidators.Map()
	whitelistedValsMap := types.GetWhitelistedValsMap(params.WhitelistedValidators)

	// Set Liquid validators for added whitelist validators
	for _, wv := range params.WhitelistedValidators {
		if _, ok := liquidValsMap[wv.ValidatorAddress]; !ok {
			lv := types.LiquidValidator{
				OperatorAddress: wv.ValidatorAddress,
			}
			if k.IsActiveLiquidValidator(ctx, lv, whitelistedValsMap) {
				k.SetLiquidValidator(ctx, lv)
				liquidValidators = append(liquidValidators, lv)
				ctx.EventManager().EmitEvents(sdk.Events{
					sdk.NewEvent(
						types.EventTypeAddLiquidValidator,
						sdk.NewAttribute(types.AttributeKeyLiquidValidator, lv.OperatorAddress),
					),
				})
				k.Logger(ctx).Info(types.EventTypeAddLiquidValidator, types.AttributeKeyLiquidValidator, lv.OperatorAddress)
			}
		}
	}

	// rebalancing based updated liquid validators status with threshold, try by cachedCtx
	// tombstone status also handled on Rebalance
	if redelegate {
		redelegations = k.Rebalance(
			ctx,
			types.LiquidStakeProxyAcc,
			liquidValidators,
			whitelistedValsMap,
			types.RebalancingTrigger,
		)

		// if there are inactive liquid validators, do not unbond,
		// instead let validator selection and rebalancing take care of it.

		return redelegations
	}
	return nil
}

// AutocompoundStakingRewards withdraws staking rewards and re-stakes when over threshold.
func (k Keeper) AutocompoundStakingRewards(ctx sdk.Context, whitelistedValsMap types.WhitelistedValsMap) {
	// withdraw rewards of LiquidStakeProxyAcc
	k.WithdrawLiquidRewards(ctx, types.LiquidStakeProxyAcc)

	// skip when no active liquid validator
	activeVals, err := k.GetActiveLiquidValidators(ctx, whitelistedValsMap)
	if err != nil {
		return
	}
	if len(activeVals) == 0 {
		return
	}

	// get all the APY components
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return
	}
	totalSupply := k.bankKeeper.GetSupply(ctx, bondDenom).Amount
	bondedTokens := k.bankKeeper.GetBalance(ctx, k.stakingKeeper.GetBondedPool(ctx).GetAddress(), bondDenom).Amount
	minter, err := k.mintKeeper.Minter.Get(ctx)
	if err != nil {
		return
	}
	inflation := minter.Inflation
	// calculate the hourly APY
	bondRatio := math.LegacyDec(bondedTokens).Quo(math.LegacyDec(totalSupply))
	hourlyApy := inflation.Quo(bondRatio).
		Quo(types.DefaultLimitAutocompoundPeriodDays).
		Quo(types.DefaultLimitAutocompoundPeriodHours)

	// calculate autocompoundable amount by limiting the current net amount with the calculated APY
	nas, err := k.GetNetAmountState(ctx)
	if err != nil {
		return
	}
	autoCompoundableAmount := nas.NetAmount.Mul(hourlyApy).TruncateInt()

	// use the calculated autocompoundable amount as the limit for the transfer
	proxyAccBalance, err := k.GetProxyAccBalance(ctx, types.LiquidStakeProxyAcc)
	if err != nil {
		return
	}
	if proxyAccBalance.Amount.LT(autoCompoundableAmount) {
		autoCompoundableAmount = proxyAccBalance.Amount
	}

	// calculate autocompounding fee
	params, err := k.GetParams(ctx)
	if err != nil {
		return
	}
	autocompoundFee := sdk.NewCoin(bondDenom, math.ZeroInt())
	if !params.AutocompoundFeeRate.IsZero() && autoCompoundableAmount.IsPositive() {
		autocompoundFee = sdk.NewCoin(
			bondDenom,
			params.AutocompoundFeeRate.MulInt(autoCompoundableAmount).TruncateInt(),
		)
	}

	// re-staking of the accumulated rewards
	cachedCtx, writeCache := ctx.CacheContext()
	delegableAmount := autoCompoundableAmount.Sub(autocompoundFee.Amount)
	err = k.LiquidDelegate(cachedCtx, types.LiquidStakeProxyAcc, activeVals, delegableAmount, whitelistedValsMap)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to re-stake the accumulated rewards",
			types.ErrorKeyVal,
			err,
		)
		return
		// skip errors as they might occur due to reaching global liquid cap
	}
	writeCache()

	// move autocompounding fee from the balance to fee account
	feeAccountAddr := sdk.MustAccAddressFromBech32(params.FeeAccountAddress)
	err = k.bankKeeper.SendCoins(ctx, types.LiquidStakeProxyAcc, feeAccountAddr, sdk.NewCoins(autocompoundFee))
	if err != nil {
		k.Logger(ctx).Error(
			"failed to send autocompound fee to fee account",
			types.ErrorKeyVal,
			err,
		)
		return
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeAutocompound,
			sdk.NewAttribute(types.AttributeKeyDelegator, types.LiquidStakeProxyAcc.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, delegableAmount.String()),
			sdk.NewAttribute(types.AttributeKeyAutocompoundFee, autocompoundFee.String()),
		),
	})
	k.Logger(ctx).Info(types.EventTypeAutocompound,
		types.AttributeKeyDelegator, types.LiquidStakeProxyAcc.String(),
		sdk.AttributeKeyAmount, delegableAmount.String(),
		types.AttributeKeyAutocompoundFee, autocompoundFee.String())
}
