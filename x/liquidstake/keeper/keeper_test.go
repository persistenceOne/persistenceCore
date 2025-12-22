package keeper_test

import (
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	evidencetypes "cosmossdk.io/x/evidence/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/persistenceOne/persistenceCore/v17/app/constants"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/persistenceCore/v17/app"
	"github.com/persistenceOne/persistenceCore/v17/x/liquidstake/keeper"
	"github.com/persistenceOne/persistenceCore/v17/x/liquidstake/types"
)

var BlockTime = 6 * time.Second

type KeeperTestSuite struct {
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
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
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

	time, err := sdk.ParseTime("2022-03-01T00:00:00.000000000")
	s.Require().NoError(err)
	s.ctx = s.ctx.WithBlockHeight(100).WithBlockTime(time)
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	params.UnstakeFeeRate = math.LegacyZeroDec()
	params.AutocompoundFeeRate = types.DefaultAutocompoundFeeRate
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	// call mint.BeginBlocker for init k.SetLastBlockTime(ctx, ctx.BlockTime())
	err = mint.BeginBlocker(s.ctx, *s.app.MintKeeper)
	s.Require().NoError(err)
}

func (s *KeeperTestSuite) TearDownTest() {}

func (s *KeeperTestSuite) CreateValidators(powers []int64) ([]sdk.AccAddress, []sdk.ValAddress, []cryptotypes.PubKey) {
	_, err := s.app.BeginBlocker(s.ctx)
	if err != nil {
		return nil, nil, nil
	}
	num := len(powers)
	addrs := simtestutil.AddTestAddrsIncremental(s.app.BankKeeper, s.app.StakingKeeper, s.ctx, num, math.NewInt(10000000000000))
	valAddrs := simtestutil.ConvertAddrsToValAddrs(addrs)
	pks := simtestutil.CreateTestPubKeys(num)
	skParams, err := s.app.LiquidKeeper.GetParams(s.ctx)
	if err != nil {
		s.T().Fatal(err)
	}
	skParams.ValidatorLiquidStakingCap = math.LegacyOneDec()
	_ = s.app.LiquidKeeper.SetParams(s.ctx, skParams)
	for i, power := range powers {
		val, err := stakingtypes.NewValidator(valAddrs[i].String(), pks[i], stakingtypes.Description{})
		s.Require().NoError(err)
		s.app.StakingKeeper.SetValidator(s.ctx, val)
		err = s.app.StakingKeeper.SetValidatorByConsAddr(s.ctx, val)
		s.Require().NoError(err)
		s.app.StakingKeeper.SetNewValidatorByPowerIndex(s.ctx, val)
		_ = s.app.StakingKeeper.Hooks().AfterValidatorCreated(s.ctx, valAddrs[i])
		newShares, err := s.app.StakingKeeper.Delegate(s.ctx, addrs[i], math.NewInt(power), stakingtypes.Unbonded, val, true)
		s.Require().NoError(err)
		s.Require().Equal(newShares.TruncateInt(), math.NewInt(power))
	}

	s.app.EndBlocker(s.ctx)
	return addrs, valAddrs, pks
}

func (s *KeeperTestSuite) liquidStaking(liquidStaker sdk.AccAddress, stakingAmt math.Int) error {
	ctx, writeCache := s.ctx.CacheContext()
	params, err := s.keeper.GetParams(ctx)
	s.Require().NoError(err)

	stkxprtBalanceBefore := s.app.BankKeeper.GetBalance(
		ctx, liquidStaker, params.LiquidBondDenom,
	).Amount

	stkXPRTMintAmt, err := s.keeper.LiquidStake(
		ctx,
		types.LiquidStakeProxyAcc,
		liquidStaker,
		sdk.NewCoin(sdk.DefaultBondDenom, stakingAmt),
	)
	if err != nil {
		return err
	}

	stkxprtBalanceAfter := s.app.BankKeeper.GetBalance(
		ctx, liquidStaker, params.LiquidBondDenom,
	).Amount

	s.Require().NoError(err)
	s.Require().EqualValues(
		stkXPRTMintAmt, stkxprtBalanceAfter.Sub(stkxprtBalanceBefore),
	)
	writeCache()

	return nil
}

