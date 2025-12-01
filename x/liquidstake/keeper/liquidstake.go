package keeper

import (
	"encoding/json"
	"sort"
	"time"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	liquidtypes "github.com/cosmos/gaia/v25/x/liquid/types"

	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/types"
)

func (k Keeper) LiquidBondDenom(ctx sdk.Context) (string, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return "", err
	}
	return params.LiquidBondDenom, nil
}

// GetNetAmountState calculates the sum of bondedDenom balance, total delegation tokens(slash applied LiquidTokens), total remaining reward of types.LiquidStakeProxyAcc
// During liquid unstaking, stkxprt immediately burns and the unbonding queue belongs to the requester, so the liquid staker's unbonding values are excluded on netAmount
// It is used only for calculation and query and is not stored in kv.
func (k Keeper) GetNetAmountState(ctx sdk.Context) (types.NetAmountState, error) {
	totalRemainingRewards, totalDelShares, totalLiquidTokens, err := k.CheckDelegationStates(ctx, types.LiquidStakeProxyAcc)
	if err != nil {
		return types.NetAmountState{}, err
	}
	proxyAccBalance, err := k.GetProxyAccBalance(ctx, types.LiquidStakeProxyAcc)
	if err != nil {
		return types.NetAmountState{}, err
	}
	liquidBondDenom, err := k.LiquidBondDenom(ctx)
	if err != nil {
		return types.NetAmountState{}, err
	}
	nas := types.NetAmountState{
		StkxprtTotalSupply:    k.bankKeeper.GetSupply(ctx, liquidBondDenom).Amount,
		TotalDelShares:        totalDelShares,
		TotalLiquidTokens:     totalLiquidTokens,
		TotalRemainingRewards: totalRemainingRewards,
		ProxyAccBalance:       proxyAccBalance.Amount,
	}
	totalUnbondingBalance := math.ZeroInt()
	ubds, err := k.stakingKeeper.GetAllUnbondingDelegations(ctx, types.LiquidStakeProxyAcc)
	if err != nil {
		return types.NetAmountState{}, err
	}
	for _, ubd := range ubds {
		for _, entry := range ubd.Entries {
			// use Balance(slashing applied) not InitialBalance(without slashing)
			totalUnbondingBalance = totalUnbondingBalance.Add(entry.Balance)
		}
	}

	nas.TotalUnbondingBalance = totalUnbondingBalance
	nas.NetAmount = nas.CalcNetAmount()
	nas.MintRate = nas.CalcMintRate()
	return nas, nil
}

// LiquidStake mints stkXPRT worth of staking coin value according to NetAmount and performs LiquidDelegate.
func (k Keeper) LiquidStake(
	ctx sdk.Context, proxyAcc, liquidStaker sdk.AccAddress, stakingCoin sdk.Coin,
) (stkXPRTMintAmount math.Int, err error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}

	if params.ModulePaused {
		return math.ZeroInt(), types.ErrModulePaused
	}

	// check minimum liquid stake amount
	if stakingCoin.Amount.LT(params.MinLiquidStakeAmount) {
		return math.ZeroInt(), types.ErrLessThanMinLiquidStakeAmount
	}

	// check bond denomination
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}
	if stakingCoin.Denom != bondDenom {
		return math.ZeroInt(), errorsmod.Wrapf(
			types.ErrInvalidBondDenom, "invalid coin denomination: got %s, expected %s", stakingCoin.Denom, bondDenom,
		)
	}

	whitelistedValsMap := types.GetWhitelistedValsMap(params.WhitelistedValidators)
	activeVals, err := k.GetActiveLiquidValidators(ctx, whitelistedValsMap)
	if err != nil {
		return math.ZeroInt(), err
	}

	if activeVals.Len() == 0 {
		return math.ZeroInt(), types.ErrActiveLiquidValidatorsNotExists
	}

	totalActiveWeight := activeVals.TotalWeight(whitelistedValsMap)
	activeWeightQuorum := math.LegacyNewDecFromInt(totalActiveWeight).Quo(
		math.LegacyNewDecFromInt(types.TotalValidatorWeight),
	)
	if activeWeightQuorum.LT(types.ActiveLiquidValidatorsWeightQuorum) {
		k.Logger(ctx).Error(
			"active liquid validators weight quorum not reached",
			types.ActiveWeightQuorumKeyVal,
			activeWeightQuorum.String(),
			types.MinActiveWeightQuorumKeyVal,
			types.ActiveLiquidValidatorsWeightQuorum.String(),
		)

		return math.ZeroInt(), errorsmod.Wrapf(
			types.ErrActiveLiquidValidatorsWeightQuorumNotReached, "%s < %s",
			activeWeightQuorum.String(), types.ActiveLiquidValidatorsWeightQuorum.String(),
		)
	}

	// NetAmount must be calculated before send
	nas, err := k.GetNetAmountState(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}

	// send staking coin to liquid stake proxy account to proxy delegation, need sufficient spendable balances
	err = k.bankKeeper.SendCoins(ctx, liquidStaker, proxyAcc, sdk.NewCoins(stakingCoin))
	if err != nil {
		return math.ZeroInt(), err
	}

	// mint stkxprt, MintAmount = TotalSupply * StakeAmount/NetAmount
	liquidBondDenom, err := k.LiquidBondDenom(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}
	stkXPRTMintAmount = stakingCoin.Amount

	if nas.StkxprtTotalSupply.IsPositive() {
		if nas.NetAmount.IsZero() {
			// this case must not be reachable, consider stopping module for investigation
			// c_value -> inf
			k.Logger(ctx).Error(
				"infinite c value",
				types.NetAmountStateKeyVal,
				nas.String(),
			)

			return math.ZeroInt(), types.ErrInsufficientProxyAccBalance
		}

		stkXPRTMintAmount = types.NativeTokenToStkXPRT(stakingCoin.Amount, nas.StkxprtTotalSupply, nas.NetAmount)
	}

	if !stkXPRTMintAmount.IsPositive() {
		return math.ZeroInt(), types.ErrTooSmallLiquidStakeAmount
	}

	// mint on module acc and send
	mintCoin := sdk.NewCoins(sdk.NewCoin(liquidBondDenom, stkXPRTMintAmount))
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoin)
	if err != nil {
		return stkXPRTMintAmount, err
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, liquidStaker, mintCoin)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to send minted coins to liquid staker",
			types.ErrorKeyVal,
			err,
		)

		return stkXPRTMintAmount, err
	}

	err = k.LiquidDelegate(ctx, proxyAcc, activeVals, stakingCoin.Amount, whitelistedValsMap)
	return stkXPRTMintAmount, err
}

