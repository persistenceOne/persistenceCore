package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/types"
)

var liquidValidators = []types.LiquidValidator{
	{
		OperatorAddress: "persistencevaloper15kdfwczhpmccprekhlzrvkhzw92940l3w37qqj",
	},
	{
		OperatorAddress: "persistencevaloper1x73gyvh74ahs2rt9cqrpjkkk74nczwfpnskv3rczmsf0m6aj5dksqr58m3",
	},
	{
		OperatorAddress: "persistencevaloper10ngyx42lfpylpllm4k3g7fz4gufnt3ptyhm5pn",
	},
	{
		OperatorAddress: "persistencevaloper10fcwju2n8vvffkp8judj3skqpvnphasxjar5yx",
	},
}

func TestDivideByWeight(t *testing.T) {
	testCases := []struct {
		whitelistedVals  []types.WhitelistedValidator
		addStakingAmt    math.Int
		currentDelShares []math.Int
		expectedOutputs  []math.Int
		expectedCrumb    math.Int
	}{
		{
			whitelistedVals: []types.WhitelistedValidator{
				{
					ValidatorAddress: liquidValidators[0].OperatorAddress,
					TargetWeight:     math.NewInt(1),
				},
				{
					ValidatorAddress: liquidValidators[1].OperatorAddress,
					TargetWeight:     math.NewInt(1),
				},
				{
					ValidatorAddress: liquidValidators[2].OperatorAddress,
					TargetWeight:     math.NewInt(1),
				},
			},
			addStakingAmt:    math.NewInt(10 * 1000000),
			currentDelShares: []math.Int{math.NewInt(2000000), math.NewInt(2000000), math.NewInt(1000000)},
			expectedOutputs:  []math.Int{math.NewInt(3333333), math.NewInt(3333333), math.NewInt(3333333)},
			expectedCrumb:    math.NewInt(1),
		},
		{
			whitelistedVals: []types.WhitelistedValidator{
				{
					ValidatorAddress: liquidValidators[0].OperatorAddress,
					TargetWeight:     math.NewInt(2),
				},
				{
					ValidatorAddress: liquidValidators[1].OperatorAddress,
					TargetWeight:     math.NewInt(2),
				},
				{
					ValidatorAddress: liquidValidators[2].OperatorAddress,
					TargetWeight:     math.NewInt(1),
				},
			},
			addStakingAmt:    math.NewInt(10 * 1000000),
			currentDelShares: []math.Int{math.NewInt(1000000), math.NewInt(1000000), math.NewInt(1000000)},
			expectedOutputs:  []math.Int{math.NewInt(4000000), math.NewInt(4000000), math.NewInt(2000000)},
			expectedCrumb:    math.NewInt(0),
		},
		{
			whitelistedVals: []types.WhitelistedValidator{
				{
					ValidatorAddress: liquidValidators[0].OperatorAddress,
					TargetWeight:     math.NewInt(1),
				},
				{
					ValidatorAddress: liquidValidators[1].OperatorAddress,
					TargetWeight:     math.NewInt(1),
				},
				{
					ValidatorAddress: liquidValidators[2].OperatorAddress,
					TargetWeight:     math.NewInt(1),
				},
			},
			addStakingAmt:    math.NewInt(10),
			currentDelShares: []math.Int{math.NewInt(3), math.NewInt(2), math.NewInt(1)},
			expectedOutputs:  []math.Int{math.NewInt(3), math.NewInt(3), math.NewInt(3)},
			expectedCrumb:    math.NewInt(1),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, []types.WhitelistedValidator{}, tc.whitelistedVals)
		require.IsType(t, math.Int{}, tc.addStakingAmt)
		require.IsType(t, math.Int{}, tc.expectedCrumb)
		require.IsType(t, []math.Int{}, tc.expectedOutputs)

		totalTargetAmt := math.ZeroInt()
		valsMap := types.GetWhitelistedValsMap(tc.whitelistedVals)
		var activeVals types.ActiveLiquidValidators
		for _, v := range tc.whitelistedVals {
			activeVals = append(activeVals, types.LiquidValidator{
				OperatorAddress: v.ValidatorAddress,
			})
		}
		outputs, crumb := types.DivideByWeight(activeVals, tc.addStakingAmt, valsMap)
		for _, v := range outputs {
			totalTargetAmt = totalTargetAmt.Add(v)
		}
		require.EqualValues(t, tc.expectedOutputs, outputs)
		require.EqualValues(t, tc.addStakingAmt, totalTargetAmt.Add(crumb))
		require.Equal(t, tc.expectedCrumb.String(), crumb.String())
	}
}

