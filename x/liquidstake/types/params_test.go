package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/persistenceOne/persistenceCore/v16/app/constants"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/types"
)

func TestParams(t *testing.T) {
	constants.SetUnsealedConfig()

	require.Equal(t, "persistence19zwggtdgaspa9tje6mxdap9xjpc4rayf3nd6dt5g3lwkx4y7z6dqmj3hnc", types.LiquidStakeProxyAcc.String())

	params := types.DefaultParams()

	paramsStr := `{
"liquid_bond_denom": "stk/uxprt",
"whitelisted_validators": [],
"unstake_fee_rate": "0.000000000000000000",
"min_liquid_stake_amount": "1000",
"fee_account_address": "persistence1w2q3mashs2k4wcpqzs5q5xewnhnnr7wslr34safzvwqzvuqh3gjqv4j6ev",
"autocompound_fee_rate": "0.050000000000000000",
"module_paused": true
}`
	require.Equal(t, paramsStr, params.String())

	params.WhitelistedValidators = []types.WhitelistedValidator{
		{
			ValidatorAddress: "persistencevaloper19rz0gtqf88vwk6dwz522ajpqpv5swunqm9z90m",
			TargetWeight:     math.NewInt(10),
		},
	}
	paramsStr = `{
"liquid_bond_denom": "stk/uxprt",
"whitelisted_validators": [
{
"validator_address": "persistencevaloper19rz0gtqf88vwk6dwz522ajpqpv5swunqm9z90m",
"target_weight": "10"
}
],
"unstake_fee_rate": "0.000000000000000000",
"min_liquid_stake_amount": "1000",
"fee_account_address": "persistence1w2q3mashs2k4wcpqzs5q5xewnhnnr7wslr34safzvwqzvuqh3gjqv4j6ev",
"autocompound_fee_rate": "0.050000000000000000",
"module_paused": true
}`
	require.Equal(t, paramsStr, params.String())
}

func TestWhitelistedValsMap(t *testing.T) {
	params := types.DefaultParams()
	require.EqualValues(t, params.WhitelistedValsMap(), types.WhitelistedValsMap{})

	params.WhitelistedValidators = []types.WhitelistedValidator{
		whitelistedValidators[0],
		whitelistedValidators[1],
	}

	wvm := params.WhitelistedValsMap()
	require.Len(t, params.WhitelistedValidators, len(wvm))

	for _, wv := range params.WhitelistedValidators {
		require.EqualValues(t, wvm[wv.ValidatorAddress], wv)
		require.True(t, wvm.IsListed(wv.ValidatorAddress))
	}

	require.False(t, wvm.IsListed("notExistedAddr"))
}

func TestValidateWhitelistedValidators(t *testing.T) {
	for _, tc := range []struct {
		name     string
		malleate func(*types.Params)
		errStr   string
	}{
		{
			"valid default params",
			func(params *types.Params) {},
			"",
		},
		{
			"blank liquid bond denom",
			func(params *types.Params) {
				params.LiquidBondDenom = ""
			},
			"liquid bond denom cannot be blank",
		},
		{
			"invalid liquid bond denom",
			func(params *types.Params) {
				params.LiquidBondDenom = "a"
			},
			"invalid denom: a",
		},
		{
			"duplicated whitelisted validators",
			func(params *types.Params) {
				params.WhitelistedValidators = []types.WhitelistedValidator{
					{
						ValidatorAddress: "persistencevaloper19rz0gtqf88vwk6dwz522ajpqpv5swunqm9z90m",
						TargetWeight:     math.NewInt(10),
					},
					{
						ValidatorAddress: "persistencevaloper19rz0gtqf88vwk6dwz522ajpqpv5swunqm9z90m",
						TargetWeight:     math.NewInt(10),
					},
				}
			},
			"liquidstake validator cannot be duplicated: persistencevaloper19rz0gtqf88vwk6dwz522ajpqpv5swunqm9z90m",
		},
		{
			"invalid whitelisted validator address",
			func(params *types.Params) {
				params.WhitelistedValidators = []types.WhitelistedValidator{
					{
						ValidatorAddress: "invalidaddr",
						TargetWeight:     math.NewInt(10),
					},
				}
			},
			"decoding bech32 failed: invalid separator index -1",
		},
		{
			"nil whitelisted validator target weight",
			func(params *types.Params) {
				params.WhitelistedValidators = []types.WhitelistedValidator{
					{
						ValidatorAddress: "persistencevaloper19rz0gtqf88vwk6dwz522ajpqpv5swunqm9z90m",
						TargetWeight:     math.Int{},
					},
				}
			},
			"liquidstake validator target weight must not be nil",
		},
		{
			"negative whitelisted validator target weight",
			func(params *types.Params) {
				params.WhitelistedValidators = []types.WhitelistedValidator{
					{
						ValidatorAddress: "persistencevaloper19rz0gtqf88vwk6dwz522ajpqpv5swunqm9z90m",
						TargetWeight:     math.NewInt(-1),
					},
				}
			},
			"liquidstake validator target weight must be positive: -1",
		},
		{
			"zero whitelisted validator target weight",
			func(params *types.Params) {
				params.WhitelistedValidators = []types.WhitelistedValidator{
					{
						ValidatorAddress: "persistencevaloper19rz0gtqf88vwk6dwz522ajpqpv5swunqm9z90m",
						TargetWeight:     math.ZeroInt(),
					},
				}
			},
			"liquidstake validator target weight must be positive: 0",
		},
		{
			"nil unstake fee rate",
			func(params *types.Params) {
				params.UnstakeFeeRate = math.LegacyDec{}
			},
			"unstake fee rate must not be nil",
		},
		{
			"negative unstake fee rate",
			func(params *types.Params) {
				params.UnstakeFeeRate = math.LegacyNewDec(-1)
			},
			"unstake fee rate must not be negative: -1.000000000000000000",
		},
		{
			"too large unstake fee rate",
			func(params *types.Params) {
				params.UnstakeFeeRate = math.LegacyMustNewDecFromStr("1.0000001")
			},
			"unstake fee rate too large: 1.000000100000000000",
		},
		{
			"nil min liquid stake amount",
			func(params *types.Params) {
				params.MinLiquidStakeAmount = math.Int{}
			},
			"min liquid stake amount must not be nil",
		},
		{
			"negative min liquid stake amount",
			func(params *types.Params) {
				params.MinLiquidStakeAmount = math.NewInt(-1)
			},
			"min liquid stake amount must not be negative: -1",
		},
		{
			"invalid CwLockedPoolAddress address",
			func(params *types.Params) {
				params.CwLockedPoolAddress = "cw192340924"
			},
			"cannot convert cw contract address to bech32, invalid address: cw192340924, err: decoding bech32 failed: invalid checksum (expected hn8nlx got 340924)",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			params := types.DefaultParams()
			tc.malleate(&params)
			err := params.Validate()
			if tc.errStr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.errStr)
			}
		})
	}
}