func (s *KeeperTestSuite) liquidUnstaking(
	liquidStaker sdk.AccAddress,
	ubdStkXPRTAmt math.Int,
	ubdComplete bool,
) error {
	ctx := s.ctx
	params, err := s.keeper.GetParams(ctx)
	s.Require().NoError(err)

	balanceBefore := s.app.BankKeeper.GetBalance(
		ctx,
		liquidStaker,
		sdk.DefaultBondDenom,
	).Amount

	ubdTime, unbondingAmt, _, unbondedAmt, err := s.liquidUnstakingWithResult(
		liquidStaker,
		sdk.NewCoin(params.LiquidBondDenom, ubdStkXPRTAmt),
	)
	if err != nil {
		return err
	}

	if ubdComplete {
		alv, err := s.keeper.GetActiveLiquidValidators(ctx, params.WhitelistedValsMap())
		if err != nil {
			return err
		}
		ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 200).
			WithBlockTime(ubdTime.Add(1))

		// EndBlock of staking keeper, mature UBD
		s.app.StakingKeeper.BlockValidatorUpdates(ctx)

		balanceCompleteUBD := s.app.BankKeeper.GetBalance(
			ctx,
			liquidStaker,
			sdk.DefaultBondDenom,
		)
		for _, v := range alv {
			_, err := s.app.StakingKeeper.GetUnbondingDelegation(
				ctx,
				liquidStaker,
				v.GetOperator(),
			)
			s.Require().Error(err)
		}

		s.Require().EqualValues(
			balanceCompleteUBD.Amount,
			balanceBefore.Add(unbondingAmt).Add(unbondedAmt),
		)
	}

	return nil
}

func (s *KeeperTestSuite) liquidUnstakingWithResult(
	liquidStaker sdk.AccAddress, unstakingStkXPRT sdk.Coin,
) (time.Time, math.Int, []stakingtypes.UnbondingDelegation, math.Int, error) {
	ctx, writeCache := s.ctx.CacheContext()
	params, err := s.keeper.GetParams(ctx)
	s.Require().NoError(err)
	alv, err := s.keeper.GetActiveLiquidValidators(ctx, params.WhitelistedValsMap())
	if err != nil {
		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), err
	}

	balanceBefore := s.app.BankKeeper.GetBalance(
		ctx, liquidStaker, sdk.DefaultBondDenom,
	).Amount
	stkxprtBalanceBefore := s.app.BankKeeper.GetBalance(
		ctx, liquidStaker, params.LiquidBondDenom,
	).Amount

	ubdTime, unbondingAmt, ubds, unbondedAmt, err := s.keeper.LiquidUnstake(
		ctx, types.LiquidStakeProxyAcc, liquidStaker, unstakingStkXPRT,
	)
	if err != nil {
		return ubdTime, unbondingAmt, ubds, unbondedAmt, err
	}

	balanceAfter := s.app.BankKeeper.GetBalance(
		ctx, liquidStaker, sdk.DefaultBondDenom,
	).Amount
	stkxprtBalanceAfter := s.app.BankKeeper.GetBalance(
		ctx, liquidStaker, params.LiquidBondDenom,
	).Amount
	s.Require().EqualValues(
		unstakingStkXPRT.Amount, stkxprtBalanceBefore.Sub(stkxprtBalanceAfter),
	)

	if unbondedAmt.IsPositive() {
		s.Require().EqualValues(
			unbondedAmt, balanceAfter.Sub(balanceBefore),
		)
	}

	for _, v := range alv {
		_, err := s.app.StakingKeeper.GetUnbondingDelegation(
			ctx, liquidStaker, v.GetOperator(),
		)
		s.Require().NoError(err)
	}

	writeCache()
	return ubdTime, unbondingAmt, ubds, unbondedAmt, err
}

func (s *KeeperTestSuite) RequireNetAmountStateZero() {
	nas, err := s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(nas.MintRate, math.LegacyZeroDec())
	s.Require().EqualValues(nas.StkxprtTotalSupply, math.ZeroInt())
	s.Require().EqualValues(nas.NetAmount, math.LegacyZeroDec())
	s.Require().EqualValues(nas.TotalDelShares, math.LegacyZeroDec())
	s.Require().EqualValues(nas.TotalLiquidTokens, math.ZeroInt())
	s.Require().EqualValues(nas.TotalRemainingRewards, math.LegacyZeroDec())
	s.Require().EqualValues(nas.TotalUnbondingBalance, math.LegacyZeroDec())
	s.Require().EqualValues(nas.ProxyAccBalance, math.ZeroInt())
}

// advance block time and height for complete redelegations and unbondings
func (s *KeeperTestSuite) completeRedelegationUnbonding() {
	headerInfo := s.ctx.HeaderInfo()
	headerInfo.Time = s.ctx.BlockTime().Add(stakingtypes.DefaultUnbondingTime)
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).
		WithBlockTime(s.ctx.BlockTime().Add(stakingtypes.DefaultUnbondingTime)).
		WithHeaderInfo(headerInfo)
	_, err := s.app.EndBlocker(s.ctx)
	s.Require().NoError(err)
	reds, err := s.app.StakingKeeper.GetRedelegations(s.ctx, types.LiquidStakeProxyAcc, 100)
	s.Require().NoError(err)
	s.Require().Len(reds, 0)
	ubds, err := s.app.StakingKeeper.GetUnbondingDelegations(s.ctx, types.LiquidStakeProxyAcc, 100)
	s.Require().NoError(err)
	s.Require().Len(ubds, 0)
}

