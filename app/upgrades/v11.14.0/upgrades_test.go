package v11_14_0_test

import (
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm"
	dbm "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v12/app"
	v11_14_0 "github.com/persistenceOne/persistenceCore/v12/app/upgrades/v11.14.0"
)

func TestRemoveStargazeUnbondedBalance(t *testing.T) {
	testApp := app.NewApplication(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.NewAppOptionsWithFlagHome(""), []wasm.Option{})
	ctx := testApp.NewContext(true, tmproto.Header{})
	v11_14_0.RemoveStargazeUnbondedBalance(ctx, testApp.LiquidStakeIBCKeeper)

	unbonding, ok := testApp.LiquidStakeIBCKeeper.GetUnbonding(ctx, "stargaze-1", 582)
	require.True(t, ok)
	require.Equal(t, types.NewInt(62810179898), unbonding.UnbondAmount.Amount)

	userUnbonding, ok := testApp.LiquidStakeIBCKeeper.GetUserUnbonding(ctx, "stargaze-1", "persistence1fp6qhht94pmfdq9h94dvw0tnmnlf2vutnlu7pt", 582)
	require.True(t, ok)
	require.Equal(t, types.NewInt(62810179898), userUnbonding.UnbondAmount.Amount)
}
