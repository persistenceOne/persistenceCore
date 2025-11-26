package types

import "cosmossdk.io/errors"

// Sentinel errors for the liquidstake module.
var (
	ErrActiveLiquidValidatorsNotExists              = errors.Register(ModuleName, 2, "active liquid validators not exists")
	ErrInvalidDenom                                 = errors.Register(ModuleName, 3, "invalid denom")
	ErrInvalidBondDenom                             = errors.Register(ModuleName, 4, "invalid bond denom")
	ErrInvalidLiquidBondDenom                       = errors.Register(ModuleName, 5, "invalid liquid bond denom")
	ErrNotImplementedYet                            = errors.Register(ModuleName, 6, "not implemented yet")
	ErrLessThanMinLiquidStakeAmount                 = errors.Register(ModuleName, 7, "staking amount should be over params.min_liquid_stake_amount")
	ErrInvalidStkXPRTSupply                         = errors.Register(ModuleName, 8, "invalid liquid bond denom supply")
	ErrInvalidActiveLiquidValidators                = errors.Register(ModuleName, 9, "invalid active liquid validators")
	ErrLiquidValidatorsNotExists                    = errors.Register(ModuleName, 10, "liquid validators not exists")
	ErrInsufficientProxyAccBalance                  = errors.Register(ModuleName, 11, "insufficient liquid tokens or balance of proxy account, need to wait for new liquid validator to be added or unbonding of proxy account to be completed")
	ErrTooSmallLiquidStakeAmount                    = errors.Register(ModuleName, 12, "liquid stake amount is too small, the result becomes zero")
	ErrTooSmallLiquidUnstakingAmount                = errors.Register(ModuleName, 13, "liquid unstaking amount is too small, the result becomes zero")
	ErrNoLPContractAddress                          = errors.Register(ModuleName, 14, "CW address of an LP contract is not set")
	ErrDisabledLSM                                  = errors.Register(ModuleName, 15, "LSM delegation is disabled")
	ErrLSMTokenizeFailed                            = errors.Register(ModuleName, 16, "LSM tokenization failed")
	ErrLSMRedeemFailed                              = errors.Register(ModuleName, 17, "LSM redemption failed")
	ErrLPContract                                   = errors.Register(ModuleName, 18, "CW contract execution failed")
	ErrWhitelistedValidatorsList                    = errors.Register(ModuleName, 19, "whitelisted validators list incorrect")
	ErrActiveLiquidValidatorsWeightQuorumNotReached = errors.Register(ModuleName, 20, "active liquid validators weight quorum not reached")
	ErrModulePaused                                 = errors.Register(ModuleName, 21, "module functions have been paused")
	ErrDelegationFailed                             = errors.Register(ModuleName, 22, "delegation failed")
	ErrUnbondFailed                                 = errors.Register(ModuleName, 23, "unbond failed")
	ErrInvalidResponse                              = errors.Register(ModuleName, 24, "invalid response")
	ErrUnstakeFailed                                = errors.Register(ModuleName, 25, "Unstaking failed")
	ErrRedelegateFailed                             = errors.Register(ModuleName, 26, "Redelegate failed")
)