func TestMinMaxGap(t *testing.T) {
	testCases := []struct {
		name                     string
		liquidVals               types.LiquidValidators
		targetMap                map[string]math.Int
		liquidTokenMap           map[string]math.Int
		expectedMinGapVal        types.LiquidValidator
		expectedMaxGapVal        types.LiquidValidator
		expectedAmountNeeded     math.Int
		expectedLastRedelegation bool
	}{
		{
			name:       "zero case",
			liquidVals: liquidValidators,
			targetMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.ZeroInt(),
				liquidValidators[1].OperatorAddress: math.ZeroInt(),
				liquidValidators[2].OperatorAddress: math.ZeroInt(),
				liquidValidators[3].OperatorAddress: math.ZeroInt(),
			},
			liquidTokenMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.ZeroInt(),
				liquidValidators[1].OperatorAddress: math.ZeroInt(),
				liquidValidators[2].OperatorAddress: math.ZeroInt(),
				liquidValidators[3].OperatorAddress: math.ZeroInt(),
			},
			expectedMinGapVal:        types.LiquidValidator{},
			expectedMaxGapVal:        types.LiquidValidator{},
			expectedAmountNeeded:     math.ZeroInt(),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 1-1",
			liquidVals: liquidValidators,
			targetMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(100000000),
				liquidValidators[1].OperatorAddress: math.NewInt(100000000),
				liquidValidators[2].OperatorAddress: math.NewInt(100000000),
				liquidValidators[3].OperatorAddress: math.NewInt(100000000),
			},
			liquidTokenMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(133333334),
				liquidValidators[1].OperatorAddress: math.NewInt(133333333),
				liquidValidators[2].OperatorAddress: math.NewInt(133333333),
				liquidValidators[3].OperatorAddress: math.ZeroInt(),
			},
			expectedMinGapVal:        liquidValidators[3],
			expectedMaxGapVal:        liquidValidators[0],
			expectedAmountNeeded:     math.NewInt(33333334),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 1-2",
			liquidVals: liquidValidators,
			targetMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(100000000),
				liquidValidators[1].OperatorAddress: math.NewInt(100000000),
				liquidValidators[2].OperatorAddress: math.NewInt(100000000),
				liquidValidators[3].OperatorAddress: math.NewInt(100000000),
			},
			liquidTokenMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(133333334 - 33333334),
				liquidValidators[1].OperatorAddress: math.NewInt(133333333),
				liquidValidators[2].OperatorAddress: math.NewInt(133333333),
				liquidValidators[3].OperatorAddress: math.NewInt(0 + 33333334),
			},
			expectedMinGapVal:        liquidValidators[3],
			expectedMaxGapVal:        liquidValidators[1],
			expectedAmountNeeded:     math.NewInt(33333333),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 1-3",
			liquidVals: liquidValidators,
			targetMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(100000000),
				liquidValidators[1].OperatorAddress: math.NewInt(100000000),
				liquidValidators[2].OperatorAddress: math.NewInt(100000000),
				liquidValidators[3].OperatorAddress: math.NewInt(100000000),
			},
			liquidTokenMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(133333334 - 33333334),
				liquidValidators[1].OperatorAddress: math.NewInt(133333333 - 33333333),
				liquidValidators[2].OperatorAddress: math.NewInt(133333333),
				liquidValidators[3].OperatorAddress: math.NewInt(33333334 + 33333333),
			},
			expectedMinGapVal:        liquidValidators[3],
			expectedMaxGapVal:        liquidValidators[2],
			expectedAmountNeeded:     math.NewInt(33333333),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 1-4",
			liquidVals: liquidValidators,
			targetMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(100000000),
				liquidValidators[1].OperatorAddress: math.NewInt(100000000),
				liquidValidators[2].OperatorAddress: math.NewInt(100000000),
				liquidValidators[3].OperatorAddress: math.NewInt(100000000),
			},
			liquidTokenMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(133333334 - 33333334),
				liquidValidators[1].OperatorAddress: math.NewInt(133333333 - 33333333),
				liquidValidators[2].OperatorAddress: math.NewInt(133333333 - 33333333),
				liquidValidators[3].OperatorAddress: math.NewInt(33333334 + 33333333 + 33333333),
			},
			expectedMinGapVal:        types.LiquidValidator{},
			expectedMaxGapVal:        types.LiquidValidator{},
			expectedAmountNeeded:     math.ZeroInt(),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 2-1",
			liquidVals: liquidValidators,
			targetMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(133333334),
				liquidValidators[1].OperatorAddress: math.NewInt(133333333),
				liquidValidators[2].OperatorAddress: math.NewInt(133333333),
				liquidValidators[3].OperatorAddress: math.ZeroInt(),
			},
			liquidTokenMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(100000000),
				liquidValidators[1].OperatorAddress: math.NewInt(100000000),
				liquidValidators[2].OperatorAddress: math.NewInt(100000000),
				liquidValidators[3].OperatorAddress: math.NewInt(100000000),
			},
			expectedMinGapVal:        liquidValidators[0],
			expectedMaxGapVal:        liquidValidators[3],
			expectedAmountNeeded:     math.NewInt(33333334),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 2-2",
			liquidVals: liquidValidators,
			targetMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(133333334),
				liquidValidators[1].OperatorAddress: math.NewInt(133333333),
				liquidValidators[2].OperatorAddress: math.NewInt(133333333),
				liquidValidators[3].OperatorAddress: math.ZeroInt(),
			},
			liquidTokenMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(100000000 + 33333334),
				liquidValidators[1].OperatorAddress: math.NewInt(100000000),
				liquidValidators[2].OperatorAddress: math.NewInt(100000000),
				liquidValidators[3].OperatorAddress: math.NewInt(100000000 - 33333334),
			},
			expectedMinGapVal:        liquidValidators[1],
			expectedMaxGapVal:        liquidValidators[3],
			expectedAmountNeeded:     math.NewInt(33333333),
			expectedLastRedelegation: false,
		},
		{
			name:       "rebalancing case 2-3, last redelegation",
			liquidVals: liquidValidators,
			targetMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(133333334),
				liquidValidators[1].OperatorAddress: math.NewInt(133333333),
				liquidValidators[2].OperatorAddress: math.NewInt(133333333),
				liquidValidators[3].OperatorAddress: math.ZeroInt(),
			},
			liquidTokenMap: map[string]math.Int{
				liquidValidators[0].OperatorAddress: math.NewInt(100000000 + 33333334),
				liquidValidators[1].OperatorAddress: math.NewInt(100000000 + 33333333),
				liquidValidators[2].OperatorAddress: math.NewInt(100000000),
				liquidValidators[3].OperatorAddress: math.NewInt(100000000 - 33333334 - 33333333),
			},
			expectedMinGapVal:        liquidValidators[2],
			expectedMaxGapVal:        liquidValidators[3],
			expectedAmountNeeded:     math.NewInt(33333333),
			expectedLastRedelegation: true,
		},
	}

	for _, tc := range testCases {
		minGapVal, maxGapVal, amountNeeded, last := tc.liquidVals.MinMaxGap(tc.targetMap, tc.liquidTokenMap)
		require.EqualValues(t, minGapVal, tc.expectedMinGapVal)
		require.EqualValues(t, maxGapVal, tc.expectedMaxGapVal)
		require.EqualValues(t, amountNeeded, tc.expectedAmountNeeded)
		require.EqualValues(t, last, tc.expectedLastRedelegation)
	}
}

