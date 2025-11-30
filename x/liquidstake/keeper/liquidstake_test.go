package keeper_test

import (
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/types"
)

// tests LiquidStake, LiquidUnstake
func (s *KeeperTestSuite) TestLiquidStake() {
	_, valOpers, _ := s.CreateValidators([]int64{1000000, 2000000, 3000000})
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	params.MinLiquidStakeAmount = math.NewInt(50000)
	params.ModulePaused = false
	err = s.keeper.SetParams(s.ctx, params)
	s.Require().NoError(err)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	stakingAmt := params.MinLiquidStakeAmount

	// fail, no active validator
	cachedCtx, _ := s.ctx.CacheContext()
	stkXPRTMintAmt, err := s.keeper.LiquidStake(
		cachedCtx, types.LiquidStakeProxyAcc, s.delAddrs[0],
		sdk.NewCoin(sdk.DefaultBondDenom, stakingAmt),
	)
	s.Require().ErrorIs(err, types.ErrActiveLiquidValidatorsNotExists)
	s.Require().Equal(stkXPRTMintAmt, math.ZeroInt())

	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(2000)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(2000)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(2000)},
	}
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	res := s.keeper.GetAllLiquidValidatorStates(s.ctx)
	s.Require().Equal(params.WhitelistedValidators[0].ValidatorAddress,
		res[0].OperatorAddress)
	s.Require().Equal(params.WhitelistedValidators[0].TargetWeight,
		res[0].Weight)
	s.Require().Equal(types.ValidatorStatusActive, res[0].Status)
	s.Require().Equal(math.LegacyZeroDec(), res[0].DelShares)
	s.Require().Equal(math.ZeroInt(), res[0].LiquidTokens)

	s.Require().Equal(params.WhitelistedValidators[1].ValidatorAddress,
		res[1].OperatorAddress)
	s.Require().Equal(params.WhitelistedValidators[1].TargetWeight,
		res[1].Weight)
	s.Require().Equal(types.ValidatorStatusActive, res[1].Status)
	s.Require().Equal(math.LegacyZeroDec(), res[1].DelShares)
	s.Require().Equal(math.ZeroInt(), res[1].LiquidTokens)

	s.Require().Equal(params.WhitelistedValidators[2].ValidatorAddress,
		res[2].OperatorAddress)
	s.Require().Equal(params.WhitelistedValidators[2].TargetWeight,
		res[2].Weight)
	s.Require().Equal(types.ValidatorStatusActive, res[2].Status)
	s.Require().Equal(math.LegacyZeroDec(), res[2].DelShares)
	s.Require().Equal(math.ZeroInt(), res[2].LiquidTokens)

	// liquid stake
	stkXPRTMintAmt, err = s.keeper.LiquidStake(
		s.ctx, types.LiquidStakeProxyAcc, s.delAddrs[0],
		sdk.NewCoin(sdk.DefaultBondDenom, stakingAmt),
	)
	s.Require().NoError(err)
	s.Require().Equal(stkXPRTMintAmt, stakingAmt)

	_, err = s.app.StakingKeeper.GetDelegation(
		s.ctx, s.delAddrs[0], valOpers[0],
	)
	s.Require().Error(err)
	_, err = s.app.StakingKeeper.GetDelegation(
		s.ctx, s.delAddrs[0], valOpers[1],
	)
	s.Require().Error(err)
	_, err = s.app.StakingKeeper.GetDelegation(
		s.ctx, s.delAddrs[0], valOpers[2],
	)
	s.Require().Error(err)

	proxyAccDel1, err := s.app.StakingKeeper.GetDelegation(
		s.ctx, types.LiquidStakeProxyAcc, valOpers[0],
	)
	s.Require().NoError(err)
	proxyAccDel2, err := s.app.StakingKeeper.GetDelegation(
		s.ctx, types.LiquidStakeProxyAcc, valOpers[1],
	)
	s.Require().NoError(err)
	proxyAccDel3, err := s.app.StakingKeeper.GetDelegation(
		s.ctx, types.LiquidStakeProxyAcc, valOpers[2],
	)
	s.Require().NoError(err)
	s.Require().Equal(proxyAccDel1.Shares, math.LegacyNewDec(16668))
	s.Require().Equal(proxyAccDel2.Shares, math.LegacyNewDec(16666))
	s.Require().Equal(proxyAccDel2.Shares, math.LegacyNewDec(16666))
	s.Require().Equal(stakingAmt.ToLegacyDec(),
		proxyAccDel1.Shares.Add(proxyAccDel2.Shares).Add(proxyAccDel3.Shares))

	liquidBondDenom, err := s.keeper.LiquidBondDenom(s.ctx)
	s.Require().NoError(err)
	balanceBeforeUBD := s.app.BankKeeper.GetBalance(
		s.ctx, s.delAddrs[0], sdk.DefaultBondDenom,
	)
	s.Require().Equal(balanceBeforeUBD.Amount, math.NewInt(999950000))
	ubdStkXPRT := sdk.NewCoin(liquidBondDenom, math.NewInt(10000))
	stkXPRTBalance := s.app.BankKeeper.GetBalance(
		s.ctx, s.delAddrs[0], liquidBondDenom,
	)
	stkXPRTTotalSupply := s.app.BankKeeper.GetSupply(
		s.ctx, liquidBondDenom,
	)
	s.Require().Equal(stkXPRTBalance,
		sdk.NewCoin(liquidBondDenom, math.NewInt(50000)))
	s.Require().Equal(stkXPRTBalance, stkXPRTTotalSupply)

	// liquid unstaking
	ubdTime, unbondingAmt, ubds, unbondedAmt, err := s.keeper.LiquidUnstake(
		s.ctx, types.LiquidStakeProxyAcc, s.delAddrs[0], ubdStkXPRT,
	)
	s.Require().NoError(err)
	s.Require().EqualValues(unbondedAmt, math.ZeroInt())
	s.Require().Len(ubds, 3)

	// crumb excepted on unbonding
	crumb := ubdStkXPRT.Amount.Sub(ubdStkXPRT.Amount.QuoRaw(3).MulRaw(3))
	s.Require().EqualValues(unbondingAmt, ubdStkXPRT.Amount.Sub(crumb))
	s.Require().Equal(ubds[0].DelegatorAddress, s.delAddrs[0].String())
	blocktime, err := sdk.ParseTime("2022-03-22T00:00:00.000000000")
	s.Require().NoError(err)
	s.Require().Equal(ubdTime, blocktime)
	stkXPRTBalanceAfter := s.app.BankKeeper.GetBalance(
		s.ctx, s.delAddrs[0], liquidBondDenom,
	)
	s.Require().Equal(stkXPRTBalanceAfter,
		sdk.NewCoin(liquidBondDenom, math.NewInt(40000)))

	balanceBeginUBD := s.app.BankKeeper.GetBalance(
		s.ctx, s.delAddrs[0], sdk.DefaultBondDenom,
	)
	s.Require().Equal(balanceBeginUBD.Amount, balanceBeforeUBD.Amount)

	proxyAccDel1, err = s.app.StakingKeeper.GetDelegation(
		s.ctx, types.LiquidStakeProxyAcc, valOpers[0],
	)
	s.Require().NoError(err)
	proxyAccDel2, err = s.app.StakingKeeper.GetDelegation(
		s.ctx, types.LiquidStakeProxyAcc, valOpers[1],
	)
	s.Require().NoError(err)
	proxyAccDel3, err = s.app.StakingKeeper.GetDelegation(
		s.ctx, types.LiquidStakeProxyAcc, valOpers[2],
	)
	s.Require().NoError(err)
	s.Require().Equal(stakingAmt.Sub(unbondingAmt).ToLegacyDec(),
		proxyAccDel1.GetShares().Add(proxyAccDel2.Shares).Add(proxyAccDel3.Shares))

	// complete unbonding
	s.ctx = s.ctx.WithBlockHeight(200).WithBlockTime(ubdTime.Add(1))
	updates, err := s.app.StakingKeeper.BlockValidatorUpdates(s.ctx)
	s.Require().NoError(err)
	s.Require().Empty(updates)
	balanceCompleteUBD := s.app.BankKeeper.GetBalance(
		s.ctx, s.delAddrs[0], sdk.DefaultBondDenom,
	)
	s.Require().Equal(balanceCompleteUBD.Amount,
		balanceBeforeUBD.Amount.Add(unbondingAmt))

	proxyAccDel1, err = s.app.StakingKeeper.GetDelegation(
		s.ctx, types.LiquidStakeProxyAcc, valOpers[0],
	)
	s.Require().NoError(err)
	proxyAccDel2, err = s.app.StakingKeeper.GetDelegation(
		s.ctx, types.LiquidStakeProxyAcc, valOpers[1],
	)
	s.Require().NoError(err)
	proxyAccDel3, err = s.app.StakingKeeper.GetDelegation(
		s.ctx, types.LiquidStakeProxyAcc, valOpers[2],
	)
	s.Require().NoError(err)
	s.Require().Equal(math.LegacyNewDec(13335), proxyAccDel1.Shares)
	s.Require().Equal(math.LegacyNewDec(13333), proxyAccDel2.Shares)
	s.Require().Equal(math.LegacyNewDec(13333), proxyAccDel3.Shares)

	res = s.keeper.GetAllLiquidValidatorStates(s.ctx)
	s.Require().Equal(params.WhitelistedValidators[0].ValidatorAddress,
		res[0].OperatorAddress)
	s.Require().Equal(params.WhitelistedValidators[0].TargetWeight,
		res[0].Weight)
	s.Require().Equal(types.ValidatorStatusActive, res[0].Status)
	s.Require().Equal(math.LegacyNewDec(13335), res[0].DelShares)

	s.Require().Equal(params.WhitelistedValidators[1].ValidatorAddress,
		res[1].OperatorAddress)
	s.Require().Equal(params.WhitelistedValidators[1].TargetWeight,
		res[1].Weight)
	s.Require().Equal(types.ValidatorStatusActive, res[1].Status)
	s.Require().Equal(math.LegacyNewDec(13333), res[1].DelShares)

	s.Require().Equal(params.WhitelistedValidators[2].ValidatorAddress,
		res[2].OperatorAddress)
	s.Require().Equal(params.WhitelistedValidators[2].TargetWeight,
		res[2].Weight)
	s.Require().Equal(types.ValidatorStatusActive, res[2].Status)
	s.Require().Equal(math.LegacyNewDec(13333), res[2].DelShares)

	// rewards are not autocompounded after validator set update and rebalancing
	s.advanceHeight(10, true)
	rewards, totalLiquidShares, _, err := s.keeper.CheckDelegationStates(
		s.ctx, types.LiquidStakeProxyAcc,
	)
	s.Require().NoError(err)
	s.Require().NotEqualValues(rewards, math.LegacyZeroDec())
	s.Require().EqualValues(totalLiquidShares, proxyAccDel1.Shares.Add(proxyAccDel2.Shares).Add(proxyAccDel3.Shares))

	// all remaining rewards re-staked, request last unstaking, unbond all
	params, err = s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	s.keeper.AutocompoundStakingRewards(s.ctx, types.GetWhitelistedValsMap(params.WhitelistedValidators))
	stkxprtBalanceBefore := s.app.BankKeeper.GetBalance(s.ctx, s.delAddrs[0], params.LiquidBondDenom).Amount
	rewards, _, _, err = s.keeper.CheckDelegationStates(
		s.ctx, types.LiquidStakeProxyAcc,
	)
	s.Require().NoError(err)
	s.Require().EqualValues(rewards, math.LegacyZeroDec())
	s.Require().NoError(
		s.liquidUnstaking(s.delAddrs[0], stkxprtBalanceBefore, true),
	)

	// still active liquid validator after unbond all
	alv := s.keeper.GetActiveLiquidValidators(
		s.ctx, params.WhitelistedValsMap(),
	)
	s.Require().True(len(alv) != 0)

	// no btoken supply and netAmount after unbond all
	nas, err := s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(nas.StkxprtTotalSupply, math.ZeroInt())
	s.Require().Equal(nas.TotalRemainingRewards, math.LegacyZeroDec())
	s.Require().Equal(nas.TotalDelShares, math.LegacyZeroDec())
	s.Require().Equal(nas.TotalLiquidTokens, math.ZeroInt())
	s.Require().Equal(nas.ProxyAccBalance, math.ZeroInt())
	s.Require().Equal(nas.NetAmount, math.LegacyZeroDec())
}

