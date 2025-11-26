package types

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type WhitelistedValsMap map[string]WhitelistedValidator

func (whitelistedValsMap WhitelistedValsMap) IsListed(operatorAddr string) bool {
	if _, ok := whitelistedValsMap[operatorAddr]; ok {
		return true
	}

	return false
}

func GetWhitelistedValsMap(whitelistedValidators []WhitelistedValidator) WhitelistedValsMap {
	whitelistedValsMap := make(WhitelistedValsMap)
	for _, wv := range whitelistedValidators {
		whitelistedValsMap[wv.ValidatorAddress] = wv
	}
	return whitelistedValsMap
}

// Validate validates LiquidValidator.
func (v LiquidValidator) Validate() error {
	_, valErr := sdk.ValAddressFromBech32(v.OperatorAddress)
	if valErr != nil {
		return valErr
	}
	return nil
}

func (v LiquidValidator) GetOperator() sdk.ValAddress {
	if v.OperatorAddress == "" {
		return nil
	}
	addr, err := sdk.ValAddressFromBech32(v.OperatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func (v LiquidValidator) GetDelShares(ctx sdk.Context, sk StakingKeeper) math.LegacyDec {
	del, err := sk.GetDelegation(ctx, LiquidStakeProxyAcc, v.GetOperator())
	if err != nil {
		return math.LegacyZeroDec()
	}
	return del.GetShares()
}

func (v LiquidValidator) GetLiquidTokens(ctx sdk.Context, sk StakingKeeper, onlyBonded bool) math.Int {
	delShares := v.GetDelShares(ctx, sk)
	if !delShares.IsPositive() {
		return math.ZeroInt()
	}
	val, err := sk.Validator(ctx, v.GetOperator())
	if err != nil {
		return math.ZeroInt()
	}
	if onlyBonded && !val.IsBonded() {
		return math.ZeroInt()
	}
	return val.TokensFromSharesTruncated(delShares).TruncateInt()
}

func (v LiquidValidator) GetWeight(whitelistedValsMap WhitelistedValsMap, active bool) math.Int {
	if wv, ok := whitelistedValsMap[v.OperatorAddress]; ok && active {
		return wv.TargetWeight
	}

	return math.ZeroInt()
}

func (v LiquidValidator) GetStatus(activeCondition bool) ValidatorStatus {
	if activeCondition {
		return ValidatorStatusActive
	}

	return ValidatorStatusInactive
}

// ActiveCondition checks the liquid validator could be active by below cases
// - included on whitelist
// - existed valid validator on staking module ( existed, not nil del shares and tokens, valid exchange rate)
// - not tombstoned
func ActiveCondition(validator stakingtypes.Validator, whitelisted, tombstoned bool) bool {
	return whitelisted &&
		!tombstoned &&
		// !Unspecified ==> Bonded, Unbonding, Unbonded
		validator.GetStatus() != stakingtypes.Unspecified &&
		!validator.GetTokens().IsNil() &&
		!validator.GetDelegatorShares().IsNil() &&
		!validator.InvalidExRate()
}

// LiquidValidators is a collection of LiquidValidator
type (
	LiquidValidators       []LiquidValidator
	ActiveLiquidValidators LiquidValidators
)

// MinMaxGap Return the list of LiquidValidator with the maximum gap and minimum gap from the target weight of LiquidValidators, respectively.
func (vs LiquidValidators) MinMaxGap(targetMap, liquidTokenMap map[string]math.Int) (minGapVal, maxGapVal LiquidValidator, amountNeeded math.Int, lastRedelegation bool) {
	maxGap := math.ZeroInt()
	minGap := math.ZeroInt()

	for _, val := range vs {
		gap := liquidTokenMap[val.OperatorAddress].Sub(targetMap[val.OperatorAddress])
		if gap.GT(maxGap) {
			maxGap = gap
			maxGapVal = val
		}
		if gap.LT(minGap) {
			minGap = gap
			minGapVal = val
		}
	}
	amountNeeded = math.MinInt(maxGap, minGap.Abs())
	// lastRedelegation when maxGap validator's liquid token == amountNeeded for redelegation all delShares
	lastRedelegation = amountNeeded.IsPositive() &&
		!targetMap[maxGapVal.OperatorAddress].IsPositive() &&
		liquidTokenMap[maxGapVal.OperatorAddress].Equal(amountNeeded)

	return minGapVal, maxGapVal, amountNeeded, lastRedelegation
}

func (vs LiquidValidators) Len() int {
	return len(vs)
}

func (vs LiquidValidators) TotalLiquidTokens(ctx sdk.Context, sk StakingKeeper, onlyBonded bool) (math.Int, map[string]math.Int) {
	totalLiquidTokens := math.ZeroInt()
	liquidTokenMap := map[string]math.Int{}
	for _, lv := range vs {
		liquidTokens := lv.GetLiquidTokens(ctx, sk, onlyBonded)
		liquidTokenMap[lv.OperatorAddress] = liquidTokens
		totalLiquidTokens = totalLiquidTokens.Add(liquidTokens)
	}
	return totalLiquidTokens, liquidTokenMap
}

func (vs LiquidValidators) Map() map[string]struct{} {
	valsMap := map[string]struct{}{}
	for _, val := range vs {
		valsMap[val.OperatorAddress] = struct{}{}
	}
	return valsMap
}

func (avs ActiveLiquidValidators) Len() int {
	return LiquidValidators(avs).Len()
}

func (avs ActiveLiquidValidators) TotalActiveLiquidTokens(ctx sdk.Context, sk StakingKeeper, onlyBonded bool) (math.Int, map[string]math.Int) {
	return LiquidValidators(avs).TotalLiquidTokens(ctx, sk, onlyBonded)
}

// TotalWeight for active liquid validator
func (avs ActiveLiquidValidators) TotalWeight(whitelistedValsMap WhitelistedValsMap) math.Int {
	totalWeight := math.ZeroInt()
	for _, val := range avs {
		totalWeight = totalWeight.Add(val.GetWeight(whitelistedValsMap, true))
	}
	return totalWeight
}

// NativeTokenToStkXPRT returns StkxprtTotalSupply * nativeTokenAmount / netAmount
func NativeTokenToStkXPRT(nativeTokenAmount, stkXPRTTotalSupplyAmount math.Int, netAmount math.LegacyDec) (stkXPRTAmount math.Int) {
	return math.LegacyNewDecFromInt(stkXPRTTotalSupplyAmount).MulTruncate(math.LegacyNewDecFromInt(nativeTokenAmount)).QuoTruncate(netAmount.TruncateDec()).TruncateInt()
}

// StkXPRTToNativeToken returns stkXPRTAmount * netAmount / StkxprtTotalSupply with truncations
func StkXPRTToNativeToken(stkXPRTAmount, stkXPRTTotalSupplyAmount math.Int, netAmount math.LegacyDec) (nativeTokenAmount math.LegacyDec) {
	return math.LegacyNewDecFromInt(stkXPRTAmount).MulTruncate(netAmount).Quo(math.LegacyNewDecFromInt(stkXPRTTotalSupplyAmount)).TruncateDec()
}

// DeductFeeRate returns Input * (1-FeeRate) with truncations
func DeductFeeRate(input, feeRate math.LegacyDec) (feeDeductedOutput math.LegacyDec) {
	return input.MulTruncate(math.LegacyOneDec().Sub(feeRate)).TruncateDec()
}

func (nas NetAmountState) CalcNetAmount() math.LegacyDec {
	return math.LegacyNewDecFromInt(nas.TotalLiquidTokens.Add(nas.TotalUnbondingBalance))
}

func (nas NetAmountState) CalcMintRate() math.LegacyDec {
	if nas.NetAmount.IsNil() || !nas.NetAmount.IsPositive() {
		return math.LegacyZeroDec()
	}
	return math.LegacyNewDecFromInt(nas.StkxprtTotalSupply).QuoTruncate(nas.NetAmount)
}

type LiquidValidatorStates []LiquidValidatorState

func MustMarshalLiquidValidator(cdc codec.BinaryCodec, val *LiquidValidator) []byte {
	return cdc.MustMarshal(val)
}

// must unmarshal a liquid validator from a store value
func MustUnmarshalLiquidValidator(cdc codec.BinaryCodec, value []byte) LiquidValidator {
	validator, err := UnmarshalLiquidValidator(cdc, value)
	if err != nil {
		panic(err)
	}

	return validator
}

// unmarshal a liquid validator from a store value
func UnmarshalLiquidValidator(cdc codec.BinaryCodec, value []byte) (val LiquidValidator, err error) {
	err = cdc.Unmarshal(value, &val)
	return val, err
}
