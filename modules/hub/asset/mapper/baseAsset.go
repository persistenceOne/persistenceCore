package mapper

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

type baseAsset struct {
	assetAddress types.AssetAddress
	owner        sdkTypes.AccAddress
}

func NewAsset(address string, owner sdkTypes.AccAddress) types.Asset {
	return baseAsset{
		assetAddress: newAssetAddress(address),
		owner:        owner,
	}
}

var _ types.Asset = (*baseAsset)(nil)

func (baseAsset baseAsset) GetAddress() types.AssetAddress { return baseAsset.assetAddress }
func (baseAsset baseAsset) GetOwner() sdkTypes.AccAddress  { return baseAsset.owner }