func (s *KeeperTestSuite) TestLiquidStakeFromVestingAccount() {
	_, valOpers, _ := s.CreateValidators([]int64{1000000, 2000000, 3000000})
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(3334)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(3333)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(3333)},
	}
	params.ModulePaused = false
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	from := s.delAddrs[0]
	vestingAmt := s.app.BankKeeper.GetAllBalances(s.ctx, from)
	vestingStartTime := s.ctx.BlockTime().Add(1 * time.Hour)
	vestingEndTime := s.ctx.BlockTime().Add(2 * time.Hour)
	vestingMidTime := s.ctx.BlockTime().Add(90 * time.Minute)

	vestingAccAddr := "persistence10n3ncmlsaqfuwsmfll8kq6hvt4x7c8czahev75"
	vestingAcc, err := sdk.AccAddressFromBech32(vestingAccAddr)
	s.Require().NoError(err)

	// createContinuousVestingAccount
	cVestingAcc := s.createContinuousVestingAccount(from, vestingAcc, vestingAmt, vestingStartTime, vestingEndTime)
	spendableCoins := s.app.BankKeeper.SpendableCoins(s.ctx, cVestingAcc.GetAddress())
	s.Require().True(spendableCoins.IsZero())
	lockedCoins := s.app.BankKeeper.LockedCoins(s.ctx, cVestingAcc.GetAddress())
	s.Require().EqualValues(lockedCoins, vestingAmt)

	// failed liquid stake, no spendable coins on the vesting account ( not allowed locked coins )
	err = s.liquidStaking(vestingAcc, vestingAmt.AmountOf(sdk.DefaultBondDenom))
	s.Require().ErrorIs(err, sdkerrors.ErrInsufficientFunds)

	// release some vesting coins
	s.ctx = s.ctx.WithBlockTime(vestingMidTime)
	spendableCoins = s.app.BankKeeper.SpendableCoins(s.ctx, cVestingAcc.GetAddress())
	s.Require().True(spendableCoins.IsAllPositive())
	lockedCoins = s.app.BankKeeper.LockedCoins(s.ctx, cVestingAcc.GetAddress())
	s.Require().True(lockedCoins.IsAllPositive())

	// success with released spendable coins
	err = s.liquidStaking(vestingAcc, spendableCoins.AmountOf(sdk.DefaultBondDenom))
	s.Require().NoError(err)
	nas, err := s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)
	s.Require().EqualValues(nas.TotalLiquidTokens, spendableCoins.AmountOf(sdk.DefaultBondDenom))
}

