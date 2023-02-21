package wasmbindings

import (
	"encoding/json"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/persistenceOne/persistenceCore/v7/wasmbindings/bindings"
)

// customQuerier dispatches custom CosmWasm bindings queries.
func customQuerier(qp *QueryPlugin) func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var contractQuery bindings.OracleQuery
		if err := json.Unmarshal(request, &contractQuery); err != nil {
			return nil, sdkerrors.Wrap(err, "invalid oracle query")
		}

		switch {
		case contractQuery.GetExchangeRate != nil:
			exchangeRate, err := qp.oracleKeeper.GetExchangeRate(ctx, contractQuery.GetExchangeRate.Symbol)
			if err != nil {
				return nil, status.Error(codes.NotFound, err.Error())
			}

			res := bindings.GetExchangeRateResponse{
				ExchangeRate: exchangeRate.String(),
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "invalid get exchange rate response")
			}

			return bz, nil
		case contractQuery.GetAllExchangeRates != nil:
			var decExchangeRates sdk.DecCoins
			qp.oracleKeeper.IterateExchangeRates(ctx, func(denom string, rate sdk.Dec) (stop bool) {
				decExchangeRates = decExchangeRates.Add(sdk.NewDecCoinFromDec(denom, rate))
				return false
			})

			exchangeRates := make([]string, len(decExchangeRates))
			for i, rate := range decExchangeRates {
				exchangeRates[i] = rate.Amount.String()
			}

			res := bindings.GetAllExchangeRateResponse{
				ExchangeRate: exchangeRates,
			}

			bz, err := json.Marshal(res)
			if err != nil {
				return nil, sdkerrors.Wrap(err, "invalid get all exchange rate response")
			}

			return bz, nil
		default:
			return nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown oracle query variant"}
		}
	}
}
