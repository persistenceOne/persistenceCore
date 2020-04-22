package mapper

import (
	"encoding/json"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

type baseAsset struct {
	Address types.AssetAddress  `json:"address" yaml:"address" valid:"required~address"`
	Owner   sdkTypes.AccAddress `json:"owner" yaml:"owner" valid:"required~owner"`
	Lock    bool                `json:"lock" yaml:"lock"`
}

func newAsset(address types.AssetAddress, owner sdkTypes.AccAddress, lock bool) types.Asset {
	return &baseAsset{
		Address: address,
		Owner:   owner,
		Lock:    lock,
	}
}

var _ types.Asset = (*baseAsset)(nil)

func (baseAsset baseAsset) GetAddress() types.AssetAddress      { return baseAsset.Address }
func (baseAsset baseAsset) GetOwner() sdkTypes.AccAddress       { return baseAsset.Owner }
func (baseAsset *baseAsset) SetOwner(Owner sdkTypes.AccAddress) { baseAsset.Owner = Owner }
func (baseAsset baseAsset) GetLock() bool                       { return baseAsset.Lock }
func (baseAsset *baseAsset) SetLock(Lock bool)                  { baseAsset.Lock = Lock }
func (baseAsset baseAsset) String() string {
	bytes, error := json.Marshal(baseAsset)
	if error != nil {
		panic(error)
	}
	return string(bytes)
}
