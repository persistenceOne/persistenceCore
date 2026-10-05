package v16_1_0_test

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/log"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/CosmWasm/wasmd/x/wasm"
	abci "github.com/cometbft/cometbft/abci/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v16/app"
	"github.com/persistenceOne/persistenceCore/v16/app/constants"
	v1610 "github.com/persistenceOne/persistenceCore/v16/app/upgrades/v16.1.0"
)

func TestUpgrade(t *testing.T) {
	testApp, ctx := setupTestApp(t)

	err := testApp.UpgradeKeeper.ApplyUpgrade(ctx, upgradetypes.Plan{
		Name:   v1610.UpgradeName,
		Height: 1,
	})
	require.NoError(t, err)
}

func setupTestApp(t *testing.T) (*app.Application, sdk.Context) {
	t.Helper()
	constants.SetConfig()

	testApp := app.NewApplication(
		log.NewNopLogger(),
		dbm.NewMemDB(),
		nil,
		true,
		simtestutil.NewAppOptionsWithFlagHome(t.TempDir()),
		[]wasm.Option{},
	)
	ctx := testApp.NewContext(true)

	validatorSet, err := simtestutil.CreateRandomValidatorSet()
	require.NoError(t, err)

	privateKey := secp256k1.GenPrivKey()
	account := authtypes.NewBaseAccount(
		privateKey.PubKey().Address().Bytes(),
		privateKey.PubKey(),
		0,
		0,
	)
	balance := banktypes.Balance{
		Address: account.GetAddress().String(),
		Coins: sdk.NewCoins(
			sdk.NewCoin(sdk.DefaultBondDenom, sdk.DefaultPowerReduction.MulRaw(100)),
		),
	}
	genesisState, err := simtestutil.GenesisStateWithValSet(
		testApp.AppCodec(),
		testApp.DefaultGenesis(),
		validatorSet,
		[]authtypes.GenesisAccount{account},
		balance,
	)
	require.NoError(t, err)

	genesisStateBytes, err := json.Marshal(genesisState)
	require.NoError(t, err)
	_, err = testApp.InitChainer(ctx, &abci.RequestInitChain{AppStateBytes: genesisStateBytes})
	require.NoError(t, err)

	return testApp, ctx
}
