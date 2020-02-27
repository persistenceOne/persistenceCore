package mapper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

var storeKeyPrefix = []byte{0x11}

func storeKey(shareAddress types.ShareAddress) []byte {
	return append(storeKeyPrefix, shareAddress.Bytes()...)
}

type Mapper interface {
	Create(sdkTypes.Context, types.ShareAddress, sdkTypes.AccAddress, bool) sdkTypes.Error
	Read(sdkTypes.Context, types.ShareAddress) (types.Share, sdkTypes.Error)
	Update(sdkTypes.Context, types.Share) sdkTypes.Error
	Delete(sdkTypes.Context, types.ShareAddress)
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

func (baseMapper baseMapper) Create(context sdkTypes.Context, address types.ShareAddress, owner sdkTypes.AccAddress, lock bool) sdkTypes.Error {
	share := newShare(address, owner, lock)
	bytes, err := baseMapper.codec.MarshalBinaryBare(share)
	if err != nil {
		panic(err)
	}
	shareAddress := share.GetAddress()
	kvStore := context.KVStore(baseMapper.storeKey)
	kvStore.Set(storeKey(shareAddress), bytes)
	return nil
}
func (baseMapper baseMapper) Read(context sdkTypes.Context, address types.ShareAddress) (share types.Share, error sdkTypes.Error) {
	kvStore := context.KVStore(baseMapper.storeKey)
	bytes := kvStore.Get(storeKey(address))
	if bytes == nil {
		return nil, shareNotFoundError(address.String())
	}
	err := baseMapper.codec.UnmarshalBinaryBare(bytes, &share)
	if err != nil {
		panic(err)
	}
	return share, nil
}
func (baseMapper baseMapper) Update(context sdkTypes.Context, share types.Share) sdkTypes.Error {
	bytes, err := baseMapper.codec.MarshalBinaryBare(share)
	if err != nil {
		panic(err)
	}
	shareAddress := share.GetAddress()
	kvStore := context.KVStore(baseMapper.storeKey)
	kvStore.Set(storeKey(shareAddress), bytes)
	return nil
}
func (baseMapper baseMapper) Delete(context sdkTypes.Context, address types.ShareAddress) {
	bytes, err := baseMapper.codec.MarshalBinaryBare(&baseShare{})
	if err != nil {
		panic(err)
	}
	kvStore := context.KVStore(baseMapper.storeKey)
	kvStore.Set(storeKey(address), bytes)
}
