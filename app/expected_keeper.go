package app

import sdk "github.com/cosmos/cosmos-sdk/types"

// OracleKeeper for feeder validation
type OracleKeeper interface {
	ValidateFeeder(ctx sdk.Context, validatorAddr sdk.ValAddress, feederAddr sdk.AccAddress) error
}
