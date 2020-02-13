package mapper

import (
	"github.com/persistenceOne/persistenceSDK/types"
)

type baseAssetAddress struct {
	address string
}

func newAssetAddress(address string) types.AssetAddress {
	return baseAssetAddress{
		address: address,
	}
}

var _ types.AssetAddress = (*baseAssetAddress)(nil)

func (baseAssetAddress baseAssetAddress) Bytes() []byte  { return []byte(baseAssetAddress.address) }
func (baseAssetAddress baseAssetAddress) String() string { return baseAssetAddress.address }
