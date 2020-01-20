package escrow

import (
	"github.com/persistenceOne/persistenceSDK/modules/hub/escrow/transactions/execute"
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	execute.RegisterCodec(cdc)
}

var cdc = codec.New()

func init() {
	RegisterCodec(cdc)
	cdc.Seal()
}
