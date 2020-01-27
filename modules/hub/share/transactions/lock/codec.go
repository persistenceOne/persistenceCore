package lock

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(Message{}, "share/lock", nil)
}

var cdc = codec.New()

func init() {
	RegisterCodec(cdc)
	cdc.Seal()
}
