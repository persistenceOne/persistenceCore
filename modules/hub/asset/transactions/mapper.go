package mapper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

var assetStoreKeyPrefix = []byte{0x11}

func assetStoreKey(assetAddress types.AssetAddress) []byte {
	return append(assetStoreKeyPrefix, assetAddress.Bytes()...)
}

type Mapper interface {
	GetAsset(sdkTypes.Context, types.AssetAddress) (types.Asset, error)
	SetAsset(types.Asset) error
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

func (baseMapper baseMapper) GetAsset(context sdkTypes.Context, assetAddress types.AssetAddress) (asset types.Asset, error error) {
	kvStore := context.KVStore(baseMapper.storeKey)
	bz := kvStore.Get(assetStoreKey(assetAddress))
	if bz == nil {
		return nil, nil
	}
	err := baseMapper.codec.UnmarshalBinaryBare(bz, &asset)
	if err != nil {
		return nil, nil
	}
	return asset, nil
}

func (baseMapper baseMapper) SetAsset(types.Asset) error { return nil }
