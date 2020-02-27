package share

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/mapper"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/burn"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/lock"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/send"
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
