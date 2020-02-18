package lock

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(codec *codec.Codec) {
	codec.RegisterConcrete(Message{}, "share/lock", nil)
}

var packageCodec = codec.New()

func init() {
	RegisterCodec(packageCodec)
	packageCodec.Seal()
}
