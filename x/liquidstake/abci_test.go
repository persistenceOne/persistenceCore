package liquidstake_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceCore/v16/app"
	"github.com/persistenceOne/persistenceCore/v16/app/constants"
	"github.com/stretchr/testify/suite"

	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake"
	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/keeper"
)

type ABCITestSuite struct {
	suite.Suite

	app    *app.Application
	ctx    sdk.Context
	keeper keeper.Keeper
}

func TestABCITestSuite(t *testing.T) {
	suite.Run(t, new(ABCITestSuite))
}

func (s *ABCITestSuite) SetupTest() {
	constants.SetUnsealedConfig()
	s.app = app.Setup(s.T())
	s.ctx = s.app.BaseApp.NewContext(false)

	s.keeper = *s.app.LiquidStakeKeeper
	blocktime, err := sdk.ParseTime("2022-03-01T00:00:00.000000000")
	s.Require().NoError(err)
	s.ctx = s.ctx.WithBlockHeight(100).WithBlockTime(blocktime)
}

func (s *ABCITestSuite) TestBeginBlock() {
	// Test when module is not paused
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	params.ModulePaused = false
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))

	// Call BeginBlock
	err = liquidstake.BeginBlock(s.ctx, s.keeper)
	s.Require().NoError(err)

	// Test when module is paused
	params.ModulePaused = true
	s.Require().NoError(s.keeper.SetParams(s.ctx, params))

	// Call BeginBlock
	err = liquidstake.BeginBlock(s.ctx, s.keeper)
	s.Require().NoError(err)
}
