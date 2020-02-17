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
	GetAsset(sdkTypes.Context, types.AssetAddress) (types.Asset, sdkTypes.Error)
	SetAsset(sdkTypes.Context, types.Asset) sdkTypes.Error
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

func (baseMapper baseMapper) GetAsset(context sdkTypes.Context, assetAddress types.AssetAddress) (asset types.Asset, error sdkTypes.Error) {
	kvStore := context.KVStore(baseMapper.storeKey)
	bytes := kvStore.Get(storeKey(assetAddress))
	if bytes == nil {
		return nil, AssetNotFoundError(assetAddress.String())
	}
	err := baseMapper.codec.UnmarshalBinaryBare(bytes, &asset)
	if err != nil {
		panic(err)
	}
	return asset, nil
}
func (baseMapper baseMapper) SetAsset(context sdkTypes.Context, asset types.Asset) sdkTypes.Error {
	bytes, err := baseMapper.codec.MarshalBinaryBare(asset)
	if err != nil {
		panic(err)
	}
	address := asset.GetAddress()
	kvStore := context.KVStore(baseMapper.storeKey)
	kvStore.Set(storeKey(address), bytes)
	return nil
}
