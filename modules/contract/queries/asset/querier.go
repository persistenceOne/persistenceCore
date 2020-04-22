package asset

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/persistenceSDK/modules/asset/constants"
	"github.com/persistenceOne/persistenceSDK/modules/asset/mapper"
	abciTypes "github.com/tendermint/tendermint/abci/types"
)

type Querier interface {
	Query(sdkTypes.Context, abciTypes.RequestQuery) ([]byte, error)
}

type baseQuerier struct {
	mapper mapper.Mapper
}

func NewQuerier(mapper mapper.Mapper) Querier {
	return baseQuerier{mapper: mapper}
}

var _ Querier = (*baseQuerier)(nil)

func (baseQuerier baseQuerier) Query(context sdkTypes.Context, requestQuery abciTypes.RequestQuery) ([]byte, error) {
	var query query
	if error := packageCodec.UnmarshalJSON(requestQuery.Data, &query); error != nil {
		return nil, errors.Wrap(constants.IncorrectQueryCode, error.Error())
	}
	asset, getAssetError := baseQuerier.mapper.Read(context, mapper.NewAssetAddress(query.Address))
	if getAssetError != nil {
		return nil, getAssetError
	}

	bytes, marshalJSONIndentError := codec.MarshalJSONIndent(packageCodec, asset)
	if marshalJSONIndentError != nil {
		panic(marshalJSONIndentError)
	}

	return bytes, nil
}