func (s *KeeperTestSuite) TestLiquidStakeEdgeCases() {
	_, valOpers, _ := s.CreateValidators([]int64{1000000, 2000000, 3000000})
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	stakingAmt := math.NewInt(5000000)

	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(3334)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(3333)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(3333)},
	}
	params.ModulePaused = false
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	// fail Invalid BondDenom case
	_, err = s.keeper.LiquidStake(s.ctx, types.LiquidStakeProxyAcc, s.delAddrs[0], sdk.NewCoin("bad", stakingAmt))
	s.Require().ErrorIs(err, types.ErrInvalidBondDenom)

	// liquid stake, unstaking with huge amount
	hugeAmt := math.NewInt(1_000_000_000_000_000_000)
	s.fundAddr(s.delAddrs[0], sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, hugeAmt.MulRaw(2))))
	s.Require().NoError(s.liquidStaking(s.delAddrs[0], hugeAmt))
	s.Require().NoError(s.liquidStaking(s.delAddrs[0], hugeAmt))
	s.Require().NoError(s.liquidUnstaking(s.delAddrs[0], math.NewInt(10), true))
	s.Require().NoError(s.liquidUnstaking(s.delAddrs[0], hugeAmt, true))
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	s.completeRedelegationUnbonding()
	states, err := s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)
	states.TotalLiquidTokens.Equal(hugeAmt)
}

