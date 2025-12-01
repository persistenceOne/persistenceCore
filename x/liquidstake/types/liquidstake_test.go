package types_test

import (
	"testing"

	"cosmossdk.io/math"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/persistenceCore/v16/app/constants"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/persistenceCore/v16/app"
	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/keeper"
	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/types"
)

var whitelistedValidators = []types.WhitelistedValidator{
	{
		ValidatorAddress: "cosmosvaloper10e4vsut6suau8tk9m6dnrm0slgd6npe3jx5xpv",
		TargetWeight:     math.NewInt(10000),
	},
	{
		ValidatorAddress: "cosmosvaloper18hfzxheyknesfgcrttr5dg50ffnfphtwtar9fz",
		TargetWeight:     math.NewInt(1),
	},
	{
		ValidatorAddress: "cosmosvaloper18hfzxheyknesfgcrttr5dg50ffnfphtwtar9fz",
		TargetWeight:     math.NewInt(-1),
	},
	{
		ValidatorAddress: "cosmosvaloper1ld6vlyy24906u3aqp5lj54f3nsg2592nm9nj5c",
		TargetWeight:     math.NewInt(0),
	},
}

func TestStkXPRTToNativeTokenWithFee(t *testing.T) {
	testCases := []struct {
		stkXPRTAmount            math.Int
		stkXPRTTotalSupplyAmount math.Int
		netAmount                math.LegacyDec
		feeRate                  math.LegacyDec
		expectedOutput           math.LegacyDec
	}{
		// reward added case
		{
			stkXPRTAmount:            math.NewInt(100000000),
			stkXPRTTotalSupplyAmount: math.NewInt(5000000000),
			netAmount:                math.LegacyNewDec(5100000000),
			feeRate:                  math.LegacyMustNewDecFromStr("0.0"),
			expectedOutput:           math.LegacyMustNewDecFromStr("102000000.0"),
		},
		// reward added case with fee
		{
			stkXPRTAmount:            math.NewInt(100000000),
			stkXPRTTotalSupplyAmount: math.NewInt(5000000000),
			netAmount:                math.LegacyNewDec(5100000000),
			feeRate:                  math.LegacyMustNewDecFromStr("0.005"),
			expectedOutput:           math.LegacyMustNewDecFromStr("101490000.0"),
		},
		// slashed case
		{
			stkXPRTAmount:            math.NewInt(100000000),
			stkXPRTTotalSupplyAmount: math.NewInt(5000000000),
			netAmount:                math.LegacyNewDec(4000000000),
			feeRate:                  math.LegacyMustNewDecFromStr("0.0"),
			expectedOutput:           math.LegacyMustNewDecFromStr("80000000.0"),
		},
		// slashed case with fee
		{
			stkXPRTAmount:            math.NewInt(100000000),
			stkXPRTTotalSupplyAmount: math.NewInt(5000000000),
			netAmount:                math.LegacyNewDec(4000000000),
			feeRate:                  math.LegacyMustNewDecFromStr("0.001"),
			expectedOutput:           math.LegacyMustNewDecFromStr("79920000.0"),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, math.Int{}, tc.stkXPRTAmount)
		require.IsType(t, math.Int{}, tc.stkXPRTTotalSupplyAmount)
		require.IsType(t, math.LegacyDec{}, tc.netAmount)
		require.IsType(t, math.LegacyDec{}, tc.feeRate)
		require.IsType(t, math.LegacyDec{}, tc.expectedOutput)

		output := types.StkXPRTToNativeToken(tc.stkXPRTAmount, tc.stkXPRTTotalSupplyAmount, tc.netAmount)
		if tc.feeRate.IsPositive() {
			output = types.DeductFeeRate(output, tc.feeRate)
		}
		require.EqualValues(t, tc.expectedOutput, output)
	}
}

