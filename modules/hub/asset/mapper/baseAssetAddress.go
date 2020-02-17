package mapper

import (
	"github.com/persistenceOne/persistenceSDK/types"
)

type baseAssetAddress struct {
	Address string
}

func NewAssetAddress(address string) types.AssetAddress {
	return baseAssetAddress{
		Address: address,
	}
}

var _ types.AssetAddress = (*baseAssetAddress)(nil)

func (baseAssetAddress baseAssetAddress) Bytes() []byte  { return []byte(baseAssetAddress.Address) }
func (baseAssetAddress baseAssetAddress) String() string { return baseAssetAddress.Address }
