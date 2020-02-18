package share

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/mint"
)

func RegisterCodec(codec *codec.Codec) {
	mint.RegisterCodec(codec)
}

var packageCodec = codec.New()

func init() {
	RegisterCodec(packageCodec)
	packageCodec.Seal()
}
