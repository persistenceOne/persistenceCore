package types

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

type ShareAddress interface {
	Bytes() []byte
}

type Share interface {
	GetAddress() ShareAddress
	GetOwner() sdkTypes.AccAddress
}
