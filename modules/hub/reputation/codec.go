package reputation

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceSDK/modules/hub/reputation/transactions/feedback"
)

func RegisterCodec(codec *codec.Codec) {
	feedback.RegisterCodec(codec)
}

var packageCodec = codec.New()

func init() {
	RegisterCodec(packageCodec)
	packageCodec.Seal()
}
