package asset

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/mapper"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/mint"
)

func RegisterCodec(codec *codec.Codec) {
	mint.RegisterCodec(codec)
	mapper.RegisterCodec(codec)
}

var moduleCodec = codec.New()

func init() {
	RegisterCodec(moduleCodec)
	moduleCodec.Seal()
}