// LockOnLP sends tokens to a CW contract (Superfluid LP) with time locking.
// It performs a CosmWasm execution through global message handler and may fail.
// Emits events on a successful call.
func (k Keeper) LockOnLP(ctx sdk.Context, delegator sdk.AccAddress, amount sdk.Coin) ([]*codectypes.Any, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	if len(params.CwLockedPoolAddress) == 0 {
		return nil, types.ErrNoLPContractAddress
	} else if amount.Denom != params.LiquidBondDenom {
		return nil, types.ErrInvalidDenom.Wrapf("cannot lock any denom on LP except liquid bond denom: %s", params.LiquidBondDenom)
	}

	msg := &LockLstAssetMsg{
		Asset: Asset{
			Amount: amount.Amount.String(),
			Info: AssetInfo{
				NativeToken: NativeTokenInfo{
					Denom: amount.Denom,
				},
			},
		},
	}

	callData, err := json.Marshal(&ExecMsg{
		LockLstAsset: msg,
	})
	if err != nil {
		panic("failed to marshal CW contract call LockLstAsset")
	}

	cwMsg := &wasmtypes.MsgExecuteContract{
		Sender:   delegator.String(),
		Contract: params.CwLockedPoolAddress,
		Msg:      wasmtypes.RawContractMessage(callData),
		Funds:    sdk.NewCoins(amount),
	}

	handler := k.router.Handler(cwMsg)
	if handler == nil {
		k.Logger(ctx).Error("failed to find CW contract handler")

		return nil, sdkerrors.ErrUnknownRequest.Wrapf("unrecognized message route: %s", sdk.MsgTypeURL(cwMsg))
	}

	msgResp, err := handler(ctx, cwMsg)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to execute CW contract message",
			types.ErrorKeyVal,
			err,
			types.MsgKeyVal,
			cwMsg.String(),
		)

		return nil, types.ErrLPContract.Wrapf("error: %s, message %v", err.Error(), cwMsg)
	}

	// emit the events from the dispatched actions
	ctx.EventManager().EmitEvents(msgResp.GetEvents())

	return msgResp.MsgResponses, nil
}

type ExecMsg struct {
	LockLstAsset *LockLstAssetMsg `json:"lock_lst_asset,omitempty"`
}

type LockLstAssetMsg struct {
	Asset Asset `json:"asset"`
}

type Asset struct {
	Amount string    `json:"amount"`
	Info   AssetInfo `json:"info"`
}

type AssetInfo struct {
	NativeToken NativeTokenInfo `json:"native_token"`
}

type NativeTokenInfo struct {
	Denom string `json:"denom"`
}

// DelegateWithCap is a wrapper to invoke stakingKeeper.Delegate but account for
// the amount of liquid staked shares and check against liquid staking cap.
func (k Keeper) DelegateWithCap(
	ctx sdk.Context,
	delegatorAddress sdk.AccAddress,
	validator stakingtypes.Validator,
	bondAmt math.Int,
) error {
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return err
	}
	msgDelegate := &stakingtypes.MsgDelegate{
		DelegatorAddress: delegatorAddress.String(),
		ValidatorAddress: validator.OperatorAddress,
		Amount:           sdk.NewCoin(bondDenom, bondAmt),
	}
	handler := k.router.Handler(msgDelegate)
	if handler == nil {
		k.Logger(ctx).Error("failed to find delegate handler")

		return errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message route: %s", sdk.MsgTypeURL(msgDelegate))
	}
	res, err := handler(ctx, msgDelegate)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to execute delegate msg",
			types.ErrorKeyVal,
			err,
			types.MsgKeyVal,
			msgDelegate.String(),
		)

		return errorsmod.Wrapf(types.ErrDelegationFailed, "failed to send delegate msg with err: %v", err)
	}
	ctx.EventManager().EmitEvents(res.GetEvents())

	if len(res.MsgResponses) != 1 {
		return errorsmod.Wrapf(
			types.ErrInvalidResponse,
			"expected msg response should be exactly 1, got: %v, responses: %v",
			len(res.MsgResponses), res.MsgResponses,
		)
	}

	var msgDelegateResponse stakingtypes.MsgDelegateResponse
	if err = k.cdc.Unmarshal(res.MsgResponses[0].Value, &msgDelegateResponse); err != nil {
		return errorsmod.Wrapf(
			sdkerrors.ErrJSONUnmarshal,
			"cannot unmarshal delegate tx response message: %v",
			err,
		)
	}

	return nil
}

