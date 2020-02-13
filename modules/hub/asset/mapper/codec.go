package mapper

import "github.com/cosmos/cosmos-sdk/codec"

func RegisterCodec(codec *codec.Codec) {
	codec.RegisterConcrete(&baseAsset{}, "asset/baseAsset", nil)
	codec.RegisterConcrete(&baseAssetAddress{}, "asset/baseAssetAddress", nil)
}
