package app

import sdk "github.com/cosmos/cosmos-sdk/types"

// OracleKeeper for feeder validation
type OracleKeeper interface {
	ValidateFeeder(ctx sdk.Context, feederAddr sdk.Address, validatorAddr sdk.ValAddress) error
}