// UnbondWithCap is a wrapper to invoke stakingKeeper.Unbond but updates
// the total liquid staked tokens.
func (k Keeper) UnbondWithCap(
	ctx sdk.Context,
	delegatorAddress sdk.AccAddress,
	validatorAddress sdk.ValAddress,
	amount sdk.Coin,
	userAddress sdk.AccAddress,
) (math.Int, error) {
	// perform an LSM tokenize->bank send->redeem flow: moving delegation from proxyAcc onto user's account
	lsmTokenizeMsg := &liquidtypes.MsgTokenizeShares{
		DelegatorAddress:    delegatorAddress.String(),
		ValidatorAddress:    validatorAddress.String(),
		Amount:              amount,
		TokenizedShareOwner: userAddress.String(),
	}

	handler := k.router.Handler(lsmTokenizeMsg)
	if handler == nil {
		k.Logger(ctx).Error("failed to find tokenize handler")

		return math.ZeroInt(), sdkerrors.ErrUnknownRequest.Wrapf("unrecognized message route: %s", sdk.MsgTypeURL(lsmTokenizeMsg))
	}

	// [1] tokenize delegation into LSM shares
	msgResp, err := handler(ctx, lsmTokenizeMsg)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to execute tokenize shares message",
			types.ErrorKeyVal,
			err,
			types.MsgKeyVal,
			lsmTokenizeMsg.String(),
		)

		return math.ZeroInt(), types.ErrLSMTokenizeFailed.Wrapf("error: %s; message: %v", err.Error(), lsmTokenizeMsg)
	}
	ctx.EventManager().EmitEvents(msgResp.GetEvents())

	if len(msgResp.MsgResponses) != 1 {
		return math.ZeroInt(), errorsmod.Wrapf(
			types.ErrInvalidResponse,
			"expected msg response should be exactly 1, got: %v, responses: %v",
			len(msgResp.MsgResponses), msgResp.MsgResponses,
		)
	}

	var lsmTokenizeResp liquidtypes.MsgTokenizeSharesResponse
	if err = k.cdc.Unmarshal(msgResp.MsgResponses[0].Value, &lsmTokenizeResp); err != nil {
		return math.ZeroInt(), errorsmod.Wrapf(
			sdkerrors.ErrJSONUnmarshal,
			"cannot unmarshal tokenize share tx response message: %v",
			err,
		)
	}

	// [2] send LSM shares to proxyAcc
	err = k.bankKeeper.SendCoins(ctx, delegatorAddress, userAddress, sdk.NewCoins(lsmTokenizeResp.Amount))
	if err != nil {
		return math.ZeroInt(), err
	}

	lsmRedeemMsg := &liquidtypes.MsgRedeemTokensForShares{
		DelegatorAddress: userAddress.String(),
		Amount:           lsmTokenizeResp.Amount,
	}

	handler = k.router.Handler(lsmRedeemMsg)
	if handler == nil {
		k.Logger(ctx).Error("failed to find redeem handler")
		return math.ZeroInt(), sdkerrors.ErrUnknownRequest.Wrapf("unrecognized message route: %s", sdk.MsgTypeURL(lsmRedeemMsg))
	}

	// [3] redeem LSM shares from proxyAcc, to obtain a delegation
	msgResp, err = handler(ctx, lsmRedeemMsg)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to execute redeem tokens for shares message",
			types.ErrorKeyVal,
			err,
			types.MsgKeyVal,
			lsmRedeemMsg.String(),
		)

		return math.ZeroInt(), types.ErrLSMRedeemFailed.Wrapf("error: %s; message: %v", err.Error(), lsmRedeemMsg)
	}
	ctx.EventManager().EmitEvents(msgResp.GetEvents())

	if len(msgResp.MsgResponses) != 1 {
		return math.ZeroInt(), errorsmod.Wrapf(
			types.ErrInvalidResponse,
			"expected msg response should be exactly 1, got: %v, responses: %v",
			len(msgResp.MsgResponses), msgResp.MsgResponses,
		)
	}

	var lsmRedeemResp liquidtypes.MsgRedeemTokensForSharesResponse
	if err = k.cdc.Unmarshal(msgResp.MsgResponses[0].Value, &lsmRedeemResp); err != nil {
		return math.ZeroInt(), errorsmod.Wrapf(
			sdkerrors.ErrJSONUnmarshal,
			"cannot unmarshal redeem tokens for shares tx response message: %v",
			err,
		)
	}

	// [4] unstake from user's account.
	unstakeMsg := &stakingtypes.MsgUndelegate{
		DelegatorAddress: userAddress.String(),
		ValidatorAddress: validatorAddress.String(),
		Amount:           lsmRedeemResp.Amount,
	}

	handler = k.router.Handler(unstakeMsg)
	if handler == nil {
		k.Logger(ctx).Error("failed to find undelegate handler")

		return math.ZeroInt(), sdkerrors.ErrUnknownRequest.Wrapf("unrecognized message route: %s", sdk.MsgTypeURL(unstakeMsg))
	}

	msgResp, err = handler(ctx, unstakeMsg)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to execute undelegate message",
			types.ErrorKeyVal,
			err,
			types.MsgKeyVal,
			unstakeMsg.String(),
		)

		return math.ZeroInt(), types.ErrUnstakeFailed.Wrapf("error: %s; message: %v", err.Error(), unstakeMsg)
	}
	ctx.EventManager().EmitEvents(msgResp.GetEvents())

	if len(msgResp.MsgResponses) != 1 {
		return math.ZeroInt(), errorsmod.Wrapf(
			types.ErrInvalidResponse,
			"expected msg response should be exactly 1, got: %v, responses: %v",
			len(msgResp.MsgResponses), msgResp.MsgResponses,
		)
	}

	var msgUndelegateResp stakingtypes.MsgUndelegateResponse
	if err = k.cdc.Unmarshal(msgResp.MsgResponses[0].Value, &msgUndelegateResp); err != nil {
		return math.ZeroInt(), errorsmod.Wrapf(
			sdkerrors.ErrJSONUnmarshal,
			"cannot unmarshal msg undelegate tx response message: %v",
			err,
		)
	}

	return lsmRedeemResp.Amount.Amount, nil
}

