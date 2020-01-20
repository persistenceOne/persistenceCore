package contract

import (
	"github.com/persistenceOne/persistenceSDK/modules/hub/contract/transactions/sign"
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	sign.RegisterCodec(cdc)
}

var cdc = codec.New()

func init() {
	RegisterCodec(cdc)
	cdc.Seal()
}
