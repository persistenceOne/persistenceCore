package asset

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/mapper"
)

func RegisterCodec(codec *codec.Codec) {
	codec.RegisterConcrete(query{}, "asset/query", nil)
}

var packageCodec = codec.New()

func init() {
	RegisterCodec(packageCodec)
	mapper.RegisterCodec(packageCodec)
	packageCodec.Seal()
}
