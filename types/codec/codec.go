package codec

import (
	"github.com/cosmos/cosmos-sdk/codec"
	types "github.com/persistenceOne/persistenceSDK/interfaces"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*types.Asset)(nil), nil)
	cdc.RegisterInterface((*types.AssetAddress)(nil), nil)

	//cdc.RegisterConcrete(&ModuleAccount{}, "cosmos-sdk/ModuleAccount", nil)
}
