package mapper

import (
	"encoding/json"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

type baseAsset struct {
	AssetAddress types.AssetAddress  `json:"assetAddress" yaml:"assetAddress" valid:"required~assetAddress"`
	Owner        sdkTypes.AccAddress `json:"owner" yaml:"owner" valid:"required~owner"`
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
func (baseAsset baseAsset) String() string {
	bytes, error := json.Marshal(baseAsset)
	if error != nil {
		panic(error)
	}
	return string(bytes)
}
