package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/types"
)

func (s *KeeperTestSuite) TestRebalancingCase1() {
	_, valOpers, pks := s.CreateValidators([]int64{1000000, 1000000, 1000000, 1000000, 1000000})
	blocktime, err := sdk.ParseTime("2022-03-01T00:00:00.000000000")
	s.Require().NoError(err)
	s.ctx = s.ctx.WithBlockHeight(100).WithBlockTime(blocktime)
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	params.UnstakeFeeRate = math.LegacyZeroDec()
	params.MinLiquidStakeAmount = math.NewInt(10000)
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	stakingAmt := math.NewInt(49998)
	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(3000)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(3000)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(3000)},
	}
	params.ModulePaused = false
	s.keeper.SetParams(s.ctx, params)
	reds := s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 0)

	stkXPRTMintAmt, err := s.keeper.LiquidStake(s.ctx, types.LiquidStakeProxyAcc, s.delAddrs[0], sdk.NewCoin(sdk.DefaultBondDenom, stakingAmt))
	s.Require().NoError(err)
	s.Require().Equal(stkXPRTMintAmt, stakingAmt)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 0)

	proxyAccDel1, err := s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().NoError(err)
	proxyAccDel2, err := s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().NoError(err)
	proxyAccDel3, err := s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[2])
	s.Require().NoError(err)

	s.Require().EqualValues(proxyAccDel1.Shares.TruncateInt(), math.NewInt(16668))
	s.Require().EqualValues(proxyAccDel2.Shares.TruncateInt(), math.NewInt(16665))
	s.Require().EqualValues(proxyAccDel3.Shares.TruncateInt(), math.NewInt(16665))
	totalLiquidTokens, _ := s.keeper.GetAllLiquidValidators(s.ctx).TotalLiquidTokens(s.ctx, s.app.StakingKeeper, false)
	s.Require().EqualValues(stakingAmt, totalLiquidTokens)
	s.printRedelegationsLiquidTokens()

	// update whitelist validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(2500)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(2500)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(2500)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(2500)},
	}
	s.keeper.SetParams(s.ctx, params)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 3)

	proxyAccDel1, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().NoError(err)
	proxyAccDel2, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().NoError(err)
	proxyAccDel3, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[2])
	s.Require().NoError(err)
	proxyAccDel4, err := s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[3])
	s.Require().NoError(err)

	s.Require().EqualValues(proxyAccDel1.Shares.TruncateInt(), math.NewInt(12501))
	s.Require().EqualValues(proxyAccDel2.Shares.TruncateInt(), math.NewInt(12499))
	s.Require().EqualValues(proxyAccDel3.Shares.TruncateInt(), math.NewInt(12499))
	s.Require().EqualValues(proxyAccDel4.Shares.TruncateInt(), math.NewInt(12499))
	totalLiquidTokens, _ = s.keeper.GetAllLiquidValidators(s.ctx).TotalLiquidTokens(s.ctx, s.app.StakingKeeper, false)
	s.Require().EqualValues(stakingAmt, totalLiquidTokens)
	s.printRedelegationsLiquidTokens()

	// reds := s.app.StakingKeeper.GetRedelegations(s.ctx, types.LiquidStakeProxyAcc, 20)
	s.Require().Len(reds, 3)

	nas, err := s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)

	// advance block time and height for complete redelegations
	s.completeRedelegationUnbonding()

	nas, err = s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)
	// update whitelist validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(2000)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(2000)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(2000)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(2000)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(2000)},
	}
	s.keeper.SetParams(s.ctx, params)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 4)

	proxyAccDel1, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().NoError(err)
	proxyAccDel2, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().NoError(err)
	proxyAccDel3, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[2])
	s.Require().NoError(err)
	proxyAccDel4, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[3])
	s.Require().NoError(err)
	proxyAccDel5, err := s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[4])
	s.Require().NoError(err)

	s.printRedelegationsLiquidTokens()
	s.Require().EqualValues(proxyAccDel1.Shares.TruncateInt(), math.NewInt(10002))
	s.Require().EqualValues(proxyAccDel2.Shares.TruncateInt(), math.NewInt(9999))
	s.Require().EqualValues(proxyAccDel3.Shares.TruncateInt(), math.NewInt(9999))
	s.Require().EqualValues(proxyAccDel4.Shares.TruncateInt(), math.NewInt(9999))
	s.Require().EqualValues(proxyAccDel5.Shares.TruncateInt(), math.NewInt(9999))
	totalLiquidTokens, _ = s.keeper.GetAllLiquidValidators(s.ctx).TotalLiquidTokens(s.ctx, s.app.StakingKeeper, false)
	s.Require().EqualValues(stakingAmt, totalLiquidTokens)

	// advance block time and height for complete redelegations
	s.completeRedelegationUnbonding()

	// remove whitelist validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(2500)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(2500)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(2500)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(2500)},
	}

	s.keeper.SetParams(s.ctx, params)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 4)

	proxyAccDel1, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().NoError(err)
	proxyAccDel2, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().NoError(err)
	proxyAccDel3, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[2])
	s.Require().NoError(err)
	proxyAccDel4, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[3])
	s.Require().NoError(err)
	proxyAccDel5, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[4])
	s.Require().Error(err)

	s.printRedelegationsLiquidTokens()
	s.Require().EqualValues(proxyAccDel1.Shares.TruncateInt(), math.NewInt(12501))
	s.Require().EqualValues(proxyAccDel2.Shares.TruncateInt(), math.NewInt(12499))
	s.Require().EqualValues(proxyAccDel3.Shares.TruncateInt(), math.NewInt(12499))
	s.Require().EqualValues(proxyAccDel4.Shares.TruncateInt(), math.NewInt(12499))
	totalLiquidTokens, _ = s.keeper.GetAllLiquidValidators(s.ctx).TotalLiquidTokens(s.ctx, s.app.StakingKeeper, false)
	s.Require().EqualValues(stakingAmt, totalLiquidTokens)

	// advance block time and height for complete redelegations
	s.completeRedelegationUnbonding()

	// remove whitelist validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(5000)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(5000)},
	}

	s.keeper.SetParams(s.ctx, params)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 3)

	proxyAccDel1, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().NoError(err)
	proxyAccDel2, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().NoError(err)
	proxyAccDel3, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[2])
	s.Require().Error(err)
	proxyAccDel4, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[3])
	s.Require().Error(err)
	proxyAccDel5, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[4])
	s.Require().Error(err)

	s.printRedelegationsLiquidTokens()
	s.Require().EqualValues(proxyAccDel1.Shares.TruncateInt(), math.NewInt(24999))
	s.Require().EqualValues(proxyAccDel2.Shares.TruncateInt(), math.NewInt(24999))
	totalLiquidTokens, _ = s.keeper.GetAllLiquidValidators(s.ctx).TotalLiquidTokens(s.ctx, s.app.StakingKeeper, false)
	s.Require().EqualValues(stakingAmt, totalLiquidTokens)

	// advance block time and height for complete redelegations
	s.completeRedelegationUnbonding()

	// double sign, tombstone, slash, jail
	s.doubleSign(valOpers[1], sdk.ConsAddress(pks[1].Address()))

	// check inactive with zero weight after tombstoned
	proxyAccDel2Valaddr, err := sdk.ValAddressFromBech32(proxyAccDel2.GetValidatorAddr())
	s.Require().NoError(err)
	lvState, found := s.keeper.GetLiquidValidatorState(s.ctx, proxyAccDel2Valaddr)
	s.Require().True(found)
	s.Require().Equal(lvState.Status, types.ValidatorStatusInactive)
	s.Require().Equal(lvState.Weight, math.ZeroInt())
	s.Require().NotEqualValues(lvState.DelShares, math.LegacyZeroDec())
	s.Require().NotEqualValues(lvState.LiquidTokens, math.ZeroInt())

	// rebalancing, remove tombstoned liquid validator
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 1)

	// all redelegated, no delShares ( exception, dust )
	proxyAccDel2, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[1])
	s.Require().NoError(err)
	s.Require().True(proxyAccDel2.Shares.LT(math.LegacyOneDec()))

	// liquid validator removed, invalid after tombstoned
	lvState, found = s.keeper.GetLiquidValidatorState(s.ctx, valOpers[1])
	s.Require().True(found)
	s.Require().Equal(lvState.OperatorAddress, valOpers[1].String())
	s.Require().Equal(lvState.Status, types.ValidatorStatusInactive)
	s.Require().True(proxyAccDel2.Shares.LT(math.LegacyOneDec()))
	s.Require().True(lvState.LiquidTokens.Equal(math.ZeroInt()))

	// jail last liquid validator, undelegate all liquid tokens to proxy acc
	nasBefore, err := s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)
	s.doubleSign(valOpers[0], sdk.ConsAddress(pks[0].Address()))
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 0)

	// no delegation of proxy acc
	proxyAccDel1, err = s.app.StakingKeeper.GetDelegation(s.ctx, types.LiquidStakeProxyAcc, valOpers[0])
	s.Require().NoError(err)
	val1, err := s.app.StakingKeeper.GetValidator(s.ctx, valOpers[0])
	s.Require().NoError(err)
	s.Require().Equal(val1.Status, stakingtypes.Unbonding)

	// complete unbonding
	s.completeRedelegationUnbonding()

	// check validator Unbonded
	val1, err = s.app.StakingKeeper.GetValidator(s.ctx, valOpers[0])
	s.Require().NoError(err)
	s.Require().Equal(val1.Status, stakingtypes.Unbonded)

	// no rewards, same delShares, liquid tokens as we do not unbond now
	nas, err = s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(nas.TotalRemainingRewards, math.LegacyZeroDec())
	s.Require().EqualValues(nas.TotalDelShares, nasBefore.TotalDelShares)
	s.Require().LessOrEqual(nas.TotalLiquidTokens.Int64(), nasBefore.TotalLiquidTokens.Int64()) // slashing

	// mintRate over 1 due to slashing
	s.Require().True(nas.MintRate.GT(math.LegacyOneDec()))
	stkXPRTBalanceBefore := s.app.BankKeeper.GetBalance(s.ctx, s.delAddrs[0], params.LiquidBondDenom).Amount
	s.Require().EqualValues(nas.StkxprtTotalSupply, stkXPRTBalanceBefore)
}

