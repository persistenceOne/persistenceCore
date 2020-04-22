package asset

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceSDK/modules/asset/mapper"
	"github.com/persistenceOne/persistenceSDK/types"
)

func RegisterCodec(codec *codec.Codec) {
	codec.RegisterConcrete(query{}, "asset/query", nil)
}

var packageCodec = codec.New()

func init() {
	RegisterCodec(packageCodec)
	mapper.RegisterCodec(packageCodec)
	types.RegisterCodec(packageCodec)
	packageCodec.Seal()
}
