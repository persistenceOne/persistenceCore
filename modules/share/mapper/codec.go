package mapper

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(codec *codec.Codec) {
	codec.RegisterConcrete(&baseShare{}, "share/baseShare", nil)
	codec.RegisterConcrete(&baseShareAddress{}, "share/baseShareAddress", nil)
}
