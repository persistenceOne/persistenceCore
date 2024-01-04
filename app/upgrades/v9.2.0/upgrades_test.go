//go:build genesis_test
// +build genesis_test

package v9_2_0_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/CosmWasm/wasmd/x/wasm"
	dbm "github.com/cometbft/cometbft-db"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/libs/log"
	tmtypes "github.com/cometbft/cometbft/proto/tendermint/types"
	cmttypes "github.com/cometbft/cometbft/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	oracletypes "github.com/persistenceOne/persistence-sdk/v2/x/oracle/types"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/persistenceCore/v11/app"
	v9_2_0 "github.com/persistenceOne/persistenceCore/v11/app/upgrades/v9.2.0"
)

type KeeperTestHelper struct {
	suite.Suite

	// defaults to false,
	// set to true if any method that potentially alters baseapp/abci is used.
	// this controls whether or not we can reuse the app instance, or have to set a new one.
	hasUsedAbci bool
	// defaults to false, set to true if we want to use a new app instance with caching enabled.
	// then on new setup test call, we just drop the current cache.
	// this is not always enabled, because some tests may take a painful performance hit due to CacheKv.
	withCaching bool

	App         *app.Application
	Ctx         sdk.Context
	QueryHelper *baseapp.QueryServiceTestHelper
	TestAccs    []sdk.AccAddress
}

func (s *KeeperTestHelper) Setup() {
	dir, err := os.MkdirTemp("", "persistence-test-home")
	if err != nil {
		panic(fmt.Sprintf("failed creating temporary directory: %v", err))
	}
	s.T().Cleanup(func() { os.RemoveAll(dir); s.withCaching = false })
	s.SetupWithCustomHome(dir)
	s.setupGeneral()
}

// SetupWithCustomHome initializes a new Persistence app with a custom home directory
func (s *KeeperTestHelper) SetupWithCustomHome(dir string) {

	genDoc, err := cmttypes.GenesisDocFromFile("./core-1-genesis.json")
	s.Require().NoError(err)

	consensusParams := genDoc.ConsensusParams.ToProto()

	// Try to disable genesis invariant
	s.App = app.NewApplication(log.NewNopLogger(), dbm.NewMemDB(), nil, true, app.GetEnabledProposals(), NewAppOptionsWithSkipGenesisInvariants(dir, true), []wasm.Option{})

	s.App.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: &consensusParams,
			AppStateBytes:   genDoc.AppState,
		},
	)

}

func (s *KeeperTestHelper) setupGeneral() {
	s.Ctx = s.App.BaseApp.NewContext(false, tmtypes.Header{Height: 1, ChainID: "core-1"})
	if s.withCaching {
		s.Ctx, _ = s.Ctx.CacheContext()
	}
	s.QueryHelper = &baseapp.QueryServiceTestHelper{
		GRPCQueryRouter: s.App.GRPCQueryRouter(),
		Ctx:             s.Ctx,
	}

	s.App.OracleKeeper.SetParams(s.Ctx, oracletypes.DefaultParams())
}

func setConfig() {
	cfg := sdk.GetConfig()

	cfg.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
	cfg.SetCoinType(app.CoinType)
	cfg.SetPurpose(app.Purpose)

	cfg.Seal()
}

func NewAppOptionsWithSkipGenesisInvariants(homePath string, isSkip bool) servertypes.AppOptions {
	return simtestutil.AppOptionsMap{
		flags.FlagHome:                   homePath,
		crisis.FlagSkipGenesisInvariants: isSkip,
		server.FlagInvCheckPeriod:        1, // check invariant every EndBlocker
	}
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestHelper))
}

func (s *KeeperTestHelper) TestInvariant() {
	setConfig()
	s.Setup()

	// After skip invariant in initGenesis, run EndBlocker with default ctx
	s.Require().NotPanics(func() {
		v9_2_0.Fork(s.Ctx, s.App.StakingKeeper)
	})

}
