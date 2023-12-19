package app

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/libs/log"
	tmtypes "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestFeeDenomWhiltelistDecorator(t *testing.T) {
	testDenomsWhitelist := []string{
		"uxprt",
		"ibc/C8A74ABBE2AF892E15680D916A7C22130585CE5704F9B17A10F184A90D53BECA",
	}

	testcases := []struct {
		name            string
		txFee           sdk.Coins
		denomsWhitelist []string
		expectedErr     string
		panics          bool
	}{
		{
			name:            "empty denoms list should allow all denom",
			txFee:           coins(1, "abcd", 1, "xyz"),
			denomsWhitelist: []string{},
			expectedErr:     "",
		},
		{
			name:            "invalid denoms list",
			txFee:           coins(),
			denomsWhitelist: []string{"a"},
			expectedErr:     "invalid denoms whiltelist; err: invalid denom: a",
			panics:          true,
		},
		{
			name:            "valid denoms - valid fee",
			txFee:           coins(10, "uxprt"),
			denomsWhitelist: testDenomsWhitelist,
			expectedErr:     "",
		},
		{
			name:            "valid denoms - multiple valid fees",
			txFee:           coins(10, testDenomsWhitelist[1], 10, "uxprt"),
			denomsWhitelist: testDenomsWhitelist,
			expectedErr:     "",
		},
		{
			name:            "valid denoms - invalid fee",
			txFee:           coins(10, "abcd"),
			denomsWhitelist: testDenomsWhitelist,
			expectedErr:     "fee denom is not allowed; got: abcd, allowed: uxprt,ibc/C8A74ABBE2AF892E15680D916A7C22130585CE5704F9B17A10F184A90D53BECA: invalid coins",
		},
		{
			name:            "valid denoms - multiple invalid fee",
			txFee:           coins(10, "uxprt", 10, "xyz"),
			denomsWhitelist: testDenomsWhitelist,
			expectedErr:     "fee denom is not allowed; got: xyz, allowed: uxprt,ibc/C8A74ABBE2AF892E15680D916A7C22130585CE5704F9B17A10F184A90D53BECA: invalid coins",
		},
	}

	isCheckTx := false
	ctx := sdk.NewContext(nil, tmtypes.Header{}, isCheckTx, log.NewNopLogger())
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				r := recover()
				if tc.panics {
					require.Equal(t, r, tc.expectedErr)
				} else {
					require.Nil(t, r, "Test did not panic")
				}
			}()

			fdd := NewFeeDenomWhitelistDecorator(tc.denomsWhitelist)
			antehandlerFFD := sdk.ChainAnteDecorators(fdd)
			tx := &mockFeeTx{fee: tc.txFee}

			_, err := antehandlerFFD(ctx, tx, false)

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
