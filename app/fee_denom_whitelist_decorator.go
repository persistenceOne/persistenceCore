package app

import (
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type FeeDenomWhitelistDecorator struct {
	whitelistMap map[string]bool
	whitelistStr string // this is used for err msg only
}

func NewFeeDenomWhitelistDecorator(denomsWhitelist []string) *FeeDenomWhitelistDecorator {
	whitelistMap := map[string]bool{}
	for _, denom := range denomsWhitelist {
		// must be valid denom
		if err := sdk.ValidateDenom(denom); err != nil {
			panic(fmt.Sprintf("invalid denoms whiltelist; err: %v", err))
		}
		whitelistMap[denom] = true
	}

	return &FeeDenomWhitelistDecorator{
		whitelistMap: whitelistMap,
		whitelistStr: strings.Join(denomsWhitelist, ","),
	}
}

func (fdd *FeeDenomWhitelistDecorator) allowAll() bool {
	return len(fdd.whitelistMap) == 0
}

func (fdd *FeeDenomWhitelistDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (newCtx sdk.Context, err error) {
	if fdd.allowAll() {
		return next(ctx, tx, simulate)
	}

	feeTx, ok := tx.(sdk.FeeTx)
	if !ok {
		return ctx, errorsmod.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
	}

	feeCoins := feeTx.GetFee()
	for _, coin := range feeCoins {
		if _, found := fdd.whitelistMap[coin.Denom]; !found {
			return ctx, errorsmod.Wrapf(sdkerrors.ErrInvalidCoins,
				"fee denom is not allowed; got: %v, allowed: %v",
				coin.Denom, fdd.whitelistStr)
		}
	}
	return next(ctx, tx, simulate)
}