func (s *KeeperTestSuite) TestRebalancingConsecutiveCase() {
	_, valOpers, _ := s.CreateValidators([]int64{
		1000000000000, 1000000000000, 1000000000000, 1000000000000, 1000000000000,
		1000000000000, 1000000000000, 1000000000000, 1000000000000, 1000000000000,
		1000000000000, 1000000000000, 1000000000000, 1000000000000, 1000000000000,
	})
	headerInfo := s.ctx.HeaderInfo()
	blocktime, err := sdk.ParseTime("2022-03-01T00:00:00.000000000")
	headerInfo.Time = blocktime
	s.ctx = s.ctx.WithBlockHeight(100).WithBlockTime(blocktime).
		WithHeaderInfo(headerInfo)
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	params.UnstakeFeeRate = math.LegacyZeroDec()
	params.MinLiquidStakeAmount = math.NewInt(10000)
	err = s.keeper.SetParams(s.ctx, params)
	s.Require().NoError(err)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	stakingAmt := math.NewInt(10000000000000)
	s.fundAddr(s.delAddrs[0], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, stakingAmt)))
	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[5].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[6].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[7].String(), TargetWeight: math.NewInt(500)},
	}
	params.ModulePaused = false
	s.keeper.SetParams(s.ctx, params)
	reds := s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 0)

	stkXPRTMintAmt, err := s.keeper.LiquidStake(s.ctx, types.LiquidStakeProxyAcc, s.delAddrs[0], sdk.NewCoin(sdk.DefaultBondDenom, stakingAmt))
	s.Require().NoError(err)
	s.Require().Equal(stkXPRTMintAmt, stakingAmt)
	// assert rebalanced
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 0)
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)
	s.printRedelegationsLiquidTokens()

	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[5].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[6].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[7].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[8].String(), TargetWeight: math.NewInt(50)},
	}
	s.keeper.SetParams(s.ctx, params)
	headerInfo = s.ctx.HeaderInfo()
	headerInfo.Time = s.ctx.BlockTime().Add(time.Hour * 24)
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24)).
		WithHeaderInfo(headerInfo)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 8)
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)
	s.printRedelegationsLiquidTokens()
	// assert rebalanced
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 0)

	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[5].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[6].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[7].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[8].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[9].String(), TargetWeight: math.NewInt(500)},
	}
	s.keeper.SetParams(s.ctx, params)
	headerInfo = s.ctx.HeaderInfo()
	headerInfo.Time = s.ctx.BlockTime().Add(time.Hour * 24)
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24)).
		WithHeaderInfo(headerInfo)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 9)
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)
	s.printRedelegationsLiquidTokens()
	// assert rebalanced
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 0)

	// complete redelegations
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24 * 20).Add(time.Hour))
	_, err = s.app.StakingKeeper.EndBlocker(s.ctx)
	s.Require().NoError(err)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 0)
	// assert rebalanced
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)
	s.printRedelegationsLiquidTokens()

	// remove active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[5].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[6].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[7].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[8].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[9].String(), TargetWeight: math.NewInt(500)},
	}
	s.keeper.SetParams(s.ctx, params)
	headerInfo = s.ctx.HeaderInfo()
	headerInfo.Time = s.ctx.BlockTime().Add(time.Hour * 24)
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24)).
		WithHeaderInfo(headerInfo)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 9)
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)
	s.printRedelegationsLiquidTokens()
	// assert rebalanced
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 0)

	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[5].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[6].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[7].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[8].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[9].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[10].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[11].String(), TargetWeight: math.NewInt(500)},
		{ValidatorAddress: valOpers[12].String(), TargetWeight: math.NewInt(500)},
	}
	s.keeper.SetParams(s.ctx, params)
	headerInfo = s.ctx.HeaderInfo()
	headerInfo.Time = s.ctx.BlockTime().Add(time.Hour * 24)
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24)).WithHeaderInfo(headerInfo)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 11)
	// fail rebalancing due to redelegation hopping
	s.Require().Equal(s.redelegationsErrorCount(reds), 11)
	s.printRedelegationsLiquidTokens()

	// complete redelegation and retry
	s.completeRedelegationUnbonding()
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.printRedelegationsLiquidTokens()
	s.Require().Len(reds, 11)
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)

	// assert rebalanced
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 0)

	// modify weight
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(600)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(600)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(600)},
		{ValidatorAddress: valOpers[4].String(), TargetWeight: math.NewInt(600)},
		{ValidatorAddress: valOpers[5].String(), TargetWeight: math.NewInt(600)},
		{ValidatorAddress: valOpers[6].String(), TargetWeight: math.NewInt(600)},
		{ValidatorAddress: valOpers[7].String(), TargetWeight: math.NewInt(300)},
		{ValidatorAddress: valOpers[8].String(), TargetWeight: math.NewInt(300)},
		{ValidatorAddress: valOpers[9].String(), TargetWeight: math.NewInt(300)},
		{ValidatorAddress: valOpers[10].String(), TargetWeight: math.NewInt(300)},
		{ValidatorAddress: valOpers[11].String(), TargetWeight: math.NewInt(300)},
		{ValidatorAddress: valOpers[12].String(), TargetWeight: math.NewInt(300)},
	}
	err = s.keeper.SetParams(s.ctx, params)
	s.Require().NoError(err)
	headerInfo = s.ctx.HeaderInfo()
	headerInfo.Time = s.ctx.BlockTime().Add(time.Hour * 24)
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24)).WithHeaderInfo(headerInfo)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 6)
	// fail rebalancing partially due to redelegation hopping
	s.Require().Equal(s.redelegationsErrorCount(reds), 3)
	s.printRedelegationsLiquidTokens()

	// additional liquid stake when not rebalanced
	_, err = s.keeper.LiquidStake(s.ctx, types.LiquidStakeProxyAcc, s.delAddrs[0], sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1000000000)))
	s.Require().NoError(err)
	s.printRedelegationsLiquidTokens()

	// complete some redelegations
	headerInfo = s.ctx.HeaderInfo()
	headerInfo.Time = s.ctx.BlockTime().Add(time.Hour * 24 * 20).Add(time.Hour)
	s.ctx = s.ctx.WithBlockHeight(s.ctx.BlockHeight() + 100).WithBlockTime(s.ctx.BlockTime().Add(time.Hour * 24 * 20).Add(time.Hour)).WithHeaderInfo(headerInfo)
	_, err = s.app.StakingKeeper.EndBlocker(s.ctx)
	s.Require().NoError(err)
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 9)

	// failed redelegations with small amount (less than rebalancing trigger)
	s.Require().Equal(s.redelegationsErrorCount(reds), 6)
	s.printRedelegationsLiquidTokens()

	// assert rebalanced
	reds = s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.Require().Len(reds, 0)
	s.Require().Equal(s.redelegationsErrorCount(reds), 0)
	s.printRedelegationsLiquidTokens()
}

