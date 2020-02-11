package types

type AssetAddress interface {
	Bytes() []byte
	String() string
}
type baseAssetAddress struct {
	address string
}

func newAssetAddress(address string) AssetAddress {
	return baseAssetAddress{
		address: address,
	}
}

var _ AssetAddress = (*baseAssetAddress)(nil)

func (baseAssetAddress baseAssetAddress) Bytes() []byte  { return []byte(baseAssetAddress.address) }
func (baseAssetAddress baseAssetAddress) String() string { return baseAssetAddress.address }