func (s *KeeperTestSuite) redelegationsErrorCount(redelegations []types.Redelegation) int {
	errCnt := 0
	for _, red := range redelegations {
		if red.Error != nil {
			errCnt++
		}
	}
	return errCnt
}

func (s *KeeperTestSuite) printRedelegationsLiquidTokens() {
	redsIng, err := s.app.StakingKeeper.GetRedelegations(s.ctx, types.LiquidStakeProxyAcc, 50)
	s.Require().NoError(err)

	if len(redsIng) != 0 {
		fmt.Println("[Redelegations]")
		for i, red := range redsIng {
			fmt.Println("\tRedelegation #", i+1)
			fmt.Println("\t\tDelegatorAddress: ", red.DelegatorAddress)
			fmt.Println("\t\tValidatorSrcAddress : ", red.ValidatorSrcAddress)
			fmt.Println("\t\tValidatorDstAddress: ", red.ValidatorDstAddress)
			fmt.Println("\t\tEntries: ")
			for _, e := range red.Entries {
				fmt.Println("\t\t\tCreationHeight: ", e.CreationHeight)
				fmt.Println("\t\t\tCompletionTime: ", e.CompletionTime)
				fmt.Println("\t\t\tInitialBalance: ", e.InitialBalance)
				fmt.Println("\t\t\tSharesDst: ", e.SharesDst)
			}
		}
		fmt.Println("")
	}
	liquidVals := s.keeper.GetAllLiquidValidators(s.ctx)
	if len(liquidVals) != 0 {
		fmt.Println("[LiquidValidators]")
		for _, v := range s.keeper.GetAllLiquidValidators(s.ctx) {
			fmt.Printf("   OperatorAddress %s; LiquidTokens: %s\n",
				v.OperatorAddress, v.GetLiquidTokens(s.ctx, s.app.StakingKeeper, false))
		}
	}
}

func (s *KeeperTestSuite) advanceHeight(height int, _ bool) {
	feeCollector := s.app.AccountKeeper.GetModuleAddress(
		authtypes.FeeCollectorName,
	)

	for i := 0; i < height; i++ {
		s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 1).
			WithBlockTime(s.ctx.BlockTime().Add(BlockTime))

		mint.BeginBlocker(s.ctx, *s.app.MintKeeper)
		feeCollectorBalance := s.app.BankKeeper.GetAllBalances(
			s.ctx, feeCollector,
		)
		rewardsToBeDistributed := feeCollectorBalance.AmountOf(
			sdk.DefaultBondDenom,
		)

		// mimic distribution.BeginBlock (AllocateTokens, get rewards from
		// feeCollector, AllocateTokensToValidator, add remaining to feePool)
		err := s.app.BankKeeper.SendCoinsFromModuleToModule(
			s.ctx, authtypes.FeeCollectorName, distrtypes.ModuleName,
			feeCollectorBalance,
		)

		s.Require().NoError(err)
		totalRewards := math.LegacyZeroDec()
		totalPower := int64(0)
		s.app.StakingKeeper.IterateBondedValidatorsByPower(
			s.ctx,
			func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
				consPower := validator.GetConsensusPower(
					s.app.StakingKeeper.PowerReduction(s.ctx),
				)
				totalPower = totalPower + consPower
				return false
			},
		)

		if totalPower != 0 {
			s.app.StakingKeeper.IterateBondedValidatorsByPower(
				s.ctx,
				func(index int64, validator stakingtypes.ValidatorI) (stop bool) {
					consPower := validator.GetConsensusPower(
						s.app.StakingKeeper.PowerReduction(s.ctx),
					)
					powerFraction := math.LegacyNewDec(consPower).QuoTruncate(
						math.LegacyNewDec(totalPower),
					)
					reward := rewardsToBeDistributed.ToLegacyDec().MulTruncate(
						powerFraction,
					)

					err = s.app.DistributionKeeper.AllocateTokensToValidator(
						s.ctx, validator,
						sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: reward}},
					)
					s.Require().NoError(err)

					totalRewards = totalRewards.Add(reward)
					return false
				},
			)
		}

		remaining := rewardsToBeDistributed.ToLegacyDec().Sub(totalRewards)
		s.Require().False(remaining.GT(math.LegacyNewDec(1)))
		feePool, err := s.app.DistributionKeeper.FeePool.Get(s.ctx)
		s.Require().NoError(err)
		feePool.CommunityPool = feePool.CommunityPool.Add(
			sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: remaining}}...,
		)

		err = s.app.DistributionKeeper.FeePool.Set(s.ctx, feePool)
		s.Require().NoError(err)

		_, err = s.app.StakingKeeper.EndBlocker(s.ctx)
		s.Require().NoError(err)
	}
}

