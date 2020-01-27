package bid

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(Message{}, "contract/bid", nil)
}

var cdc = codec.New()

func init() {
	RegisterCodec(cdc)
	cdc.Seal()
}
