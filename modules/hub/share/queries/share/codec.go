package share

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/mapper"
	"github.com/persistenceOne/persistenceSDK/types"
)

func RegisterCodec(codec *codec.Codec) {
	codec.RegisterConcrete(query{}, "share/query", nil)
}

var packageCodec = codec.New()

func init() {
	RegisterCodec(packageCodec)
	mapper.RegisterCodec(packageCodec)
	types.RegisterCodec(packageCodec)
	packageCodec.Seal()
}
