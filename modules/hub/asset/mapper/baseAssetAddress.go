package mapper

import (
	"encoding/json"
	"github.com/persistenceOne/persistenceSDK/types"
)

type baseAssetAddress struct {
	Address string `json:"address" yaml:"address" valid:"required~address"`
}

func NewAssetAddress(address string) types.AssetAddress {
	return baseAssetAddress{
		Address: address,
	}
}

var _ types.AssetAddress = (*baseAssetAddress)(nil)

func (baseAssetAddress baseAssetAddress) Bytes() []byte { return []byte(baseAssetAddress.Address) }
func (baseAssetAddress baseAssetAddress) String() string {
	bytes, error := json.Marshal(baseAssetAddress)
	if error != nil {
		panic(error)
	}
	return string(bytes)
}