func (s *KeeperTestSuite) TestAutocompoundStakingRewards() {
	_, valOpers, _ := s.CreateValidators([]int64{2000000, 2000000, 2000000})
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)

	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(5000)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(5000)},
	}
	params.ModulePaused = false
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	stakingAmt := math.NewInt(100000000)
	s.Require().NoError(s.liquidStaking(s.delAddrs[0], stakingAmt))

	// no rewards
	totalRewards, totalDelShares, totalLiquidTokens, err := s.keeper.CheckDelegationStates(s.ctx, types.LiquidStakeProxyAcc)
	s.Require().NoError(err)
	s.EqualValues(totalRewards, math.LegacyZeroDec())
	s.EqualValues(totalDelShares, stakingAmt.ToLegacyDec(), totalLiquidTokens)

	// allocate rewards
	s.advanceHeight(360, false)
	totalRewards, totalDelShares, totalLiquidTokens, err = s.keeper.CheckDelegationStates(s.ctx, types.LiquidStakeProxyAcc)
	s.Require().NoError(err)
	s.NotEqualValues(totalRewards, math.LegacyZeroDec())
	s.Equal(totalLiquidTokens, stakingAmt)

	// withdraw rewards and re-staking
	whitelistedValsMap := types.GetWhitelistedValsMap(params.WhitelistedValidators)
	s.keeper.AutocompoundStakingRewards(s.ctx, whitelistedValsMap)
	totalRewardsAfter, totalDelSharesAfter, totalLiquidTokensAfter, err := s.keeper.CheckDelegationStates(s.ctx, types.LiquidStakeProxyAcc)
	s.Require().NoError(err)
	s.EqualValues(totalRewardsAfter, math.LegacyZeroDec())

	autocompoundFee := params.AutocompoundFeeRate.Mul(totalRewards).TruncateDec()
	s.EqualValues(totalDelSharesAfter, totalRewards.Sub(autocompoundFee).Add(totalDelShares).TruncateDec(), totalLiquidTokensAfter)

	stakingParams, err := s.app.StakingKeeper.GetParams(s.ctx)
	s.Require().NoError(err)
	feeAccountBalance := s.app.BankKeeper.GetBalance(
		s.ctx,
		sdk.MustAccAddressFromBech32(params.FeeAccountAddress),
		stakingParams.BondDenom,
	)
	s.EqualValues(autocompoundFee.TruncateInt(), feeAccountBalance.Amount)
}

