package mapper

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

type baseAsset struct {
	AssetAddress types.AssetAddress
	Owner        sdkTypes.AccAddress
}

func NewAsset(address string, owner sdkTypes.AccAddress) types.Asset {
	return baseAsset{
		AssetAddress: NewAssetAddress(address),
		Owner:        owner,
	}
}

var _ types.Asset = (*baseAsset)(nil)

func (baseAsset baseAsset) GetAddress() types.AssetAddress { return baseAsset.AssetAddress }
func (baseAsset baseAsset) GetOwner() sdkTypes.AccAddress  { return baseAsset.Owner }