// doubleSign, tombstone, slash, jail
func (s *KeeperTestSuite) doubleSign(valOper sdk.ValAddress, consAddr sdk.ConsAddress) {
	liquidValidator, found := s.keeper.GetLiquidValidator(s.ctx, valOper)
	s.Require().True(found)
	val, err := s.app.StakingKeeper.GetValidator(s.ctx, valOper)
	s.Require().NoError(err)
	tokens := val.Tokens
	liquidTokens := liquidValidator.GetLiquidTokens(s.ctx, s.app.StakingKeeper, false)

	// check sign info
	info, err := s.app.SlashingKeeper.GetValidatorSigningInfo(s.ctx, consAddr)
	s.Require().NoError(err)
	s.Require().Equal(info.Address, consAddr.String())

	// HandleEquivocationEvidence call below functions
	err = s.app.SlashingKeeper.Slash(s.ctx, consAddr, math.LegacyMustNewDecFromStr("0.05"),
		s.app.StakingKeeper.TokensToConsensusPower(s.ctx, tokens), s.ctx.BlockHeight())
	s.Require().NoError(err)
	err = s.app.SlashingKeeper.Jail(s.ctx, consAddr)
	s.Require().NoError(err)
	err = s.app.SlashingKeeper.JailUntil(s.ctx, consAddr, evidencetypes.DoubleSignJailEndTime)
	s.Require().NoError(err)
	err = s.app.SlashingKeeper.Tombstone(s.ctx, consAddr)
	s.Require().NoError(err)

	// should be jailed and tombstoned
	valI, err := s.app.StakingKeeper.Validator(s.ctx, liquidValidator.GetOperator())
	s.Require().NoError(err)
	s.Require().True(valI.IsJailed())
	s.Require().True(s.app.SlashingKeeper.IsTombstoned(s.ctx, consAddr))

	// check tombstoned on sign info
	info, err = s.app.SlashingKeeper.GetValidatorSigningInfo(s.ctx, consAddr)
	s.Require().NoError(err)
	s.Require().True(info.Tombstoned)
	val, _ = s.app.StakingKeeper.GetValidator(s.ctx, valOper)
	s.Require().True(s.keeper.IsTombstoned(s.ctx, val))
	liquidTokensSlashed := liquidValidator.GetLiquidTokens(s.ctx, s.app.StakingKeeper, false)
	tokensSlashed := val.Tokens
	s.Require().True(tokensSlashed.LT(tokens))
	s.Require().True(liquidTokensSlashed.LT(liquidTokens))

	s.app.StakingKeeper.BlockValidatorUpdates(s.ctx)
	val, _ = s.app.StakingKeeper.GetValidator(s.ctx, valOper)

	// set unbonding status, no more rewards before return Bonded
	s.Require().Equal(val.Status, stakingtypes.Unbonding)
}

func (s *KeeperTestSuite) createContinuousVestingAccount(
	from, to sdk.AccAddress, amt sdk.Coins,
	startTime, endTime time.Time,
) vestingtypes.ContinuousVestingAccount {
	baseAccount := s.app.AccountKeeper.NewAccountWithAddress(s.ctx, to)
	_, ok := baseAccount.(*authtypes.BaseAccount)
	s.Require().True(ok)
	baseVestingAccount, err := vestingtypes.NewBaseVestingAccount(
		baseAccount.(*authtypes.BaseAccount), amt, endTime.Unix(),
	)
	s.Require().NoError(err)

	cVestingAcc := vestingtypes.NewContinuousVestingAccountRaw(
		baseVestingAccount, startTime.Unix(),
	)

	s.app.AccountKeeper.SetAccount(s.ctx, cVestingAcc)
	err = s.app.BankKeeper.SendCoins(s.ctx, from, to, amt)
	s.Require().NoError(err)

	return *cVestingAcc
}

func (s *KeeperTestSuite) fundAddr(addr sdk.AccAddress, amt sdk.Coins) {
	err := s.app.BankKeeper.MintCoins(s.ctx, "mint", amt)
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, "mint", addr, amt)
	s.Require().NoError(err)
}
