package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/persistenceOne/persistenceCore/v17/x/liquidstake/types"
)

func (s *KeeperTestSuite) TestGRPCParams() {
	resp, err := s.querier.Params(sdk.WrapSDKContext(s.ctx), &types.QueryParamsRequest{})
	s.Require().NoError(err)
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(params, resp.Params)
}

func (s *KeeperTestSuite) TestGRPCQueries() {
	_, valOpers, _ := s.CreateValidators([]int64{1000000, 2000000, 3000000})
	params, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	params.MinLiquidStakeAmount = math.NewInt(50000)
	err = s.keeper.SetParams(s.ctx, params)
	s.Require().NoError(err)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	// add active validator
	params.WhitelistedValidators = []types.WhitelistedValidator{
		{ValidatorAddress: valOpers[0].String(), TargetWeight: math.NewInt(3334)},
		{ValidatorAddress: valOpers[1].String(), TargetWeight: math.NewInt(3333)},
		{ValidatorAddress: valOpers[2].String(), TargetWeight: math.NewInt(3333)},
	}
	s.keeper.SetParams(s.ctx, params)
	s.keeper.UpdateLiquidValidatorSet(s.ctx, true)

	// Test LiquidValidators grpc query
	res, err := s.keeper.GetAllLiquidValidatorStates(s.ctx)
	s.Require().NoError(err)
	resp, err := s.querier.LiquidValidators(sdk.WrapSDKContext(s.ctx), &types.QueryLiquidValidatorsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(resp.LiquidValidators, res)

	resp, err = s.querier.LiquidValidators(sdk.WrapSDKContext(s.ctx), nil)
	s.Require().Nil(resp)
	s.Require().ErrorIs(err, status.Error(codes.InvalidArgument, "invalid request"))

	// Test States grpc query
	respStates, err := s.querier.States(sdk.WrapSDKContext(s.ctx), &types.QueryStatesRequest{})
	s.Require().NoError(err)
	resNetAmountState, err := s.keeper.GetNetAmountState(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(respStates.NetAmountState, resNetAmountState)

	respStates, err = s.querier.States(sdk.WrapSDKContext(s.ctx), nil)
	s.Require().Nil(respStates)
	s.Require().ErrorIs(err, status.Error(codes.InvalidArgument, "invalid request"))

	// Test Params grpc query
	respParams, err := s.querier.Params(s.ctx, &types.QueryParamsRequest{})
	s.Require().NoError(err)
	resParams, err := s.keeper.GetParams(s.ctx)
	s.Require().NoError(err)
	s.Require().Equal(respParams.Params.LiquidBondDenom, resParams.LiquidBondDenom)
	s.Require().Equal(respParams.Params.WhitelistedValidators[0].ValidatorAddress, valOpers[0].String())
	s.Require().Equal(respParams.Params.WhitelistedValidators[1].ValidatorAddress, valOpers[1].String())
	s.Require().Equal(respParams.Params.WhitelistedValidators[2].ValidatorAddress, valOpers[2].String())
}