func (s *KeeperTestSuite) TestLiquidUnstakeEdgeCases() {
	mintParams, err := s.app.MintKeeper.Params.Get(s.ctx)
	s.Require().NoError(err)
	mintParams.InflationMax = math.LegacyNewDec(0)
	mintParams.InflationMin = math.LegacyNewDec(0)
	mintParams.InflationRateChange = math.LegacyNewDec(0)
	err = s.app.MintKeeper.Params.Set(s.ctx, mintParams)
	s.Require().NoError(err)

	_, valOpers, _ := s.CreateValidators([]int64{1000000, 2000000, 3000000})
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)
	stakingAmt := math.NewInt(5000000)

	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(3334)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(3333)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(3333)},
	}
	params.ModulePaused = false
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	// success liquid stake
	s.Require().NoError(s.liquidStaking(s.delAddrs[0], stakingAmt))

	// fail when liquid unstaking with too small amount
	_, _, _, _, err = s.liquidUnstakingWithResult(s.delAddrs[0], sdk.NewCoin(params.LiquidBondDenom, math.NewInt(2)))
	s.Require().ErrorIs(err, types.ErrTooSmallLiquidUnstakingAmount)

	// fail when liquid unstaking with zero amount
	_, _, _, _, err = s.liquidUnstakingWithResult(s.delAddrs[0], sdk.NewCoin(params.LiquidBondDenom, math.NewInt(0)))
	s.Require().ErrorIs(err, types.ErrTooSmallLiquidUnstakingAmount)

	// fail when invalid liquid bond denom
	_, _, _, _, err = s.liquidUnstakingWithResult(s.delAddrs[0], sdk.NewCoin("stake", math.NewInt(10000)))
	s.Require().ErrorIs(err, types.ErrInvalidLiquidBondDenom)

	// verify that there is no problem performing liquid unstaking as much as the MaxEntries
	stakingParams, err := s.app.StakingKeeper.GetParams(s.ctx)
	s.Require().NoError(err)
	for i := uint32(0); i < stakingParams.MaxEntries; i++ {
		s.Require().NoError(s.liquidUnstaking(s.delAddrs[0], math.NewInt(1000), false))
	}

	// on sdk 0.47+ shouldn't fail in an attempt to go beyond MaxEntries
	err = s.liquidUnstaking(s.delAddrs[0], math.NewInt(1000), false)
	s.Require().NoError(err)

	dels, err := s.app.StakingKeeper.GetUnbondingDelegations(s.ctx, s.delAddrs[0], 100)
	s.Require().NoError(err)
	for _, ubd := range dels {
		s.Require().EqualValues(1, len(ubd.Entries))
	}

	// set empty whitelisted, active liquid validator
	params.WhitelistedValidators = []types.WhitelistedValidator{}
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	// error case where there is a quantity that are unbonding balance or remaining rewards that is not re-stake or withdrawn in netAmount.
	// NOT APPLICABLE since we do not validator unbond if validator goes inactive.
	//_, _, _, _, err = s.liquidUnstakingWithResult(s.delAddrs[0], sdk.NewCoin(params.LiquidBondDenom, math.NewInt(1000)))
	//s.Require().ErrorIs(err, types.ErrInsufficientProxyAccBalance)

	// success after complete unbonding, Not applicable
	s.completeRedelegationUnbonding()
	// ubdTime, unbondingAmt, ubds, unbondedAmt, err := s.liquidUnstakingWithResult(s.delAddrs[0], sdk.NewCoin(params.LiquidBondDenom, math.NewInt(1000)))
	// s.Require().NoError(err)
	// s.Require().EqualValues(unbondedAmt, math.NewInt(1000))
	// s.Require().EqualValues(unbondingAmt, math.ZeroInt())
	// s.Require().EqualValues(ubdTime, time.Time{})
	// s.Require().Len(ubds, 0)
}