// RedelegateWithCap is a wrapper to invoke stakingKeeper.Redelegate but account for
// the amount of liquid staked shares and check against liquid staking cap.
func (k Keeper) RedelegateWithCap(
	ctx sdk.Context,
	delegatorAddress sdk.AccAddress,
	validatorSrc sdk.ValAddress,
	validatorDst sdk.ValAddress,
	bondAmt math.Int,
) (time.Time, error) {
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return time.Time{}, err
	}
	msgRedelegate := &stakingtypes.MsgBeginRedelegate{
		DelegatorAddress:    delegatorAddress.String(),
		ValidatorSrcAddress: validatorSrc.String(),
		ValidatorDstAddress: validatorDst.String(),
		Amount:              sdk.NewCoin(bondDenom, bondAmt),
	}
	handler := k.router.Handler(msgRedelegate)
	if handler == nil {
		k.Logger(ctx).Error("failed to find redelegate handler")

		return time.Time{}, errorsmod.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized message route: %s", sdk.MsgTypeURL(msgRedelegate))
	}
	res, err := handler(ctx, msgRedelegate)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to execute redelegate msg",
			types.ErrorKeyVal,
			err,
			types.MsgKeyVal,
			msgRedelegate.String(),
		)

		return time.Time{}, errorsmod.Wrapf(types.ErrRedelegateFailed, "failed to send redelegate msg with err: %v", err)
	}
	ctx.EventManager().EmitEvents(res.GetEvents())

	if len(res.MsgResponses) != 1 {
		return time.Time{}, errorsmod.Wrapf(
			types.ErrInvalidResponse,
			"expected msg response should be exactly 1, got: %v, responses: %v",
			len(res.MsgResponses), res.MsgResponses,
		)
	}

	var msgRedelegateResponse stakingtypes.MsgBeginRedelegateResponse
	if err = k.cdc.Unmarshal(res.MsgResponses[0].Value, &msgRedelegateResponse); err != nil {
		return time.Time{}, errorsmod.Wrapf(
			sdkerrors.ErrJSONUnmarshal,
			"cannot unmarshal redelegate tx response message: %v",
			err,
		)
	}

	return msgRedelegateResponse.CompletionTime, nil
}

// LSMDelegate captures a staked amount from existing delegation using LSM, re-stakes from proxyAcc and
// mints stkXPRT worth of stk coin value according to NetAmount and performs LiquidDelegate.
func (k Keeper) LSMDelegate(
	ctx sdk.Context,
	delegator sdk.AccAddress,
	validator sdk.ValAddress,
	proxyAcc sdk.AccAddress,
	preLsmStake sdk.Coin,
) (stkXPRTMintAmount math.Int, err error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}
	if params.ModulePaused {
		return math.ZeroInt(), types.ErrModulePaused
	} else if params.LsmDisabled {
		return math.ZeroInt(), types.ErrDisabledLSM
	}

	// check minimum liquid stake amount
	if preLsmStake.Amount.LT(params.MinLiquidStakeAmount) {
		return math.ZeroInt(), types.ErrLessThanMinLiquidStakeAmount
	}

	// check bond denomination
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}
	if preLsmStake.Denom != bondDenom {
		return math.ZeroInt(), errorsmod.Wrapf(
			types.ErrInvalidBondDenom, "invalid coin denomination: got %s, expected %s", preLsmStake.Denom, bondDenom,
		)
	}

	whitelistedValsMap := types.GetWhitelistedValsMap(params.WhitelistedValidators)
	activeVals, err := k.GetActiveLiquidValidators(ctx, whitelistedValsMap)
	if err != nil {
		return math.ZeroInt(), err
	}

	if activeVals.Len() == 0 {
		return math.ZeroInt(), types.ErrActiveLiquidValidatorsNotExists
	}

	totalActiveWeight := activeVals.TotalWeight(whitelistedValsMap)
	activeWeightQuorum := math.LegacyNewDecFromInt(totalActiveWeight).Quo(
		math.LegacyNewDecFromInt(types.TotalValidatorWeight),
	)
	if activeWeightQuorum.LT(types.ActiveLiquidValidatorsWeightQuorum) {
		k.Logger(ctx).Error(
			"active liquid validators weight quorum not reached",
			types.ActiveWeightQuorumKeyVal,
			activeWeightQuorum.String(),
			types.MinActiveWeightQuorumKeyVal,
			types.ActiveLiquidValidatorsWeightQuorum.String(),
		)

		return math.ZeroInt(), errorsmod.Wrapf(
			types.ErrActiveLiquidValidatorsWeightQuorumNotReached, "%s < %s",
			activeWeightQuorum.String(), types.ActiveLiquidValidatorsWeightQuorum.String(),
		)
	}

	if !whitelistedValsMap.IsListed(validator.String()) {
		return math.ZeroInt(), types.ErrLiquidValidatorsNotExists.Wrap("delegation from a non allowed validator")
	}

	// NetAmount must be calculated before send
	nas, err := k.GetNetAmountState(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}

	// perform an LSM tokenize->bank send->redeem flow: moving delegation from user's account onto proxyAcc

	lsmTokenizeMsg := &liquidtypes.MsgTokenizeShares{
		DelegatorAddress:    delegator.String(),
		ValidatorAddress:    validator.String(),
		Amount:              preLsmStake,
		TokenizedShareOwner: proxyAcc.String(),
	}

	handler := k.router.Handler(lsmTokenizeMsg)
	if handler == nil {
		k.Logger(ctx).Error("failed to find tokenize handler")

		return math.ZeroInt(), sdkerrors.ErrUnknownRequest.Wrapf("unrecognized message route: %s", sdk.MsgTypeURL(lsmTokenizeMsg))
	}

	// [1] tokenize delegation into LSM shares
	msgResp, err := handler(ctx, lsmTokenizeMsg)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to execute tokenize shares message",
			types.ErrorKeyVal,
			err,
			types.MsgKeyVal,
			lsmTokenizeMsg.String(),
		)

		return math.ZeroInt(), types.ErrLSMTokenizeFailed.Wrapf("error: %s; message: %v", err.Error(), lsmTokenizeMsg)
	}
	ctx.EventManager().EmitEvents(msgResp.GetEvents())

	if len(msgResp.MsgResponses) != 1 {
		return math.ZeroInt(), errorsmod.Wrapf(
			types.ErrInvalidResponse,
			"expected msg response should be exactly 1, got: %v, responses: %v",
			len(msgResp.MsgResponses), msgResp.MsgResponses,
		)
	}

	var lsmTokenizeResp liquidtypes.MsgTokenizeSharesResponse
	if err = k.cdc.Unmarshal(msgResp.MsgResponses[0].Value, &lsmTokenizeResp); err != nil {
		return math.ZeroInt(), errorsmod.Wrapf(
			sdkerrors.ErrJSONUnmarshal,
			"cannot unmarshal tokenize share tx response message: %v",
			err,
		)
	}

	// [2] send LSM shares to proxyAcc
	err = k.bankKeeper.SendCoins(ctx, delegator, proxyAcc, sdk.NewCoins(lsmTokenizeResp.Amount))
	if err != nil {
		return stkXPRTMintAmount, err
	}

	lsmRedeemMsg := &liquidtypes.MsgRedeemTokensForShares{
		DelegatorAddress: proxyAcc.String(),
		Amount:           lsmTokenizeResp.Amount,
	}

	handler = k.router.Handler(lsmRedeemMsg)
	if handler == nil {
		k.Logger(ctx).Error("failed to find redeem handler")

		return math.ZeroInt(), sdkerrors.ErrUnknownRequest.Wrapf("unrecognized message route: %s", sdk.MsgTypeURL(lsmRedeemMsg))
	}

	// [3] redeem LSM shares from proxyAcc, to obtain a delegation
	msgResp, err = handler(ctx, lsmRedeemMsg)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to execute redeem tokens for shares message",
			types.ErrorKeyVal,
			err,
			types.MsgKeyVal,
			lsmRedeemMsg.String(),
		)

		return math.ZeroInt(), types.ErrLSMRedeemFailed.Wrapf("error: %s; message: %v", err.Error(), lsmRedeemMsg)
	}
	ctx.EventManager().EmitEvents(msgResp.GetEvents())

	if len(msgResp.MsgResponses) != 1 {
		return math.ZeroInt(), errorsmod.Wrapf(
			types.ErrInvalidResponse,
			"expected msg response should be exactly 1, got: %v, responses: %v",
			len(msgResp.MsgResponses), msgResp.MsgResponses,
		)
	}

	var lsmRedeemResp liquidtypes.MsgRedeemTokensForSharesResponse
	if err = k.cdc.Unmarshal(msgResp.MsgResponses[0].Value, &lsmRedeemResp); err != nil {
		return math.ZeroInt(), errorsmod.Wrapf(
			sdkerrors.ErrJSONUnmarshal,
			"cannot unmarshal redeem tokens for shares tx response message: %v",
			err,
		)
	}

	// mint stkxprt, MintAmount = TotalSupply * StakeAmount/NetAmount
	liquidBondDenom, err := k.LiquidBondDenom(ctx)
	if err != nil {
		return math.ZeroInt(), err
	}
	stkXPRTMintAmount = lsmRedeemResp.Amount.Amount

	if nas.StkxprtTotalSupply.IsPositive() {
		stkXPRTMintAmount = types.NativeTokenToStkXPRT(
			stkXPRTMintAmount,
			nas.StkxprtTotalSupply,
			nas.NetAmount,
		)
	}

	if !stkXPRTMintAmount.IsPositive() {
		return math.ZeroInt(), types.ErrTooSmallLiquidStakeAmount
	}

	// mint stkXPRT on module acc
	mintCoin := sdk.NewCoins(sdk.NewCoin(liquidBondDenom, stkXPRTMintAmount))
	err = k.bankKeeper.MintCoins(ctx, types.ModuleName, mintCoin)
	if err != nil {
		return stkXPRTMintAmount, err
	}

	// send stkXPRT to delegator acc
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, delegator, mintCoin)
	if err != nil {
		return stkXPRTMintAmount, err
	}

	// but immediately lock new stkXPRT into LP on behalf of the delegator
	_, err = k.LockOnLP(ctx, delegator, sdk.NewCoin(liquidBondDenom, stkXPRTMintAmount))
	if err != nil {
		return stkXPRTMintAmount, err
	}

	return stkXPRTMintAmount, err
}

