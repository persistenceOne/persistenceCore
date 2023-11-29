package app

import (
	"fmt"
	"math"
	"strings"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
)

// GetTxFeeChecker implements sdk's default behaviour
// with an additional check that fee denoms must be subset of
// min required fee denoms.
// see: https://github.com/persistenceOne/cosmos-sdk/blob/v0.47.3-lsm5/x/auth/ante/validator_tx_fee.go#L13-L46
func GetTxFeeChecker(minGasPricesStr string) ante.TxFeeChecker {
	// ctx.MinGasPrices() does not return zero coins as they were removed while parsing
	// see: https://github.com/persistenceOne/cosmos-sdk/blob/v0.47.3-lsm5/baseapp/options.go#L28
	// This behaviour is considered valid: https://github.com/cosmos/cosmos-sdk/issues/17755#issuecomment-1721493256
	// therefore parsing (again) min gas prices from the config string.
	minGasPrices := mustParseMinGasPrices(minGasPricesStr)

	return func(ctx sdk.Context, tx sdk.Tx) (sdk.Coins, int64, error) {
		feeTx, ok := tx.(sdk.FeeTx)
		if !ok {
			return nil, 0, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "Tx must be a FeeTx")
		}

		feeCoins := feeTx.GetFee()
		gas := feeTx.GetGas()

		// Ensure that the provided fees meet a minimum threshold for the validator,
		// if this is a CheckTx. This is only for local mempool purposes, and thus
		// is only ran on check tx.
		if ctx.IsCheckTx() && len(minGasPrices) != 0 {
			requiredFees := make(sdk.Coins, len(minGasPrices))

			// Determine the required fees by multiplying each required minimum gas
			// price by the gas limit, where fee = ceil(minGasPrice * gasLimit).
			glDec := sdkmath.LegacyNewDec(int64(gas))
			for i, gp := range minGasPrices {
				fee := gp.Amount.Mul(glDec)
				requiredFees[i] = sdk.NewCoin(gp.Denom, fee.Ceil().RoundInt())
			}

			if !denomsSubsetOf(feeCoins, requiredFees) {
				return nil, 0, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee,
					"fee is not a subset of required fees; got: %s, required: %s",
					feeCoins, requiredFees)
			}

			// TODO(ajeet): should we allow zero fee if minimum-gas-prices = "0token1,0.25token2"???
			// (currently, any non zero token (if defined) must be present in the fee)
			// If yes, then !amt.IsZero() check can be removed from IsAnyGTE method.
			if !minGasPrices.IsZero() && !feeCoins.IsAnyGTE(requiredFees) {
				return nil, 0, sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "insufficient fees; got: %s required: %s", feeCoins, requiredFees)
			}
		}

		priority := getTxPriority(feeCoins, int64(gas))
		return feeCoins, priority, nil
	}
}

// getTxPriority returns a naive tx priority based on the amount of the smallest denomination of the gas price
// provided in a transaction.
// NOTE: This implementation should be used with a great consideration as it opens potential attack vectors
// where txs with multiple coins could not be prioritize as expected.
//
// This piece of code is copied from sdk's default fee checkers.
// see: https://github.com/persistenceOne/cosmos-sdk/blob/v0.47.3-lsm5/x/auth/ante/validator_tx_fee.go#L48-L66
func getTxPriority(fee sdk.Coins, gas int64) int64 {
	var priority int64
	for _, c := range fee {
		p := int64(math.MaxInt64)
		gasPrice := c.Amount.QuoRaw(gas)
		if gasPrice.IsInt64() {
			p = gasPrice.Int64()
		}
		if priority == 0 || p < priority {
			priority = p
		}
	}

	return priority
}

// mustParseMinGasPrices is similar to sdk.ParseDecCoins except
// it does not sanitizes zero coins.
func mustParseMinGasPrices(coinsStr string) sdk.DecCoins {
	coinsStr = strings.TrimSpace(coinsStr)
	if len(coinsStr) == 0 {
		// Should not be reachable
		panic("minimum gas prices must be set")
	}

	coinStrs := strings.Split(coinsStr, ",")
	decCoins := make(sdk.DecCoins, len(coinStrs))
	for i, coinStr := range coinStrs {
		coin, err := sdk.ParseDecCoin(coinStr)
		if err != nil {
			panic(fmt.Sprintf("invalid minimum gas prices: %v", err))
		}

		decCoins[i] = coin
	}

	return decCoins.Sort()
}

// denomsSubsetOf is similar to sdk's Coins.DenomsSubsetOf() except
// it allows zero coin in the requiredFees set.
// Ex.
// feeCoins = [1 coinA], requiredFees = [0 coinA, 1 coinB] -> returns true
func denomsSubsetOf(feeCoins, requiredFees sdk.Coins) bool {
	if len(feeCoins) > len(requiredFees) {
		return false
	}

	// NOTE: empty set (len(feeCoins) == 0) is subset of anything.
	for _, feeCoin := range feeCoins {
		if found, _ := requiredFees.Find(feeCoin.Denom); !found {
			return false
		}
	}
	return true
}
