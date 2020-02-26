package mapper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

var storeKeyPrefix = []byte{0x11}

func storeKey(assetAddress types.AssetAddress) []byte {
	return append(storeKeyPrefix, assetAddress.Bytes()...)
}

type Mapper interface {
	Get(sdkTypes.Context, types.AssetAddress) (types.Asset, sdkTypes.Error)
	Set(sdkTypes.Context, types.Asset) sdkTypes.Error
	Delete(sdkTypes.Context, types.AssetAddress)
}

type baseMapper struct {
	storeKey sdkTypes.StoreKey
	codec    *codec.Codec
}

func NewMapper(codec *codec.Codec, storeKey sdkTypes.StoreKey) Mapper {
	return baseMapper{
		storeKey: storeKey,
		codec:    codec,
	}
}

var _ Mapper = (*baseMapper)(nil)

func (baseMapper baseMapper) Get(context sdkTypes.Context, assetAddress types.AssetAddress) (asset types.Asset, error sdkTypes.Error) {
	kvStore := context.KVStore(baseMapper.storeKey)
	bytes := kvStore.Get(storeKey(assetAddress))
	if bytes == nil {
		return nil, assetNotFoundError(assetAddress.String())
	}
	err := baseMapper.codec.UnmarshalBinaryBare(bytes, &asset)
	if err != nil {
		panic(err)
	}
	return asset, nil
}
func (baseMapper baseMapper) Set(context sdkTypes.Context, asset types.Asset) sdkTypes.Error {
	bytes, err := baseMapper.codec.MarshalBinaryBare(asset)
	if err != nil {
		panic(err)
	}
	assetAddress := asset.GetAddress()
	kvStore := context.KVStore(baseMapper.storeKey)
	kvStore.Set(storeKey(assetAddress), bytes)
	return nil
}
func (baseMapper baseMapper) Delete(context sdkTypes.Context, assetAddress types.AssetAddress) {
	bytes, err := baseMapper.codec.MarshalBinaryBare(&baseAsset{})
	if err != nil {
		panic(err)
	}
	kvStore := context.KVStore(baseMapper.storeKey)
	kvStore.Set(storeKey(assetAddress), bytes)
}
