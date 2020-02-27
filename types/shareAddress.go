package types

type ShareAddress interface {
	Bytes() []byte
	String() string
}