func (s *KeeperTestSuite) TestLimitAutocompoundStakingRewards() {
	_, valOpers, _ := s.CreateValidators([]int64{2000000, 2000000, 2000000})
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)

	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(5000)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(5000)},
	}
	params.ModulePaused = false
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	stakingAmt := math.NewInt(100000000)
	s.Require().NoError(s.liquidStaking(s.delAddrs[0], stakingAmt))

	// allocate rewards
	s.advanceHeight(360, false)
	totalRewards, _, totalLiquidTokens, err := s.keeper.CheckDelegationStates(s.ctx, types.LiquidStakeProxyAcc)
	s.Require().NoError(err)
	s.NotEqualValues(totalRewards, math.LegacyZeroDec())
	s.Equal(totalLiquidTokens, stakingAmt)

	// unilaterally send tokens to the proxy account
	s.fundAddr(types.LiquidStakeProxyAcc, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(1000000000))))

	// withdraw rewards and re-stake
	whitelistedValsMap := types.GetWhitelistedValsMap(params.WhitelistedValidators)
	s.keeper.AutocompoundStakingRewards(s.ctx, whitelistedValsMap)

	// tokens still remaining in the proxy account as the balance was higher than the APY limit
	proxyAccBalanceAfter, err := s.keeper.GetProxyAccBalance(s.ctx, types.LiquidStakeProxyAcc)
	s.Require().NoError(err)
	s.NotEqual(proxyAccBalanceAfter.Amount, math.ZeroInt())
}