// LiquidDelegate delegates staking amount to active validators by proxy account.
func (k Keeper) LiquidDelegate(ctx sdk.Context, proxyAcc sdk.AccAddress, activeVals types.ActiveLiquidValidators, stakingAmt math.Int, whitelistedValsMap types.WhitelistedValsMap) (err error) {
	// crumb may occur due to a decimal point error in dividing the staking amount into the weight of liquid validators, It added on first active liquid validator
	weightedAmt, crumb := types.DivideByWeight(activeVals, stakingAmt, whitelistedValsMap)
	if len(weightedAmt) == 0 {
		k.Logger(ctx).Error(
			"invalid active liquid validators",
			types.ActiveValidatorsKeyVal,
			activeVals,
			types.AmountKeyVal,
			stakingAmt.String(),
			types.WhitelistedValidatorsMapKeyVal,
			whitelistedValsMap,
		)

		return types.ErrInvalidActiveLiquidValidators
	}
	weightedAmt[0] = weightedAmt[0].Add(crumb)
	for i, val := range activeVals {
		if !weightedAmt[i].IsPositive() {
			continue
		}
		validator, _ := k.stakingKeeper.GetValidator(ctx, val.GetOperator())
		err = k.DelegateWithCap(ctx, proxyAcc, validator, weightedAmt[i])
		if err != nil {
			return errorsmod.Wrapf(err, "failed to delegate to validator %s", val.GetOperator())
		}
	}
	return nil
}

