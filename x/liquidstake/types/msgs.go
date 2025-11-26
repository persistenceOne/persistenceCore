package types

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ sdk.Msg = (*MsgLiquidStake)(nil)
	_ sdk.Msg = (*MsgLiquidUnstake)(nil)
	_ sdk.Msg = (*MsgUpdateParams)(nil)
	_ sdk.Msg = (*MsgStakeToLP)(nil)
	_ sdk.Msg = (*MsgUpdateWhitelistedValidators)(nil)
	_ sdk.Msg = (*MsgSetModulePaused)(nil)
)

// Message types for the liquidstake module
const (
	MsgTypeLiquidStake                 = "liquid_stake"
	MsgTypeLiquidUnstake               = "liquid_unstake"
	MsgTypeStakeToLP                   = "stake_to_lp"
	MsgTypeUpdateParams                = "update_params"
	MsgTypeUpdateWhitelistedValidators = "update_whitelisted_validators"
	MsgTypeSetModulePaused             = "set_module_paused"
)

// NewMsgLiquidStake creates a new MsgLiquidStake.
func NewMsgLiquidStake(
	liquidStaker sdk.AccAddress,
	amount sdk.Coin,
) *MsgLiquidStake {
	return &MsgLiquidStake{
		DelegatorAddress: liquidStaker.String(),
		Amount:           amount,
	}
}

func (m *MsgLiquidStake) Route() string { return RouterKey }

func (m *MsgLiquidStake) Type() string { return MsgTypeLiquidStake }

func (m *MsgLiquidStake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid delegator address %q: %v", m.DelegatorAddress, err)
	}
	if ok := m.Amount.IsZero(); ok {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "staking amount must not be zero")
	}
	if err := m.Amount.Validate(); err != nil {
		return err
	}
	return nil
}

func (m *MsgLiquidStake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgLiquidStake) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m *MsgLiquidStake) GetDelegator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgStakeToLP creates a new MsgStakeToLP.
func NewMsgStakeToLP(
	liquidStaker sdk.AccAddress,
	validator sdk.ValAddress,
	stakedAmount,
	liquidAmount sdk.Coin,
) *MsgStakeToLP {
	return &MsgStakeToLP{
		DelegatorAddress: liquidStaker.String(),
		ValidatorAddress: validator.String(),
		StakedAmount:     stakedAmount,
		LiquidAmount:     liquidAmount,
	}
}

func (m *MsgStakeToLP) Route() string { return RouterKey }

func (m *MsgStakeToLP) Type() string { return MsgTypeStakeToLP }

func (m *MsgStakeToLP) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid delegator address %q: %v", m.DelegatorAddress, err)
	}

	if _, err := sdk.ValAddressFromBech32(m.ValidatorAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address %q: %v", m.ValidatorAddress, err)
	}

	if (m.StakedAmount == sdk.Coin{}) || m.StakedAmount.IsZero() {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "staked amount must not be zero")
	} else if err := m.StakedAmount.Validate(); err != nil {
		return err
	}

	if (m.LiquidAmount != sdk.Coin{}) && !m.LiquidAmount.IsZero() {
		if err := m.LiquidAmount.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (m *MsgStakeToLP) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgStakeToLP) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m *MsgStakeToLP) GetValidator() sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(m.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func (m *MsgStakeToLP) GetDelegator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgLiquidUnstake creates a new MsgLiquidUnstake.
func NewMsgLiquidUnstake(
	liquidStaker sdk.AccAddress,
	amount sdk.Coin,
) *MsgLiquidUnstake {
	return &MsgLiquidUnstake{
		DelegatorAddress: liquidStaker.String(),
		Amount:           amount,
	}
}

func (m *MsgLiquidUnstake) Route() string { return RouterKey }

func (m *MsgLiquidUnstake) Type() string { return MsgTypeLiquidUnstake }

func (m *MsgLiquidUnstake) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.DelegatorAddress); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid delegator address %q: %v", m.DelegatorAddress, err)
	}
	if ok := m.Amount.IsZero(); ok {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "unstaking amount must not be zero")
	}
	if err := m.Amount.Validate(); err != nil {
		return err
	}
	return nil
}

func (m *MsgLiquidUnstake) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m *MsgLiquidUnstake) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m *MsgLiquidUnstake) GetDelegator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// NewMsgUpdateParams creates a new MsgUpdateParams.
func NewMsgUpdateParams(authority sdk.AccAddress, amount Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: authority.String(),
		Params:    amount,
	}
}

func (m *MsgUpdateParams) Route() string {
	return RouterKey
}

// Type should return the action
func (m *MsgUpdateParams) Type() string {
	return MsgTypeUpdateParams
}

// GetSignBytes encodes the message for signing
func (m *MsgUpdateParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (m *MsgUpdateParams) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address %q: %v", m.Authority, err)
	}

	err := m.Params.Validate()
	if err != nil {
		return err
	}
	return nil
}

// NewMsgUpdateWhitelistedValidators creates a new MsgUpdateWhitelistedValidators.
func NewMsgUpdateWhitelistedValidators(authority sdk.AccAddress, list []WhitelistedValidator) *MsgUpdateWhitelistedValidators {
	return &MsgUpdateWhitelistedValidators{
		Authority:             authority.String(),
		WhitelistedValidators: list,
	}
}

func (m *MsgUpdateWhitelistedValidators) Route() string {
	return RouterKey
}

// Type should return the action
func (m *MsgUpdateWhitelistedValidators) Type() string {
	return MsgTypeUpdateWhitelistedValidators
}

// GetSignBytes encodes the message for signing
func (m *MsgUpdateWhitelistedValidators) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgUpdateWhitelistedValidators) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}

func (m *MsgUpdateWhitelistedValidators) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address %q: %v", m.Authority, err)
	}

	err := validateWhitelistedValidators(m.WhitelistedValidators)
	if err != nil {
		return err
	}

	return nil
}

func (w *WhitelistedValidator) GetValidatorAddress() sdk.ValAddress {
	valAddr, err := sdk.ValAddressFromBech32(w.ValidatorAddress)
	if err != nil {
		panic(err)
	}

	return valAddr
}

// NewMsgSetModulePaused creates a new MsgSetModulePaused.
func NewMsgSetModulePaused(authority sdk.AccAddress, isPaused bool) *MsgSetModulePaused {
	return &MsgSetModulePaused{
		Authority: authority.String(),
		IsPaused:  isPaused,
	}
}

func (m *MsgSetModulePaused) Route() string {
	return RouterKey
}

// Type should return the action
func (m *MsgSetModulePaused) Type() string {
	return MsgTypeSetModulePaused
}

// GetSignBytes encodes the message for signing
func (m *MsgSetModulePaused) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m *MsgSetModulePaused) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(m.Authority)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{addr}
}

func (m *MsgSetModulePaused) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
		return errors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address %q: %v", m.Authority, err)
	}

	return nil
}
