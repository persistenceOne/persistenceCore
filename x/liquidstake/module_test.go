package liquidstake_test

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceCore/v17/app/constants"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/persistenceCore/v17/app"
	"github.com/persistenceOne/persistenceCore/v17/x/liquidstake"
	"github.com/persistenceOne/persistenceCore/v17/x/liquidstake/types"
)

type ModuleTestSuite struct {
	suite.Suite

	app         *app.Application
	ctx         sdk.Context
	appModule   liquidstake.AppModule
	basicModule liquidstake.AppModuleBasic
	cdc         codec.Codec
}

func TestModuleTestSuite(t *testing.T) {
	suite.Run(t, new(ModuleTestSuite))
}

func (s *ModuleTestSuite) SetupTest() {
	constants.SetUnsealedConfig()

	s.app = app.Setup(s.T())
	s.ctx = s.app.BaseApp.NewContext(false)

	blocktime, err := sdk.ParseTime("2022-03-01T00:00:00.000000000")
	s.Require().NoError(err)
	s.ctx = s.ctx.WithBlockHeight(100).WithBlockTime(blocktime)
	s.appModule = liquidstake.NewAppModule(*s.app.LiquidStakeKeeper)
	s.basicModule = liquidstake.AppModuleBasic{}
	s.cdc = s.app.AppCodec()
}

func (s *ModuleTestSuite) TestAppModuleBasic() {
	// Test Name
	s.Require().Equal("liquidstake", s.basicModule.Name())

	// Test RegisterLegacyAminoCodec
	cdc := codec.NewLegacyAmino()
	s.basicModule.RegisterLegacyAminoCodec(cdc)

	// Test DefaultGenesis
	defaultGenesis := s.basicModule.DefaultGenesis(s.cdc)
	var genesisState types.GenesisState
	s.Require().NoError(s.cdc.UnmarshalJSON(defaultGenesis, &genesisState))

	// Test ValidateGenesis
	s.Require().NoError(s.basicModule.ValidateGenesis(s.cdc, nil, defaultGenesis))

	// Test RegisterInterfaces
	registry := codectypes.NewInterfaceRegistry()
	s.basicModule.RegisterInterfaces(registry)

	// Test GetTxCmd
	s.Require().NotNil(s.basicModule.GetTxCmd())
}

func (s *ModuleTestSuite) TestAppModule() {
	// Test Name
	s.Require().Equal("liquidstake", s.appModule.Name())

	// Test QuerierRoute
	s.Require().Equal("liquidstake", s.appModule.QuerierRoute())

	// Test RegisterInvariants
	// This is a no-op, just call it for coverage
	s.appModule.RegisterInvariants(nil)

	// Skip RegisterServices test as it requires a fully configured router
	// which is difficult to mock in unit tests

	// Test InitGenesis
	defaultGenesis := s.basicModule.DefaultGenesis(s.cdc)
	var genesisState types.GenesisState
	s.Require().NoError(s.cdc.UnmarshalJSON(defaultGenesis, &genesisState))
	validatorUpdates := s.appModule.InitGenesis(s.ctx, s.cdc, defaultGenesis)
	s.Require().Len(validatorUpdates, 0)

	// Test ExportGenesis
	exportedGenesis := s.appModule.ExportGenesis(s.ctx, s.cdc)
	var exportedState types.GenesisState
	s.Require().NoError(s.cdc.UnmarshalJSON(exportedGenesis, &exportedState))

	// Test ConsensusVersion
	s.Require().Equal(uint64(1), s.appModule.ConsensusVersion())

	// Test BeginBlock
	s.appModule.BeginBlock(s.ctx)
}

// Test invalid genesis state
func TestValidateGenesisFail(t *testing.T) {
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	basicModule := liquidstake.AppModuleBasic{}

	// Create invalid genesis state with negative unstake fee rate
	invalidGenesis := types.GenesisState{
		Params: types.Params{
			UnstakeFeeRate: math.LegacyNewDecWithPrec(-1, 2), // -0.01
		},
	}
	invalidGenesisBz, err := cdc.MarshalJSON(&invalidGenesis)
	require.NoError(t, err)

	// Validate should fail
	err = basicModule.ValidateGenesis(cdc, nil, invalidGenesisBz)
	require.Error(t, err)
}

// Test invalid JSON in genesis state
func TestInvalidJSONGenesis(t *testing.T) {
	cdc := codec.NewProtoCodec(codectypes.NewInterfaceRegistry())
	basicModule := liquidstake.AppModuleBasic{}

	// Invalid JSON
	invalidJSON := []byte(`{"params": {invalid json}}`)

	// Validate should fail
	err := basicModule.ValidateGenesis(cdc, nil, invalidJSON)
	require.Error(t, err)
}