// LiquidUnstake burns unstakingStkXPRT and performs LiquidUnbond to active liquid validators with del shares worth of shares according to NetAmount with each validators current weight.
func (k Keeper) LiquidUnstake(
	ctx sdk.Context, proxyAcc, liquidStaker sdk.AccAddress, unstakingStkXPRT sdk.Coin,
) (time.Time, math.Int, []stakingtypes.UnbondingDelegation, math.Int, error) {
	params, err := k.GetParams(ctx)
	if err != nil {
		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), err
	}
	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), err
	}

	if params.ModulePaused {
		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), types.ErrModulePaused
	}

	// check bond denomination
	liquidBondDenom, err := k.LiquidBondDenom(ctx)
	if err != nil {
		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), err
	}
	if unstakingStkXPRT.Denom != liquidBondDenom {
		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), errorsmod.Wrapf(
			types.ErrInvalidLiquidBondDenom, "invalid coin denomination: got %s, expected %s", unstakingStkXPRT.Denom, liquidBondDenom,
		)
	}

	// Get NetAmount states
	nas, err := k.GetNetAmountState(ctx)
	if err != nil {
		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), err
	}

	if unstakingStkXPRT.Amount.GT(nas.StkxprtTotalSupply) || nas.StkxprtTotalSupply.IsZero() {
		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), types.ErrInvalidStkXPRTSupply
	}

	// UnstakeAmount = NetAmount * StkXPRTAmount/TotalSupply * (1-UnstakeFeeRate)
	unbondingAmount := types.StkXPRTToNativeToken(unstakingStkXPRT.Amount, nas.StkxprtTotalSupply, nas.NetAmount)
	unbondingAmount = types.DeductFeeRate(unbondingAmount, params.UnstakeFeeRate)
	unbondingAmountInt := unbondingAmount.TruncateInt()

	if !unbondingAmountInt.IsPositive() {
		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), types.ErrTooSmallLiquidUnstakingAmount
	}

	// burn stkxprt
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, liquidStaker, types.ModuleName, sdk.NewCoins(unstakingStkXPRT))
	if err != nil {
		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), err
	}
	err = k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(liquidBondDenom, unstakingStkXPRT.Amount)))
	if err != nil {
		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), err
	}

	liquidVals := k.GetAllLiquidValidators(ctx)
	totalLiquidTokens, liquidTokenMap := liquidVals.TotalLiquidTokens(ctx, k.stakingKeeper, false)

	// if no totalLiquidTokens, withdraw directly from balance of proxy acc
	if !totalLiquidTokens.IsPositive() {
		if nas.ProxyAccBalance.GTE(unbondingAmountInt) {
			err = k.bankKeeper.SendCoins(
				ctx,
				types.LiquidStakeProxyAcc,
				liquidStaker,
				sdk.NewCoins(sdk.NewCoin(
					bondDenom,
					unbondingAmountInt,
				)),
			)
			if err != nil {
				return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), err
			}

			return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, unbondingAmountInt, nil
		}

		k.Logger(ctx).Error(
			"non-positive total liquid tokens",
			types.ValidatorsKeyVal,
			liquidVals,
			types.TotalLiquidTokensKeyVal,
			totalLiquidTokens.String(),
		)

		// error case where there is a quantity that are unbonding balance or remaining rewards that is not re-stake or withdrawn in netAmount.
		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), types.ErrInsufficientProxyAccBalance
	}

	// fail when no liquid validators to unbond
	if liquidVals.Len() == 0 {
		k.Logger(ctx).Error(
			"no liquid validators to unbond",
			types.ValidatorsKeyVal,
			liquidVals,
		)

		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), types.ErrLiquidValidatorsNotExists
	}

	// prioritize inactive liquid validators in the list to be used in DivideByCurrentWeight
	liquidVals = k.PrioritiseInactiveLiquidValidators(ctx, liquidVals)

	// crumb may occur due to a decimal error in dividing the unstaking stkXPRT into the weight of liquid validators, it will remain in the NetAmount
	unbondingAmounts, crumb := types.DivideByCurrentWeight(liquidVals, unbondingAmount, totalLiquidTokens, liquidTokenMap)
	if !unbondingAmount.Sub(crumb).IsPositive() {
		return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), types.ErrTooSmallLiquidUnstakingAmount
	}

	totalReturnAmount := math.ZeroInt()

	var ubdTime time.Time
	ubds := make([]stakingtypes.UnbondingDelegation, 0, len(liquidVals))
	for i, val := range liquidVals {
		// skip zero weight liquid validator
		if !unbondingAmounts[i].IsPositive() {
			continue
		}

		var ubd stakingtypes.UnbondingDelegation
		var returnAmount math.Int
		var weightedShare math.LegacyDec

		// calculate delShares from tokens with validation
		weightedShare, err = k.stakingKeeper.ValidateUnbondAmount(ctx, proxyAcc, val.GetOperator(), unbondingAmounts[i].TruncateInt())
		if err != nil {
			k.Logger(ctx).Error(
				"failed to validate unbond amount",
				types.ErrorKeyVal,
				err,
				types.ValidatorKeyVal,
				val.GetOperator().String(),
				types.AmountKeyVal,
				unbondingAmounts[i].TruncateInt().String(),
			)

			return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), err
		}

		if !weightedShare.IsPositive() {
			continue
		}

		// unbond with weightedShare
		ubdTime, returnAmount, ubd, err = k.LiquidUnbond(ctx, proxyAcc, liquidStaker, val.GetOperator(), weightedShare, true, sdk.NewCoin(bondDenom, unbondingAmounts[i].TruncateInt()))
		if err != nil {
			return time.Time{}, math.ZeroInt(), []stakingtypes.UnbondingDelegation{}, math.ZeroInt(), err
		}

		ubds = append(ubds, ubd)
		totalReturnAmount = totalReturnAmount.Add(returnAmount)
	}

	return ubdTime, totalReturnAmount, ubds, math.ZeroInt(), nil
}

