package types

import (
	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewGenesisState returns new GenesisState instance.
func NewGenesisState(params Params, liquidValidators []LiquidValidator) *GenesisState {
	return &GenesisState{
		Params:           params,
		LiquidValidators: liquidValidators,
	}
}

// DefaultGenesisState returns the default genesis state.
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(
		DefaultParams(),
		[]LiquidValidator{},
	)
}

// ValidateGenesis validates GenesisState.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}
	for _, lv := range data.LiquidValidators {
		if err := lv.Validate(); err != nil {
			return errors.Wrapf(
				sdkerrors.ErrInvalidAddress,
				"invalid liquid validator %s: %v", lv, err)
		}
	}
	return nil
}
