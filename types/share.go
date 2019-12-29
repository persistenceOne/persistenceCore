package types

type ShareAddress interface {
	Bytes() []byte
}

type Share interface {
	GetAddress() ShareAddress
}
