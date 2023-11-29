package app

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestTxFeeChecker(t *testing.T) {
	testcases := []struct {
		name            string
		minGasPricesStr string
		feeCoins        sdk.Coins
		expectedErr     string
		panics          bool
	}{
		{
			name:            "empty min gas prices",
			minGasPricesStr: "",
			feeCoins:        coins(),
			expectedErr:     "minimum gas prices must be set",
			panics:          true,
		},
		{
			name:            "zero min gas prices - empty fee",
			minGasPricesStr: "0token1",
			feeCoins:        coins(),
			expectedErr:     "",
		},
		{
			name:            "zero min gas prices - zero fee",
			minGasPricesStr: "0token1",
			feeCoins:        coins(0, "token1"),
			expectedErr:     "",
		},
		{
			name:            "zero min gas prices - valid fee",
			minGasPricesStr: "0token1",
			feeCoins:        coins(10, "token1"),
			expectedErr:     "",
		},
		{
			name:            "zero min gas prices - unknown fee denom",
			minGasPricesStr: "0token1",
			feeCoins:        coins(1, "token1", 1, "token2"),
			expectedErr:     "fee is not a subset of required fees; got: 1token1,1token2, required: 0token1: insufficient fee",
		},
		{
			name:            "one token min gas prices - empty fee",
			minGasPricesStr: "10token1",
			feeCoins:        coins(),
			expectedErr:     "insufficient fees; got:  required: 10token1: insufficient fee",
		},
		{
			name:            "one token min gas prices - insufficient fee",
			minGasPricesStr: "10token1",
			feeCoins:        coins(5, "token1"),
			expectedErr:     "insufficient fees; got: 5token1 required: 10token1: insufficient fee",
		},
		{
			name:            "one token min gas prices - valid fee",
			minGasPricesStr: "10token1",
			feeCoins:        coins(10, "token1"),
			expectedErr:     "",
		},
		{
			name:            "one token min gas prices - unknown fee denom",
			minGasPricesStr: "10token1",
			feeCoins:        coins(10, "token1", 10, "token2"),
			expectedErr:     "fee is not a subset of required fees; got: 10token1,10token2, required: 10token1: insufficient fee",
		},
		{
			name:            "two tokens min gas prices - empty fee",
			minGasPricesStr: "10token1,10token2",
			feeCoins:        coins(),
			expectedErr:     "insufficient fees; got:  required: 10token1,10token2: insufficient fee",
		},
		{
			name:            "two tokens min gas prices - insufficient fee",
			minGasPricesStr: "10token1,10token2",
			feeCoins:        coins(5, "token1"),
			expectedErr:     "insufficient fees; got: 5token1 required: 10token1,10token2: insufficient fee",
		},
		{
			name:            "two tokens min gas prices - insufficient fee",
			minGasPricesStr: "10token1,10token2",
			feeCoins:        coins(5, "token1", 5, "token2"),
			expectedErr:     "insufficient fees; got: 5token1,5token2 required: 10token1,10token2: insufficient fee",
		},
		{
			name:            "two tokens min gas prices - valid fee (token1)",
			minGasPricesStr: "10token1,10token2",
			feeCoins:        coins(10, "token1"),
			expectedErr:     "",
		},
		{
			name:            "two tokens min gas prices - valid fee (token2)",
			minGasPricesStr: "10token1,10token2",
			feeCoins:        coins(10, "token2"),
			expectedErr:     "",
		},
		{
			name:            "two tokens min gas prices - valid fee (token1 + partial token2)",
			minGasPricesStr: "10token1,10token2",
			feeCoins:        coins(10, "token1", 5, "token2"),
			expectedErr:     "",
		},
		{
			name:            "two tokens min gas prices - valid fee (partial token1 + token2)",
			minGasPricesStr: "10token1,10token2",
			feeCoins:        coins(5, "token1", 10, "token2"),
			expectedErr:     "",
		},
		{
			name:            "two tokens min gas prices - valid fee (token1 + token2)",
			minGasPricesStr: "10token1,10token2",
			feeCoins:        coins(10, "token1", 10, "token2"),
			expectedErr:     "",
		},
		{
			name:            "two tokens min gas prices - unknown fee denom",
			minGasPricesStr: "10token1,10token2",
			feeCoins:        coins(10, "token1", 10, "token3"),
			expectedErr:     "fee is not a subset of required fees; got: 10token1,10token3, required: 10token1,10token2: insufficient fee",
		},
		{
			// TODO(ajeet): should this test not return error?
			name:            "multi token min gas prices with one zero token - empty fee",
			minGasPricesStr: "0token1,10token2",
			feeCoins:        coins(),
			expectedErr:     "insufficient fees; got:  required: 0token1,10token2: insufficient fee",
		},
		{
			// TODO(ajeet): should this test not return error?
			name:            "multi token min gas prices with one zero token - fee with zero token denom",
			minGasPricesStr: "0token1,10token2",
			feeCoins:        coins(10, "token1"),
			expectedErr:     "insufficient fees; got: 10token1 required: 0token1,10token2: insufficient fee",
		},
		{
			name:            "multi token min gas prices with one zero token - valid fee",
			minGasPricesStr: "0token1,10token2",
			feeCoins:        coins(10, "token2"),
			expectedErr:     "",
		},
		{
			name:            "multi token min gas prices with one zero token - valid fee",
			minGasPricesStr: "0token1,10token2",
			feeCoins:        coins(5, "token1", 10, "token2"),
			expectedErr:     "",
		},
		{
			name:            "all zero token min gas prices - empty fee",
			minGasPricesStr: "0token1,0token2",
			feeCoins:        coins(),
			expectedErr:     "",
		},
		{
			name:            "all zero token min gas prices - valid fee",
			minGasPricesStr: "0token1,0token2",
			feeCoins:        coins(10, "token1"),
			expectedErr:     "",
		},
		{
			name:            "all zero token min gas prices - valid fee",
			minGasPricesStr: "0token1,0token2",
			feeCoins:        coins(10, "token2"),
			expectedErr:     "",
		},
		{
			name:            "all zero token min gas prices - unknow fee denom",
			minGasPricesStr: "0token1,0token2",
			feeCoins:        coins(5, "token3"),
			expectedErr:     "fee is not a subset of required fees; got: 5token3, required: 0token1,0token2: insufficient fee",
		},
	}

	isCheckTx := true
	ctx := sdk.NewContext(nil, types.Header{}, isCheckTx, log.NewNopLogger())

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.panics {
					require.Contains(t, r, tc.expectedErr)
				} else {
					require.Nil(t, r, "Test did not panic")
				}
			}()

			txFeeChecker := GetTxFeeChecker(tc.minGasPricesStr)
			tx := &mockFeeTx{fee: tc.feeCoins}
			_, _, err := txFeeChecker(ctx, tx)
			if tc.expectedErr == "" {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}

func coins(amountDenomPairs ...interface{}) sdk.Coins {
	coins := sdk.Coins{}
	for i := 0; i < len(amountDenomPairs); i += 2 {
		coins = append(coins, sdk.Coin{
			Amount: math.NewInt(int64(amountDenomPairs[i].(int))),
			Denom:  amountDenomPairs[i+1].(string),
		})
	}
	return coins
}

var _ sdk.Tx = (*mockFeeTx)(nil)
var _ sdk.FeeTx = (*mockFeeTx)(nil)

type mockFeeTx struct{ fee sdk.Coins }

func (m *mockFeeTx) GetFee() sdk.Coins { return m.fee }
func (m *mockFeeTx) GetGas() uint64    { return 1 }

func (*mockFeeTx) FeeGranter() sdk.AccAddress { return nil }
func (*mockFeeTx) FeePayer() sdk.AccAddress   { return nil }
func (*mockFeeTx) GetMsgs() []sdk.Msg         { return nil }
func (*mockFeeTx) ValidateBasic() error       { return nil }
