package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/persistenceCore/x/halving/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Params(context context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(context)
	params := k.GetParams(ctx)
	return &types.QueryParamsResponse{Params: params}, nil
}
