package share

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/share/mapper"
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

func (baseQuerier baseQuerier) Query(context sdkTypes.Context, requestQuery abciTypes.RequestQuery) ([]byte, sdkTypes.Error) {
	var query query
	if error := packageCodec.UnmarshalJSON(requestQuery.Data, &query); error != nil {
		return nil, incorrectQueryError(error.Error())
	}
	share, getShareError := baseQuerier.mapper.Read(context, mapper.NewShareAddress(query.Address))
	if getShareError != nil {
		return nil, getShareError
	}

	bytes, marshalJSONIndentError := codec.MarshalJSONIndent(packageCodec, share)
	if marshalJSONIndentError != nil {
		panic(marshalJSONIndentError)
	}

	return bytes, nil
}
