package types

type AssetAddress interface {
	Bytes() []byte
}

type Asset interface {
	GetAddress() AssetAddress
}
