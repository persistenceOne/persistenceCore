package types

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type Asset interface {
	GetAddress() AssetAddress
	GetOwner() sdkTypes.AccAddress
}
type baseAsset struct {
	assetAddress AssetAddress
	owner        sdkTypes.AccAddress
}

func NewAsset(address string, owner sdkTypes.AccAddress) Asset {
	return baseAsset{
		assetAddress: newAssetAddress(address),
		owner:        owner,
	}
}

var _ Asset = (*baseAsset)(nil)

func (baseAsset baseAsset) GetAddress() AssetAddress      { return baseAsset.assetAddress }
func (baseAsset baseAsset) GetOwner() sdkTypes.AccAddress { return baseAsset.owner }
