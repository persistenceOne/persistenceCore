package asset

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceSDK/modules/asset/mapper"
	"github.com/persistenceOne/persistenceSDK/modules/asset/transactions/burn"
	"github.com/persistenceOne/persistenceSDK/modules/asset/transactions/lock"
	"github.com/persistenceOne/persistenceSDK/modules/asset/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/modules/asset/transactions/send"
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