func (s *KeeperTestSuite) TestRemoveAllLiquidValidator() {
	_, valOpers, _ := s.CreateValidators([]int64{2000000, 2000000, 2000000})
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)

	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(2000)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(2000)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(2000)},
	}
	params.ModulePaused = false
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	stakingAmt := math.NewInt(100000000)
	s.Require().NoError(s.liquidStaking(s.delAddrs[0], stakingAmt))

	// allocate rewards
	s.advanceHeight(1, false)
	nasBefore, err := s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)
	s.Require().NotEqualValues(math.LegacyZeroDec(), nasBefore.TotalRemainingRewards)
	s.Require().NotEqualValues(math.LegacyZeroDec(), nasBefore.TotalDelShares)
	s.Require().NotEqualValues(math.LegacyZeroDec(), nasBefore.NetAmount)
	s.Require().NotEqualValues(math.ZeroInt(), nasBefore.TotalLiquidTokens)
	s.Require().EqualValues(math.ZeroInt(), nasBefore.ProxyAccBalance)

	// remove all whitelist
	params.WhitelistedValidators = []types.WhitelistedValidator{}
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	// no liquid validator
	lvs := s.keeper.GetAllLiquidValidators(s.ctx)
	s.Require().Len(lvs, 3) // now we do not remove inactive validators

	nasAfter, err := s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)

	s.Require().EqualValues(nasBefore.NetAmount.TruncateInt(), nasAfter.NetAmount.TruncateInt())

	s.completeRedelegationUnbonding()
	nasAfter2, err := s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(nasAfter.ProxyAccBalance, nasAfter2.ProxyAccBalance)                  // should be equal since no unbonding
	s.Require().EqualValues(nasBefore.NetAmount.TruncateInt(), nasAfter2.NetAmount.TruncateInt()) // should be equal since no unbonding
}

