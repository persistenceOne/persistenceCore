package main

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	oracletypes "github.com/persistenceOne/persistence-sdk/v3/x/oracle/types"
	buildertypes "github.com/skip-mev/pob/x/builder/types"
)

func (s *TestSuite) VerifyParams() {
	ctx := context.Background()
	client := s.GetChainClient("test-core-2").Client

	s.T().Log("verify min commission rate & LSM params")
	params, err := stakingtypes.NewQueryClient(client).Params(ctx, &stakingtypes.QueryParamsRequest{})
	s.Require().NoError(err)
	s.Require().Equal(sdk.NewDecWithPrec(5, 2), params.Params.MinCommissionRate)
	s.Require().Equal(sdk.NewDec(250), params.Params.ValidatorBondFactor)
	s.Require().Equal(sdk.NewDecWithPrec(1, 1), params.Params.GlobalLiquidStakingCap)
	s.Require().Equal(sdk.NewDecWithPrec(5, 1), params.Params.ValidatorLiquidStakingCap)

	s.T().Log("verify mev aution is disabled")
	pobParams, err := buildertypes.NewQueryClient(client).Params(ctx, &buildertypes.QueryParamsRequest{})
	s.Require().NoError(err)
	s.Require().Zero(pobParams.Params.MaxBundleSize)

	s.T().Log("verify oracle is disabled")
	oracleParams, err := oracletypes.NewQueryClient(client).Params(ctx, &oracletypes.QueryParamsRequest{})
	s.Require().NoError(err)
	s.Require().Empty(oracleParams.Params.AcceptList)

	s.T().Log("verify ibc allowed clients includes localhost")
	ibcClientParams, err := ibcclienttypes.NewQueryClient(client).ClientParams(ctx, &ibcclienttypes.QueryClientParamsRequest{})
	s.Require().NoError(err)
	s.Require().Contains(ibcClientParams.Params.AllowedClients, exported.Localhost)

	s.T().Log("verify gov MinInitialDepositParam")
	govParams, err := govtypes.NewQueryClient(client).Params(ctx, &govtypes.QueryParamsRequest{})
	s.Require().NoError(err)
	s.Require().Equal("0.25", govParams.Params.MinInitialDepositRatio)
}