// LiquidUnbond unbond delegation shares to active validators by proxy account.
func (k Keeper) LiquidUnbond(
	ctx sdk.Context, proxyAcc, liquidStaker sdk.AccAddress, valAddr sdk.ValAddress, shares math.LegacyDec, checkMaxEntries bool, unbondAmount sdk.Coin,
) (time.Time, math.Int, stakingtypes.UnbondingDelegation, error) {
	_, err := k.stakingKeeper.GetValidator(ctx, valAddr)
	if err != nil {
		return time.Time{}, math.ZeroInt(), stakingtypes.UnbondingDelegation{}, stakingtypes.ErrNoDelegatorForAddress
	}

	// If checkMaxEntries is true, perform a maximum limit unbonding entries check.
	hasMaxUbdEntries, err := k.stakingKeeper.HasMaxUnbondingDelegationEntries(ctx, liquidStaker, valAddr)
	if err != nil {
		return time.Time{}, math.ZeroInt(), stakingtypes.UnbondingDelegation{}, err
	}
	if checkMaxEntries && hasMaxUbdEntries {
		return time.Time{}, math.ZeroInt(), stakingtypes.UnbondingDelegation{}, stakingtypes.ErrMaxUnbondingDelegationEntries
	}

	// unbond from proxy account
	returnAmount, err := k.UnbondWithCap(ctx, proxyAcc, valAddr, unbondAmount, liquidStaker)
	if err != nil {
		return time.Time{}, math.ZeroInt(), stakingtypes.UnbondingDelegation{}, err
	}

	//// Unbonding from proxy account, but queues to liquid staker.
	ubdtime, err := k.stakingKeeper.UnbondingTime(ctx)
	if err != nil {
		return time.Time{}, math.ZeroInt(), stakingtypes.UnbondingDelegation{}, err
	}
	completionTime := ctx.BlockHeader().Time.Add(ubdtime)
	ubd, err := k.stakingKeeper.GetUnbondingDelegation(ctx, liquidStaker, valAddr)
	if err != nil {
		k.Logger(ctx).Error(
			"failed to find unbonding delegation",
			types.DelegatorKeyVal,
			liquidStaker.String(),
			types.ValidatorKeyVal,
			valAddr.String(),
		)

		return time.Time{}, math.ZeroInt(), stakingtypes.UnbondingDelegation{}, types.ErrInvalidResponse.Wrap("expected undelegation entry, found none")
	}

	return completionTime, returnAmount, ubd, nil
}

// PrioritiseInactiveLiquidValidators sorts LiquidValidators array to have inactive validators first. Used for the case when
// unbonding should begin from the inactive validators first.
func (k Keeper) PrioritiseInactiveLiquidValidators(
	ctx sdk.Context,
	vs types.LiquidValidators,
) types.LiquidValidators {
	sort.SliceStable(vs, func(i, j int) bool {
		vs1, vs1err := k.stakingKeeper.GetValidator(ctx, vs[i].GetOperator())
		vs2, vs2err := k.stakingKeeper.GetValidator(ctx, vs[j].GetOperator())

		if vs1err != nil && vs2err == nil {
			// only one case when less
			return true
		} else if vs1err == nil && vs2err == nil {
			// both exist, compare status

			vs1Active := vs[i].GetStatus(types.ActiveCondition(
				vs1,
				true,
				k.IsTombstoned(ctx, vs1),
			))
			vs2Active := vs[j].GetStatus(types.ActiveCondition(
				vs2,
				true,
				k.IsTombstoned(ctx, vs2),
			))

			if vs1Active != types.ValidatorStatusActive &&
				vs2Active == types.ValidatorStatusActive {
				// only one case when is less
				return true
			}

			// not less, or are equal
			return false
		}

		// not less, or are equal
		return false
	})

	return vs
}

// CheckDelegationStates returns total remaining rewards, delshares, liquid tokens of delegations by proxy account
func (k Keeper) CheckDelegationStates(ctx sdk.Context, proxyAcc sdk.AccAddress) (math.LegacyDec, math.LegacyDec, math.Int, error) {
	totalRewards := math.LegacyZeroDec()
	totalDelShares := math.LegacyZeroDec()
	totalLiquidTokens := math.ZeroInt()

	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return totalRewards, totalDelShares, totalLiquidTokens, err
	}
	// Cache ctx for calculate rewards
	cachedCtx, _ := ctx.CacheContext()
	k.stakingKeeper.IterateDelegations(
		cachedCtx, proxyAcc,
		func(_ int64, del stakingtypes.DelegationI) (stop bool) {
			valAddr := del.GetValidatorAddr()
			valAddrBz, err := sdk.ValAddressFromBech32(valAddr)
			if err != nil {
				return false
			}
			val, err := k.stakingKeeper.Validator(cachedCtx, valAddrBz)
			if err != nil {
				return false
			}
			endingPeriod, err := k.distrKeeper.IncrementValidatorPeriod(cachedCtx, val)
			if err != nil {
				return false
			}
			delReward, err := k.distrKeeper.CalculateDelegationRewards(cachedCtx, val, del, endingPeriod)
			if err != nil {
				return false
			}
			delShares := del.GetShares()
			if delShares.IsPositive() {
				totalDelShares = totalDelShares.Add(delShares)
				liquidTokens := val.TokensFromSharesTruncated(delShares).TruncateInt()
				totalLiquidTokens = totalLiquidTokens.Add(liquidTokens)
				totalRewards = totalRewards.Add(delReward.AmountOf(bondDenom).TruncateDec())
			}
			return false
		},
	)

	return totalRewards, totalDelShares, totalLiquidTokens, nil
}

func (k Keeper) WithdrawLiquidRewards(ctx sdk.Context, proxyAcc sdk.AccAddress) {
	// iterate over all the delegations (even those out of the active set) and withdraw rewards
	k.stakingKeeper.IterateDelegations(
		ctx, proxyAcc,
		func(_ int64, del stakingtypes.DelegationI) (stop bool) {
			// construct the withdrawal rewards message
			msgWithdraw := &distributiontypes.MsgWithdrawDelegatorReward{
				DelegatorAddress: proxyAcc.String(),
				ValidatorAddress: del.GetValidatorAddr(),
			}

			// run the message handler
			handler := k.router.Handler(msgWithdraw)
			if handler == nil {
				k.Logger(ctx).Error("could not find distribution handler for withdraw rewards msg")
				return true
			}
			res, err := handler(ctx, msgWithdraw)
			if err != nil {
				k.Logger(ctx).Error(
					"failed to execute withdraw rewards msg",
					types.MsgKeyVal,
					msgWithdraw.String(),
					types.ErrorKeyVal,
					err,
				)
				// no need to return here, will be picked up in the next epoch
			} else {
				// emit the events
				ctx.EventManager().EmitEvents(res.GetEvents())
			}

			return false
		},
	)
}

// GetLiquidValidator get a single liquid validator
func (k Keeper) GetLiquidValidator(ctx sdk.Context, addr sdk.ValAddress) (val types.LiquidValidator, found bool) {
	store := k.storeService.OpenKVStore(ctx)

	value, err := store.Get(types.GetLiquidValidatorKey(addr))
	if err != nil {
		return val, false
	}
	if value == nil {
		return val, false
	}

	val = types.MustUnmarshalLiquidValidator(k.cdc, value)
	return val, true
}