func (s *KeeperTestSuite) TestUndelegatedFundsNotBecomeFees() {
	_, valOpers, _ := s.CreateValidators([]int64{2000000, 2000000, 2000000, 2000000})
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	// configure validators
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(2000)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(2000)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(2000)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(2000)},
	}
	params.ModulePaused = false
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	// stake funds
	stakingAmt := math.NewInt(100000000)
	s.Require().NoError(s.liquidStaking(s.delAddrs[0], stakingAmt))

	// remove one validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(3000)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(3000)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(3000)},
	}
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	// unbonding should occur
	nas, err := s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)
	s.Require().NotEqual(nas.TotalUnbondingBalance, 0)

	// query fee account balance before unbonding finishes
	stakingParams, err := s.app.StakingKeeper.GetParams(s.ctx)
	s.Require().NoError(err)
	feeAccountBalance := s.app.BankKeeper.GetBalance(
		s.ctx,
		sdk.MustAccAddressFromBech32(params.FeeAccountAddress),
		stakingParams.BondDenom,
	)
	s.Require().Equal(math.ZeroInt(), feeAccountBalance.Amount)

	// complete unbondings
	s.completeRedelegationUnbonding()
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	// fee account has funds, but its from undelegated tokens
	feeAccountBalanceAfterUndelegation := s.app.BankKeeper.GetBalance(
		s.ctx,
		sdk.MustAccAddressFromBech32(params.FeeAccountAddress),
		stakingParams.BondDenom,
	)

	s.Require().Equal(math.ZeroInt(), feeAccountBalanceAfterUndelegation.Amount)
}
