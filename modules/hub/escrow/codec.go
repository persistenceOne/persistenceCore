package escrow

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceSDK/modules/hub/escrow/transactions/execute"
)

func RegisterCodec(codec *codec.Codec) {
	execute.RegisterCodec(codec)
}

var packageCodec = codec.New()

func init() {
	RegisterCodec(packageCodec)
	packageCodec.Seal()
}
