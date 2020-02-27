package share

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceSDK/modules/share/mapper"
	"github.com/persistenceOne/persistenceSDK/modules/share/transactions/burn"
	"github.com/persistenceOne/persistenceSDK/modules/share/transactions/lock"
	"github.com/persistenceOne/persistenceSDK/modules/share/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/modules/share/transactions/send"
)

func RegisterCodec(codec *codec.Codec) {
	mapper.RegisterCodec(codec)

	burn.RegisterCodec(codec)
	lock.RegisterCodec(codec)
	mint.RegisterCodec(codec)
	send.RegisterCodec(codec)
}

var packageCodec = codec.New()

func init() {
	RegisterCodec(packageCodec)
	packageCodec.Seal()
}
