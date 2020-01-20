package reputation

import (
	"github.com/persistenceOne/persistenceSDK/modules/hub/reputation/transactions/feedback"
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	feedback.RegisterCodec(cdc)
}

var cdc = codec.New()

func init() {
	RegisterCodec(cdc)
	cdc.Seal()
}
