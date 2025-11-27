package types_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/types"
)

func TestMsgLiquidStake(t *testing.T) {
	delegatorAddr := sdk.AccAddress(crypto.AddressHash([]byte("delegatorAddr")))
	stakingCoin := sdk.NewCoin("uxprt", math.NewInt(1))

	testCases := []struct {
		expectedErr string
		msg         *types.MsgLiquidStake
	}{
		{
			"", // empty means no error expected
			types.NewMsgLiquidStake(delegatorAddr, stakingCoin),
		},
		{
			"invalid delegator address \"\": empty address string is not allowed: invalid address",
			types.NewMsgLiquidStake(sdk.AccAddress{}, stakingCoin),
		},
		{
			"staking amount must not be zero: invalid request",
			types.NewMsgLiquidStake(delegatorAddr, sdk.NewCoin("token", math.NewInt(0))),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgLiquidStake{}, tc.msg)
		require.Equal(t, types.MsgTypeLiquidStake, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.GetDelegator(), signers[0])
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgLiquidUnstake(t *testing.T) {
	delegatorAddr := sdk.AccAddress(crypto.AddressHash([]byte("delegatorAddr")))
	stakingCoin := sdk.NewCoin("stk/uxprt", math.NewInt(1))

	testCases := []struct {
		expectedErr string
		msg         *types.MsgLiquidUnstake
	}{
		{
			"", // empty means no error expected
			types.NewMsgLiquidUnstake(delegatorAddr, stakingCoin),
		},
		{
			"invalid delegator address \"\": empty address string is not allowed: invalid address",
			types.NewMsgLiquidUnstake(sdk.AccAddress{}, stakingCoin),
		},
		{
			"unstaking amount must not be zero: invalid request",
			types.NewMsgLiquidUnstake(delegatorAddr, sdk.NewCoin("btoken", math.NewInt(0))),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgLiquidUnstake{}, tc.msg)
		require.Equal(t, types.MsgTypeLiquidUnstake, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.GetDelegator(), signers[0])
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgStakeToLP(t *testing.T) {
	delegatorAddr := sdk.AccAddress(crypto.AddressHash([]byte("delegatorAddr")))
	validatorAddr := sdk.ValAddress(crypto.AddressHash([]byte("validatorAddr")))
	stakingCoin := sdk.NewCoin("uxprt", math.NewInt(1))
	zeroStakingCoin := sdk.NewCoin("uxprt", math.NewInt(0))

	testCases := []struct {
		expectedErr string
		msg         *types.MsgStakeToLP
	}{
		{
			"", // empty means no error expected
			types.NewMsgStakeToLP(delegatorAddr, validatorAddr, stakingCoin, zeroStakingCoin),
		},
		{
			"",
			types.NewMsgStakeToLP(delegatorAddr, validatorAddr, stakingCoin, sdk.Coin{}),
		},
		{
			"",
			types.NewMsgStakeToLP(delegatorAddr, validatorAddr, stakingCoin, stakingCoin),
		},
		{
			"invalid delegator address \"\": empty address string is not allowed: invalid address",
			types.NewMsgStakeToLP(sdk.AccAddress{}, validatorAddr, stakingCoin, zeroStakingCoin),
		},
		{
			"invalid validator address \"\": empty address string is not allowed: invalid address",
			types.NewMsgStakeToLP(delegatorAddr, sdk.ValAddress{}, stakingCoin, zeroStakingCoin),
		},
		{
			"staked amount must not be zero: invalid request",
			types.NewMsgStakeToLP(delegatorAddr, validatorAddr, zeroStakingCoin, zeroStakingCoin),
		},
		{
			"staked amount must not be zero: invalid request",
			types.NewMsgStakeToLP(delegatorAddr, validatorAddr, zeroStakingCoin, stakingCoin),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgStakeToLP{}, tc.msg)
		require.Equal(t, types.MsgTypeStakeToLP, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.GetDelegator(), signers[0])
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgSetModulePaused(t *testing.T) {
	authorityAddr := sdk.AccAddress(crypto.AddressHash([]byte("authorityAddr")))

	testCases := []struct {
		expectedErr string
		msg         *types.MsgSetModulePaused
	}{
		{
			"", // empty means no error expected
			types.NewMsgSetModulePaused(authorityAddr, true),
		},
		{
			"", // empty means no error expected
			types.NewMsgSetModulePaused(authorityAddr, false),
		},
		{
			"invalid authority address \"\": empty address string is not allowed: invalid address",
			types.NewMsgSetModulePaused(sdk.AccAddress{}, true),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgSetModulePaused{}, tc.msg)
		require.Equal(t, types.MsgTypeSetModulePaused, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.Authority, signers[0].String())
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}

func TestMsgUpdateParams(t *testing.T) {
	authorityAddr := sdk.AccAddress(crypto.AddressHash([]byte("authorityAddr")))
	validParams := types.DefaultParams()

	testCases := []struct {
		expectedErr string
		msg         *types.MsgUpdateParams
	}{
		{
			"", // empty means no error expected
			types.NewMsgUpdateParams(authorityAddr, validParams),
		},
		{
			"invalid authority address \"\": empty address string is not allowed: invalid address",
			types.NewMsgUpdateParams(sdk.AccAddress{}, validParams),
		},
	}

	for _, tc := range testCases {
		require.IsType(t, &types.MsgUpdateParams{}, tc.msg)
		require.Equal(t, types.MsgTypeUpdateParams, tc.msg.Type())
		require.Equal(t, types.RouterKey, tc.msg.Route())
		require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(tc.msg)), tc.msg.GetSignBytes())

		err := tc.msg.ValidateBasic()
		if tc.expectedErr == "" {
			require.Nil(t, err)
			signers := tc.msg.GetSigners()
			require.Len(t, signers, 1)
			require.Equal(t, tc.msg.Authority, signers[0].String())
		} else {
			require.EqualError(t, err, tc.expectedErr)
		}
	}
}