func TestNativeToStkXPRTTo(t *testing.T) {
	testCases := []struct {
		nativeTokenAmount        math.Int
		stkXPRTTotalSupplyAmount math.Int
		netAmount                math.LegacyDec
		expectedOutput           math.Int
	}{
		{
			nativeTokenAmount:        math.NewInt(100000000),
			stkXPRTTotalSupplyAmount: math.NewInt(5000000000),
			netAmount:                math.LegacyNewDec(5000000000),
			expectedOutput:           math.NewInt(100000000),
		},
		{
			nativeTokenAmount:        math.NewInt(100000000),
			stkXPRTTotalSupplyAmount: math.NewInt(5000000000),
			netAmount:                math.LegacyNewDec(4000000000),
			expectedOutput:           math.NewInt(125000000),
		},
		{
			nativeTokenAmount:        math.NewInt(100000000),
			stkXPRTTotalSupplyAmount: math.NewInt(5000000000),
			netAmount:                math.LegacyNewDec(55000000000),
			expectedOutput:           math.NewInt(9090909),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, math.Int{}, tc.nativeTokenAmount)
		require.IsType(t, math.Int{}, tc.stkXPRTTotalSupplyAmount)
		require.IsType(t, math.LegacyDec{}, tc.netAmount)
		require.IsType(t, math.Int{}, tc.expectedOutput)

		output := types.NativeTokenToStkXPRT(tc.nativeTokenAmount, tc.stkXPRTTotalSupplyAmount, tc.netAmount)
		require.EqualValues(t, tc.expectedOutput, output)
	}
}

func TestActiveCondition(t *testing.T) {
	testCases := []struct {
		validator      stakingtypes.Validator
		whitelisted    bool
		tombstoned     bool
		expectedOutput bool
	}{
		// active case 1
		{
			validator: stakingtypes.Validator{
				OperatorAddress: whitelistedValidators[0].ValidatorAddress,
				Jailed:          false,
				Status:          stakingtypes.Bonded,
				Tokens:          math.NewInt(100000000),
				DelegatorShares: math.LegacyNewDec(100000000),
			},
			whitelisted:    true,
			tombstoned:     false,
			expectedOutput: true,
		},
		// active case 2
		{
			validator: stakingtypes.Validator{
				OperatorAddress: whitelistedValidators[0].ValidatorAddress,
				Jailed:          true,
				Status:          stakingtypes.Bonded,
				Tokens:          math.NewInt(100000000),
				DelegatorShares: math.LegacyNewDec(100000000),
			},
			whitelisted:    true,
			tombstoned:     false,
			expectedOutput: true,
		},
		// inactive case 1 (not whitelisted)
		{
			validator: stakingtypes.Validator{
				OperatorAddress: whitelistedValidators[0].ValidatorAddress,
				Jailed:          false,
				Status:          stakingtypes.Bonded,
				Tokens:          math.NewInt(100000000),
				DelegatorShares: math.LegacyNewDec(100000000),
			},
			whitelisted:    false,
			tombstoned:     false,
			expectedOutput: false,
		},
		// inactive case 2 (invalid tokens, delShares)
		{
			validator: stakingtypes.Validator{
				OperatorAddress: whitelistedValidators[0].ValidatorAddress,
				Jailed:          false,
				Status:          stakingtypes.Bonded,
				Tokens:          math.Int{},
				DelegatorShares: math.LegacyDec{},
			},
			whitelisted:    true,
			tombstoned:     false,
			expectedOutput: false,
		},
		// inactive case 3 (zero tokens)
		{
			validator: stakingtypes.Validator{
				OperatorAddress: whitelistedValidators[0].ValidatorAddress,
				Jailed:          false,
				Status:          stakingtypes.Bonded,
				Tokens:          math.NewInt(0),
				DelegatorShares: math.LegacyNewDec(100000000),
			},
			whitelisted:    true,
			tombstoned:     false,
			expectedOutput: false,
		},
		// inactive case 4 (invalid status)
		{
			validator: stakingtypes.Validator{
				OperatorAddress: whitelistedValidators[0].ValidatorAddress,
				Jailed:          false,
				Status:          stakingtypes.Unspecified,
				Tokens:          math.NewInt(100000000),
				DelegatorShares: math.LegacyNewDec(100000000),
			},
			whitelisted:    true,
			tombstoned:     false,
			expectedOutput: false,
		},
		// inactive case 5 (tombstoned)
		{
			validator: stakingtypes.Validator{
				OperatorAddress: whitelistedValidators[0].ValidatorAddress,
				Jailed:          false,
				Status:          stakingtypes.Unbonding,
				Tokens:          math.NewInt(100000000),
				DelegatorShares: math.LegacyNewDec(100000000),
			},
			whitelisted:    true,
			tombstoned:     true,
			expectedOutput: false,
		},
	}

	for _, tc := range testCases {
		require.IsType(t, stakingtypes.Validator{}, tc.validator)
		output := types.ActiveCondition(tc.validator, tc.whitelisted, tc.tombstoned)
		require.EqualValues(t, tc.expectedOutput, output)
	}
}

