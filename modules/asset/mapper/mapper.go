package mapper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/persistenceSDK/modules/asset/constants"
	"github.com/persistenceOne/persistenceSDK/types"
)

var storeKeyPrefix = []byte{0x11}

func storeKey(assetAddress types.AssetAddress) []byte {
	return append(storeKeyPrefix, assetAddress.Bytes()...)
}

type Mapper interface {
	Create(sdkTypes.Context, types.AssetAddress, sdkTypes.AccAddress, bool) error
	Read(sdkTypes.Context, types.AssetAddress) (types.Asset, error)
	Update(sdkTypes.Context, types.Asset) error
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

func (baseMapper baseMapper) Create(context sdkTypes.Context, address types.AssetAddress, owner sdkTypes.AccAddress, lock bool) error {
	asset := newAsset(address, owner, lock)
	bytes, err := baseMapper.codec.MarshalBinaryBare(asset)
	if err != nil {
		panic(err)
	}
	assetAddress := asset.GetAddress()
	kvStore := context.KVStore(baseMapper.storeKey)
	kvStore.Set(storeKey(assetAddress), bytes)
	return nil
}
func (baseMapper baseMapper) Read(context sdkTypes.Context, address types.AssetAddress) (asset types.Asset, error error) {
	kvStore := context.KVStore(baseMapper.storeKey)
	bytes := kvStore.Get(storeKey(address))
	if bytes == nil {
		return nil, errors.Wrap(constants.AssetNotFoundCode, address.String())
	}
	err := baseMapper.codec.UnmarshalBinaryBare(bytes, &asset)
	if err != nil {
		panic(err)
	}
	return asset, nil
}
func (baseMapper baseMapper) Update(context sdkTypes.Context, asset types.Asset) error {
	bytes, err := baseMapper.codec.MarshalBinaryBare(asset)
	if err != nil {
		panic(err)
	}
	assetAddress := asset.GetAddress()
	kvStore := context.KVStore(baseMapper.storeKey)
	kvStore.Set(storeKey(assetAddress), bytes)
	return nil
}
func (baseMapper baseMapper) Delete(context sdkTypes.Context, address types.AssetAddress) {
	bytes, err := baseMapper.codec.MarshalBinaryBare(&baseAsset{})
	if err != nil {
		panic(err)
	}
	kvStore := context.KVStore(baseMapper.storeKey)
	kvStore.Set(storeKey(address), bytes)
}
