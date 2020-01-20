package share

import (
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/mint"
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	mint.RegisterCodec(cdc)
}

var cdc = codec.New()

func init() {
	RegisterCodec(cdc)
	cdc.Seal()
}