type TypesTestSuite struct {
	suite.Suite

	app      *app.Application
	ctx      sdk.Context
	keeper   keeper.Keeper
	querier  keeper.Querier
	addrs    []sdk.AccAddress
	delAddrs []sdk.AccAddress
	valAddrs []sdk.ValAddress
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

func (s *TypesTestSuite) SetupTest() {
	constants.SetUnsealedConfig()

	s.app = app.Setup(s.T())
	s.ctx = s.app.BaseApp.NewContext(false)

	stakingParams := stakingtypes.DefaultParams()
	stakingParams.MaxEntries = 7
	stakingParams.MaxValidators = 30
	s.Require().NoError(s.app.StakingKeeper.SetParams(s.ctx, stakingParams))

	s.keeper = *s.app.LiquidStakeKeeper
	s.querier = keeper.Querier{Keeper: s.keeper}
	s.addrs = simtestutil.AddTestAddrs(s.app.BankKeeper, s.app.StakingKeeper, s.ctx, 10, math.NewInt(1_000_000_000))
	s.delAddrs = simtestutil.AddTestAddrs(s.app.BankKeeper, s.app.StakingKeeper, s.ctx, 10, math.NewInt(1_000_000_000))
	s.valAddrs = simtestutil.ConvertAddrsToValAddrs(s.delAddrs)

	blocktime, err := sdk.ParseTime("2022-03-01T00:00:00.000000000")
	s.ctx = s.ctx.WithBlockHeight(100).WithBlockTime(blocktime)
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	params.UnstakeFeeRate = math.LegacyZeroDec()
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	// call mint.BeginBlocker for init k.SetLastBlockTime(ctx, ctx.BlockTime())
	err = mint.BeginBlocker(s.ctx, *s.app.MintKeeper)
	s.Require().NoError(err)
}

func (s *TypesTestSuite) CreateValidators(powers []int64) ([]sdk.AccAddress, []sdk.ValAddress, []cryptotypes.PubKey) {
	s.app.BeginBlocker(s.ctx)
	num := len(powers)
	addrs := simtestutil.AddTestAddrsIncremental(s.app.BankKeeper, s.app.StakingKeeper, s.ctx, num, math.NewInt(1000000000))
	valAddrs := simtestutil.ConvertAddrsToValAddrs(addrs)
	pks := simtestutil.CreateTestPubKeys(num)

	for i, power := range powers {
		val, err := stakingtypes.NewValidator(valAddrs[i].String(), pks[i], stakingtypes.Description{})
		s.Require().NoError(err)
		s.app.StakingKeeper.SetValidator(s.ctx, val)
		err = s.app.StakingKeeper.SetValidatorByConsAddr(s.ctx, val)
		s.Require().NoError(err)
		s.app.StakingKeeper.SetNewValidatorByPowerIndex(s.ctx, val)
		s.app.StakingKeeper.Hooks().AfterValidatorCreated(s.ctx, valAddrs[i])
		newShares, err := s.app.StakingKeeper.Delegate(s.ctx, addrs[i], math.NewInt(power), stakingtypes.Unbonded, val, true)
		s.Require().NoError(err)
		s.Require().Equal(newShares.TruncateInt(), math.NewInt(power))
	}

	s.app.EndBlocker(s.ctx)
	return addrs, valAddrs, pks
}

func (s *TypesTestSuite) TestLiquidStake() {
	_, valOpers, _ := s.CreateValidators([]int64{1000000, 2000000, 3000000})
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)

	params.MinLiquidStakeAmount = math.NewInt(50000)
	params.ModulePaused = false
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	stakingAmt := params.MinLiquidStakeAmount

	// fail, no active validator
	cachedCtx, _ := s.ctx.CacheContext()
	stkXPRTMintAmt, err := s.keeper.LiquidStake(cachedCtx, types.LiquidStakeProxyAcc, s.delAddrs[0], sdk.NewCoin(sdk.DefaultBondDenom, stakingAmt))
	s.Require().ErrorIs(err, types.ErrActiveLiquidValidatorsNotExists)
	s.Require().Equal(stkXPRTMintAmt, math.ZeroInt())

	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(3333)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(3333)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(3333)},
	}
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	res := s.keeper.GetAllLiquidValidatorStates(s.ctx)
	s.Require().Equal(params.WhitelistedValidators[0].ValidatorAddress, res[0].OperatorAddress)
	s.Require().Equal(params.WhitelistedValidators[0].TargetWeight, res[0].Weight)
	s.Require().Equal(types.ValidatorStatusActive, res[0].Status)
	s.Require().Equal(math.LegacyZeroDec(), res[0].DelShares)
	s.Require().Equal(math.ZeroInt(), res[0].LiquidTokens)

	s.Require().Equal(params.WhitelistedValidators[1].ValidatorAddress, res[1].OperatorAddress)
	s.Require().Equal(params.WhitelistedValidators[1].TargetWeight, res[1].Weight)
	s.Require().Equal(types.ValidatorStatusActive, res[1].Status)
	s.Require().Equal(math.LegacyZeroDec(), res[1].DelShares)
	s.Require().Equal(math.ZeroInt(), res[1].LiquidTokens)

	s.Require().Equal(params.WhitelistedValidators[2].ValidatorAddress, res[2].OperatorAddress)
	s.Require().Equal(params.WhitelistedValidators[2].TargetWeight, res[2].Weight)
	s.Require().Equal(types.ValidatorStatusActive, res[2].Status)
	s.Require().Equal(math.LegacyZeroDec(), res[2].DelShares)
	s.Require().Equal(math.ZeroInt(), res[2].LiquidTokens)

	// liquid stake
	stkXPRTMintAmt, err = s.keeper.LiquidStake(s.ctx, types.LiquidStakeProxyAcc, s.delAddrs[0], sdk.NewCoin(sdk.DefaultBondDenom, stakingAmt))
	s.Require().NoError(err)
	s.Require().Equal(stkXPRTMintAmt, stakingAmt)

	_, err = s.app.StakingKeeper.GetDelegation(s.ctx, s.delAddrs[0], valOpers[0])
	s.Require().Error(err)
	_, err = s.app.StakingKeeper.GetDelegation(s.ctx, s.delAddrs[0], valOpers[1])
	s.Require().Error(err)
	_, err = s.app.StakingKeeper.GetDelegation(s.ctx, s.delAddrs[0], valOpers[2])
	s.Require().Error(err)

	proxyAccDel1, err := s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().NoError(err)
	proxyAccDel2, err := s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().NoError(err)
	proxyAccDel3, err := s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[2])
	s.Require().NoError(err)
	s.Require().Equal(proxyAccDel1.Shares, math.LegacyNewDec(16668)) // 16666 + add crumb 2 to 1st active validator
	s.Require().Equal(proxyAccDel2.Shares, math.LegacyNewDec(16666))
	s.Require().Equal(proxyAccDel2.Shares, math.LegacyNewDec(16666))
	s.Require().Equal(stakingAmt.ToLegacyDec(), proxyAccDel1.Shares.Add(proxyAccDel2.Shares).Add(proxyAccDel3.Shares))

	liquidBondDenom, err := s.keeper.LiquidBondDenom(s.ctx)
	s.Require().NoError(err)
	balanceBeforeUBD := s.app.BankKeeper.GetBalance(s.ctx, s.delAddrs[0], sdk.DefaultBondDenom)
	s.Require().Equal(balanceBeforeUBD.Amount, math.NewInt(999950000))
	ubdStkXPRT := sdk.NewCoin(liquidBondDenom, math.NewInt(10000))
	stkXPRTBalance := s.app.BankKeeper.GetBalance(s.ctx, s.delAddrs[0], liquidBondDenom)
	stkXPRTTotalSupply := s.app.BankKeeper.GetSupply(s.ctx, liquidBondDenom)
	s.Require().Equal(stkXPRTBalance, sdk.NewCoin(liquidBondDenom, math.NewInt(50000)))
	s.Require().Equal(stkXPRTBalance, stkXPRTTotalSupply)

	// liquid unstaking
	ubdTime, unbondingAmt, ubds, unbondedAmt, err := s.keeper.LiquidUnstake(s.ctx, types.LiquidStakeProxyAcc, s.delAddrs[0], ubdStkXPRT)
	s.Require().NoError(err)
	s.Require().EqualValues(unbondedAmt, math.ZeroInt())
	s.Require().Len(ubds, 3)

	// crumb excepted on unbonding
	crumb := ubdStkXPRT.Amount.Sub(ubdStkXPRT.Amount.QuoRaw(3).MulRaw(3)) // 1
	s.Require().EqualValues(unbondingAmt, ubdStkXPRT.Amount.Sub(crumb))   // 9999
	s.Require().Equal(ubds[0].DelegatorAddress, s.delAddrs[0].String())
	blocktime, err := sdk.ParseTime("2022-03-22T00:00:00.000000000")
	s.Require().NoError(err)
	s.Require().Equal(ubdTime, blocktime)
	stkXPRTBalanceAfter := s.app.BankKeeper.GetBalance(s.ctx, s.delAddrs[0], liquidBondDenom)
	s.Require().Equal(stkXPRTBalanceAfter, sdk.NewCoin(liquidBondDenom, math.NewInt(40000)))

	balanceBeginUBD := s.app.BankKeeper.GetBalance(s.ctx, s.delAddrs[0], sdk.DefaultBondDenom)
	s.Require().Equal(balanceBeginUBD.Amount, balanceBeforeUBD.Amount)

	proxyAccDel1, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().NoError(err)
	proxyAccDel2, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().NoError(err)
	proxyAccDel3, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[2])
	s.Require().NoError(err)
	s.Require().Equal(stakingAmt.Sub(unbondingAmt).ToLegacyDec(), proxyAccDel1.GetShares().Add(proxyAccDel2.Shares).Add(proxyAccDel3.Shares))

	// complete unbonding
	s.ctx = s.ctx.WithBlockHeight(200).WithBlockTime(ubdTime.Add(1))
	updates, err := s.app.StakingKeeper.BlockValidatorUpdates(s.ctx) // EndBlock of staking keeper, mature UBD
	s.Require().NoError(err)
	s.Require().Empty(updates)
	balanceCompleteUBD := s.app.BankKeeper.GetBalance(s.ctx, s.delAddrs[0], sdk.DefaultBondDenom)
	s.Require().Equal(balanceCompleteUBD.Amount, balanceBeforeUBD.Amount.Add(unbondingAmt))

	proxyAccDel1, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().NoError(err)
	proxyAccDel2, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().NoError(err)
	proxyAccDel3, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[2])
	s.Require().NoError(err)
	// crumb added to first valid active liquid validator
	s.Require().Equal(math.LegacyNewDec(13335), proxyAccDel1.Shares)
	s.Require().Equal(math.LegacyNewDec(13333), proxyAccDel2.Shares)
	s.Require().Equal(math.LegacyNewDec(13333), proxyAccDel3.Shares)

	res = s.keeper.GetAllLiquidValidatorStates(s.ctx)
	s.Require().Equal(params.WhitelistedValidators[0].ValidatorAddress, res[0].OperatorAddress)
	s.Require().Equal(params.WhitelistedValidators[0].TargetWeight, res[0].Weight)
	s.Require().Equal(types.ValidatorStatusActive, res[0].Status)
	s.Require().Equal(math.LegacyNewDec(13335), res[0].DelShares)

	s.Require().Equal(params.WhitelistedValidators[1].ValidatorAddress, res[1].OperatorAddress)
	s.Require().Equal(params.WhitelistedValidators[1].TargetWeight, res[1].Weight)
	s.Require().Equal(types.ValidatorStatusActive, res[1].Status)
	s.Require().Equal(math.LegacyNewDec(13333), res[1].DelShares)

	s.Require().Equal(params.WhitelistedValidators[2].ValidatorAddress, res[2].OperatorAddress)
	s.Require().Equal(params.WhitelistedValidators[2].TargetWeight, res[2].Weight)
	s.Require().Equal(types.ValidatorStatusActive, res[2].Status)
	s.Require().Equal(math.LegacyNewDec(13333), res[2].DelShares)

	vs := s.keeper.GetAllLiquidValidators(s.ctx)
	s.Require().Len(vs.Map(), 3)

	whitelistedValsMap := types.GetWhitelistedValsMap(params.WhitelistedValidators)
	avs, err := s.keeper.GetActiveLiquidValidators(s.ctx, whitelistedValsMap)
	s.Require().NoError(err)
	alt, _ := avs.TotalActiveLiquidTokens(s.ctx, s.app.StakingKeeper, true)
	s.Require().EqualValues(alt, math.NewInt(40001))
}
