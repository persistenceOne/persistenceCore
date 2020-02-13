package types

type AssetAddress interface {
	Bytes() []byte
	String() string
}
