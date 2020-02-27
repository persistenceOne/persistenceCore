package types

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

type Share interface {
	//Immutable
	GetAddress() ShareAddress
	//Mutable
	GetOwner() sdkTypes.AccAddress
	SetOwner(sdkTypes.AccAddress)
	GetLock() bool
	SetLock(bool)
	String() string
}
