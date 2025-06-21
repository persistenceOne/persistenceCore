package v11_16_0_test

import (
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v12/app"
	v11_16_0 "github.com/persistenceOne/persistenceCore/v12/app/upgrades/v11.16.0"
)

func TestRemoveUnbondedBalance(t *testing.T) {
	testApp := app.NewApplication(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.NewAppOptionsWithFlagHome(""), []wasm.Option{})
	ctx := testApp.NewContext(true, tmproto.Header{})
	v11_16_0.RemoveUnbondedBalance(ctx, testApp.LiquidStakeIBCKeeper,
		"chihuahua-1", 704, sdk.NewInt64Coin("stk/uhuahua", 5000000), sdk.NewInt64Coin("uhuahua", 5010101), "persistence1fp6qhht94pmfdq9h94dvw0tnmnlf2vutnlu7pt")

	unbonding, ok := testApp.LiquidStakeIBCKeeper.GetUnbonding(ctx, "chihuahua-1", 704)
	require.True(t, ok)
	require.Equal(t, types.NewInt(5010101), unbonding.UnbondAmount.Amount)

	userUnbonding, ok := testApp.LiquidStakeIBCKeeper.GetUserUnbonding(ctx, "chihuahua-1", "persistence1fp6qhht94pmfdq9h94dvw0tnmnlf2vutnlu7pt", 704)
	require.True(t, ok)
	require.Equal(t, types.NewInt(5010101), userUnbonding.UnbondAmount.Amount)
}