func (s *KeeperTestSuite) TestShareInflation() {
	_, valOpers, _ := s.CreateValidators([]int64{1000000, 2000000, 3000000, 4000000})
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	// set minimum amount and unstake fee to 0 for testing
	params.MinLiquidStakeAmount = math.NewInt(0)
	params.UnstakeFeeRate = math.LegacyZeroDec()
	err = s.keeper.SetParams(s.ctx, params)
	s.Require().NoError(err)

	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(2500)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(2500)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(2500)},
		{ValidatorAddress: valOpers[3].String(), TargetWeight: math.NewInt(2500)},
	}
	params.ModulePaused = false
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	initialStakingAmt := math.NewInt(1)          // little amount
	initializingStakingAmt := math.NewInt(10000) // normal amount
	attacker := s.delAddrs[0]
	user := s.delAddrs[1]
	protocol := s.delAddrs[3]

	// 0. [a solution?] be first depositor
	mintAmount0, err := s.keeper.LiquidStake(s.ctx, types.LiquidStakeProxyAcc,
		protocol, sdk.NewCoin(sdk.DefaultBondDenom, initializingStakingAmt))
	s.Require().NoError(err)
	s.Require().Equal(mintAmount0, initializingStakingAmt)

	// 1. attacker becomes first depositor and liquid stake
	mintAmount, err := s.keeper.LiquidStake(s.ctx, types.LiquidStakeProxyAcc,
		attacker, sdk.NewCoin(sdk.DefaultBondDenom, initialStakingAmt))
	s.Require().NoError(err)
	s.Require().Equal(mintAmount, initialStakingAmt)

	// 2. The user sends a liquid stake message, but their tx got front-run by the attacker
	// ideally, the user should get 1000 stkXPRT (1 * 1000 / 1)
	// stkXPRT to mint = stkXPRT supply * sent XPRT / total XPRT
	userStakeAmount := math.NewInt(1_000)

	// 3. attacker's tx got accepted first which sends funds directly to proxy account
	attackerTransferAmount := userStakeAmount.Quo(math.NewInt(2))
	err = s.app.BankKeeper.SendCoins(s.ctx, attacker, types.LiquidStakeProxyAcc,
		sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, attackerTransferAmount)))
	s.Require().NoError(err)

	// 4. user tx went through and the mint rate is not affected by the XPRT sent by the attacker
	// stkXPRT to mint = 1 * 1000 / 1 = 1
	mintAmount, err = s.keeper.LiquidStake(s.ctx, types.LiquidStakeProxyAcc, user, sdk.NewCoin(sdk.DefaultBondDenom, userStakeAmount))
	s.Require().NoError(err)
	s.Require().Equal(mintAmount, math.NewInt(1_000))

	// 5. attacker unstakes the shares immediately
	liquidBondDenom, err := s.keeper.LiquidBondDenom(s.ctx)
	s.Require().NoError(err)
	_, unbondingAmt, _, _, err := s.keeper.LiquidUnstake(s.ctx, types.LiquidStakeProxyAcc, attacker, sdk.NewCoin(liquidBondDenom, math.NewInt(1)))
	// s.Require().NoError(err)
	s.Require().ErrorContains(err, "liquid unstaking amount is too small")

	attackerProfit := unbondingAmt.Sub(initialStakingAmt).Sub(attackerTransferAmount)
	s.Require().LessOrEqual(attackerProfit.Int64(), math.ZeroInt().Int64())
}