func TestDivideByCurrentWeight(t *testing.T) {
	testCases := []struct {
		liquidValidators []types.LiquidValidatorState
		addStakingAmt    math.LegacyDec
		expectedOutputs  []math.LegacyDec
		expectedCrumb    math.LegacyDec
	}{
		{
			liquidValidators: []types.LiquidValidatorState{
				{
					OperatorAddress: "a",
					Status:          types.ValidatorStatusActive,
					LiquidTokens:    math.NewIntFromUint64(2 * 1000000),
				},
				{
					OperatorAddress: "b",
					Status:          types.ValidatorStatusActive,
					LiquidTokens:    math.NewIntFromUint64(2 * 1000000),
				},
				{
					OperatorAddress: "c",
					Status:          types.ValidatorStatusActive,
					LiquidTokens:    math.NewIntFromUint64(1 * 1000000),
				},
			},
			addStakingAmt:   math.LegacyNewDec(10 * 1000000),
			expectedOutputs: []math.LegacyDec{math.LegacyNewDec(4 * 1000000), math.LegacyNewDec(4 * 1000000), math.LegacyNewDec(2 * 1000000)},
			expectedCrumb:   math.LegacyNewDec(0),
		},
		{
			liquidValidators: []types.LiquidValidatorState{
				{
					OperatorAddress: "a",
					Status:          types.ValidatorStatusActive,
					LiquidTokens:    math.NewIntFromUint64(1 * 1000000),
					Weight:          math.NewInt(2),
				},
				{
					OperatorAddress: "b",
					Status:          types.ValidatorStatusActive,
					LiquidTokens:    math.NewIntFromUint64(1 * 1000000),
					Weight:          math.NewInt(2),
				},
				{
					OperatorAddress: "c",
					Status:          types.ValidatorStatusActive,
					LiquidTokens:    math.NewIntFromUint64(1 * 1000000),
					Weight:          math.NewInt(1),
				},
			},
			addStakingAmt:   math.LegacyNewDec(10 * 1000000),
			expectedOutputs: []math.LegacyDec{math.LegacyMustNewDecFromStr("3333333.000000000000000000"), math.LegacyMustNewDecFromStr("3333333.000000000000000000"), math.LegacyMustNewDecFromStr("3333333.000000000000000000")},
			expectedCrumb:   math.LegacyMustNewDecFromStr("1.000000000000000000"),
		},
		{
			liquidValidators: []types.LiquidValidatorState{
				{
					OperatorAddress: "a",
					Status:          types.ValidatorStatusActive,
					LiquidTokens:    math.NewIntFromUint64(3),
				},
				{
					OperatorAddress: "b",
					Status:          types.ValidatorStatusActive,
					LiquidTokens:    math.NewIntFromUint64(2),
				},
				{
					OperatorAddress: "c",
					Status:          types.ValidatorStatusActive,
					LiquidTokens:    math.NewIntFromUint64(1),
				},
			},
			addStakingAmt:   math.LegacyNewDec(10),
			expectedOutputs: []math.LegacyDec{math.LegacyMustNewDecFromStr("4.000000000000000000"), math.LegacyMustNewDecFromStr("3.000000000000000000"), math.LegacyMustNewDecFromStr("1.000000000000000000")},
			expectedCrumb:   math.LegacyMustNewDecFromStr("2.000000000000000000"),
		},
		{
			liquidValidators: []types.LiquidValidatorState{
				{
					OperatorAddress: "a",
					Status:          types.ValidatorStatusActive,
					LiquidTokens:    math.NewIntFromUint64(10000000),
				},
				{
					OperatorAddress: "b",
					Status:          types.ValidatorStatusActive,
					LiquidTokens:    math.NewIntFromUint64(2000000),
				},
				{
					OperatorAddress: "c",
					Status:          types.ValidatorStatusActive,
					LiquidTokens:    math.NewIntFromUint64(3000001),
				},
			},
			addStakingAmt:   math.LegacyNewDec(10000000),
			expectedOutputs: []math.LegacyDec{math.LegacyMustNewDecFromStr("6666666.000000000000000000"), math.LegacyMustNewDecFromStr("1333333.000000000000000000"), math.LegacyMustNewDecFromStr("2000000.000000000000000000")},
			expectedCrumb:   math.LegacyMustNewDecFromStr("1.000000000000000000"),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, []types.LiquidValidatorState{}, tc.liquidValidators)
		require.IsType(t, math.LegacyDec{}, tc.addStakingAmt)
		require.IsType(t, math.LegacyDec{}, tc.expectedCrumb)
		require.IsType(t, []math.LegacyDec{}, tc.expectedOutputs)

		totalTargetAmt := math.LegacyZeroDec()
		totalLiquidTokens := math.ZeroInt()
		liquidTokenMap := map[string]math.Int{}
		var lvs types.LiquidValidators
		for _, v := range tc.liquidValidators {
			totalLiquidTokens = totalLiquidTokens.Add(v.LiquidTokens)
			liquidTokenMap[v.OperatorAddress] = v.LiquidTokens
			lvs = append(lvs, types.LiquidValidator{
				OperatorAddress: v.OperatorAddress,
			})
		}
		outputs, crumb := types.DivideByCurrentWeight(lvs, tc.addStakingAmt, totalLiquidTokens, liquidTokenMap)
		for _, v := range outputs {
			totalTargetAmt = totalTargetAmt.Add(v)
		}
		require.EqualValues(t, tc.expectedOutputs, outputs)
		require.EqualValues(t, tc.addStakingAmt, totalTargetAmt.Add(crumb))
		require.Equal(t, tc.expectedCrumb.String(), crumb.String())
	}
}