// SetLiquidValidator set the main record holding liquid validator details
func (k Keeper) SetLiquidValidator(ctx sdk.Context, val types.LiquidValidator) {
	store := k.storeService.OpenKVStore(ctx)

	bz := types.MustMarshalLiquidValidator(k.cdc, &val)
	store.Set(types.GetLiquidValidatorKey(val.GetOperator()), bz)
}

// RemoveLiquidValidator remove a liquid validator on kv store
func (k Keeper) RemoveLiquidValidator(ctx sdk.Context, val types.LiquidValidator) {
	store := k.storeService.OpenKVStore(ctx)

	store.Delete(types.GetLiquidValidatorKey(val.GetOperator()))
}

// GetAllLiquidValidators gets the set of all liquid validators, with no pagination limits.
func (k Keeper) GetAllLiquidValidators(ctx sdk.Context) (vals types.LiquidValidators) {
	store := k.storeService.OpenKVStore(ctx)

	vals = types.LiquidValidators{}
	iterator, err := store.Iterator(types.LiquidValidatorsKey, storetypes.PrefixEndBytes(types.LiquidValidatorsKey))
	if err != nil {
		panic(err)
	}

	defer iterator.Close() //nolint:errcheck

	for ; iterator.Valid(); iterator.Next() {
		val := types.MustUnmarshalLiquidValidator(k.cdc, iterator.Value())
		vals = append(vals, val)
	}

	return vals
}

// GetActiveLiquidValidators get the set of active liquid validators.
func (k Keeper) GetActiveLiquidValidators(ctx sdk.Context, whitelistedValsMap types.WhitelistedValsMap) (vals types.ActiveLiquidValidators, err error) {
	store := k.storeService.OpenKVStore(ctx)

	iterator, err := store.Iterator(types.LiquidValidatorsKey, storetypes.PrefixEndBytes(types.LiquidValidatorsKey))
	if err != nil {
		return vals, err
	}
	defer iterator.Close() //nolint:errcheck

	for ; iterator.Valid(); iterator.Next() {
		val := types.MustUnmarshalLiquidValidator(k.cdc, iterator.Value())
		if k.IsActiveLiquidValidator(ctx, val, whitelistedValsMap) {
			vals = append(vals, val)
		}
	}
	return vals, nil
}

func (k Keeper) GetAllLiquidValidatorStates(ctx sdk.Context) (liquidValidatorStates []types.LiquidValidatorState, err error) {
	lvs := k.GetAllLiquidValidators(ctx)
	params, err := k.GetParams(ctx)
	if err != nil {
		return nil, err
	}

	whitelistedValsMap := params.WhitelistedValsMap()
	for _, lv := range lvs {
		active := k.IsActiveLiquidValidator(ctx, lv, whitelistedValsMap)
		lvState := types.LiquidValidatorState{
			OperatorAddress: lv.OperatorAddress,
			Weight:          lv.GetWeight(whitelistedValsMap, active),
			Status:          lv.GetStatus(active),
			DelShares:       lv.GetDelShares(ctx, k.stakingKeeper),
			LiquidTokens:    lv.GetLiquidTokens(ctx, k.stakingKeeper, false),
		}
		liquidValidatorStates = append(liquidValidatorStates, lvState)
	}
	return liquidValidatorStates, nil
}

func (k Keeper) GetLiquidValidatorState(ctx sdk.Context, addr sdk.ValAddress) (liquidValidatorState types.LiquidValidatorState, found bool) {
	lv, found := k.GetLiquidValidator(ctx, addr)
	if !found {
		return types.LiquidValidatorState{
			OperatorAddress: addr.String(),
			Weight:          math.ZeroInt(),
			Status:          types.ValidatorStatusUnspecified,
			DelShares:       math.LegacyZeroDec(),
			LiquidTokens:    math.ZeroInt(),
		}, false
	}
	params, err := k.GetParams(ctx)
	if err != nil {
		panic(err)
	}
	whitelistedValsMap := params.WhitelistedValsMap()
	active := k.IsActiveLiquidValidator(ctx, lv, whitelistedValsMap)
	return types.LiquidValidatorState{
		OperatorAddress: lv.OperatorAddress,
		Weight:          lv.GetWeight(whitelistedValsMap, active),
		Status:          lv.GetStatus(active),
		DelShares:       lv.GetDelShares(ctx, k.stakingKeeper),
		LiquidTokens:    lv.GetLiquidTokens(ctx, k.stakingKeeper, false),
	}, true
}

func (k Keeper) IsActiveLiquidValidator(ctx sdk.Context, lv types.LiquidValidator, whitelistedValsMap types.WhitelistedValsMap) bool {
	val, err := k.stakingKeeper.GetValidator(ctx, lv.GetOperator())
	if err != nil {
		return false
	}
	return types.ActiveCondition(val, whitelistedValsMap.IsListed(lv.OperatorAddress), k.IsTombstoned(ctx, val))
}

func (k Keeper) IsTombstoned(ctx sdk.Context, val stakingtypes.Validator) bool {
	consPk, err := val.ConsPubKey()
	if err != nil {
		return false
	}
	return k.slashingKeeper.IsTombstoned(ctx, sdk.ConsAddress(consPk.Address()))
}

func (k Keeper) GetWeightMap(ctx sdk.Context, liquidVals types.LiquidValidators, whitelistedValsMap types.WhitelistedValsMap) (map[string]math.Int, math.Int) {
	weightMap := map[string]math.Int{}
	totalWeight := math.ZeroInt()
	for _, val := range liquidVals {
		weight := val.GetWeight(whitelistedValsMap, k.IsActiveLiquidValidator(ctx, val, whitelistedValsMap))
		totalWeight = totalWeight.Add(weight)
		weightMap[val.OperatorAddress] = weight
	}
	return weightMap, totalWeight
}
