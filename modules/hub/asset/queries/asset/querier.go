package asset

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/mapper"
	abciTypes "github.com/tendermint/tendermint/abci/types"
)

type Querier interface {
	Query(sdkTypes.Context, abciTypes.RequestQuery) ([]byte, sdkTypes.Error)
}

type baseQuerier struct {
	mapper mapper.Mapper
}

func NewQuerier(mapper mapper.Mapper) Querier {
	return baseQuerier{mapper: mapper}
}

var _ Querier = (*baseQuerier)(nil)

type query struct {
	Address string
}

func (baseQuerier baseQuerier) Query(context sdkTypes.Context, requestQuery abciTypes.RequestQuery) ([]byte, sdkTypes.Error) {
	var query query
	if error := packageCodec.UnmarshalJSON(requestQuery.Data, &query); error != nil {
		return nil, IncorrectQueryError(error.Error())
	}
	asset, getAssetError := baseQuerier.mapper.GetAsset(context, mapper.NewAssetAddress(query.Address))
	if getAssetError != nil {
		return nil, getAssetError
	}

	bytes, marshalJSONIndentError := codec.MarshalJSONIndent(packageCodec, asset)
	if marshalJSONIndentError != nil {
		panic(marshalJSONIndentError)
	}

	return bytes, nil
}
