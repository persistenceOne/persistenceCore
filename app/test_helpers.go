package app

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"cosmossdk.io/log"
	sdkmath "cosmossdk.io/math"
	pruningtypes "cosmossdk.io/store/pruning/types"
	"github.com/CosmWasm/wasmd/x/wasm"
	abci "github.com/cometbft/cometbft/abci/types"
	cmttypes "github.com/cometbft/cometbft/types"
	dbm "github.com/cosmos/cosmos-db"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/testutil/mock"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/persistenceOne/persistenceCore/v16/app/constants"
	"github.com/stretchr/testify/require"
)

// NewTestNetworkFixture returns a new persistenceCore AppConstructor for network simulation tests.
func NewTestNetworkFixture() network.TestFixture {
	dir, err := os.MkdirTemp("", "persistenceCore")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	defer os.RemoveAll(dir)

	app := NewApplication(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.NewAppOptionsWithFlagHome(dir), []wasm.Option{})

	appCtr := func(val network.ValidatorI) servertypes.Application {
		return NewApplication(
			val.GetCtx().Logger, dbm.NewMemDB(), nil, true,
			simtestutil.NewAppOptionsWithFlagHome(val.GetCtx().Config.RootDir), []wasm.Option{},
			bam.SetPruning(pruningtypes.NewPruningOptionsFromString(val.GetAppConfig().Pruning)),
			bam.SetMinGasPrices(val.GetAppConfig().MinGasPrices),
			bam.SetChainID(val.GetCtx().Viper.GetString(flags.FlagChainID)),
		)
	}

	return network.TestFixture{
		AppConstructor: appCtr,
		GenesisState:   app.DefaultGenesis(),
		EncodingConfig: testutil.TestEncodingConfig{
			InterfaceRegistry: app.InterfaceRegistry(),
			Codec:             app.AppCodec(),
			TxConfig:          app.TxConfig(),
			Amino:             app.LegacyAmino(),
		},
	}
}

func setup(withGenesis bool, t *testing.T) (*Application, GenesisState) {
	db := dbm.NewMemDB()

	appOptions := make(simtestutil.AppOptionsMap, 0)
	appOptions[flags.FlagHome] = t.TempDir()

	app := NewApplication(log.NewNopLogger(), db, nil, true, appOptions, []wasm.Option{})
	if withGenesis {
		return app, app.DefaultGenesis()
	}
	return app, GenesisState{}
}

//// NewTestAppWithCustomOptions initializes a new TestApp with custom options.
//func NewTestAppWithCustomOptions(t *testing.T, isCheckTx bool, options SetupOptions) *Application {
//	t.Helper()
//
//	privVal := mock.NewPV()
//	pubKey, err := privVal.GetPubKey()
//	require.NoError(t, err)
//	// create validator set with single validator
//	validator := cmttypes.NewValidator(pubKey, 1)
//	valSet := cmttypes.NewValidatorSet([]*cmttypes.Validator{validator})
//
//	// generate genesis account
//	senderPrivKey := secp256k1.GenPrivKey()
//	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
//	balance := banktypes.Balance{
//		Address: acc.GetAddress().String(),
//		Coins:   sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdkmath.NewInt(100000000000000))),
//	}
//
//	app := NewApplication(options.Logger, options.DB, nil, true, options.AppOpts)
//	genesisState := app.DefaultGenesis()
//	genesisState, err = simtestutil.GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, []authtypes.GenesisAccount{acc}, balance)
//	require.NoError(t, err)
//
//	if !isCheckTx {
//		// init chain must be called to stop deliverState from being nil
//		stateBytes, err := cmtjson.MarshalIndent(genesisState, "", " ")
//		require.NoError(t, err)
//
//		// Initialize the chain
//		_, err = app.InitChain(&abci.RequestInitChain{
//			Validators:      []abci.ValidatorUpdate{},
//			ConsensusParams: simtestutil.DefaultConsensusParams,
//			AppStateBytes:   stateBytes,
//		})
//		require.NoError(t, err)
//	}
//
//	return app
//}

// Setup initializes a new TestApp. A Nop logger is set in TestApp.
func Setup(t *testing.T) *Application {
	t.Helper()

	privVal := mock.NewPV()
	pubKey, err := privVal.GetPubKey()
	require.NoError(t, err)

	// create validator set with single validator
	validator := cmttypes.NewValidator(pubKey, 1)
	valSet := cmttypes.NewValidatorSet([]*cmttypes.Validator{validator})

	// generate genesis account
	senderPrivKey := secp256k1.GenPrivKey()
	acc := authtypes.NewBaseAccount(senderPrivKey.PubKey().Address().Bytes(), senderPrivKey.PubKey(), 0, 0)
	balance := banktypes.Balance{
		Address: acc.GetAddress().String(),
		Coins:   sdk.NewCoins(sdk.NewCoin(constants.BondDenom, sdkmath.NewInt(100000000000000))),
	}

	app := SetupWithGenesisValSet(t, valSet, []authtypes.GenesisAccount{acc}, balance)

	return app
}

// SetupWithGenesisValSet initializes a new TestApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit in the default token of the TestApp from first genesis
// account. A Nop logger is set in TestApp.
func SetupWithGenesisValSet(t *testing.T, valSet *cmttypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *Application {
	t.Helper()

	app, genesisState := setup(true, t)
	genesisState, err := simtestutil.GenesisStateWithValSet(app.AppCodec(), genesisState, valSet, genAccs, balances...)
	require.NoError(t, err)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	_, err = app.InitChain(&abci.RequestInitChain{
		Validators:      []abci.ValidatorUpdate{},
		ConsensusParams: simtestutil.DefaultConsensusParams,
		AppStateBytes:   stateBytes,
	},
	)
	require.NoError(t, err)

	require.NoError(t, err)
	_, err = app.FinalizeBlock(&abci.RequestFinalizeBlock{
		Height:             app.LastBlockHeight() + 1,
		Hash:               app.LastCommitID().Hash,
		NextValidatorsHash: valSet.Hash(),
	})
	require.NoError(t, err)

	return app
}
