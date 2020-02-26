package types

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type Asset interface {
	//Immutable
	GetAddress() AssetAddress
	//Mutable
	GetOwner() sdkTypes.AccAddress
	SetOwner(sdkTypes.AccAddress)
	GetLock() bool
	SetLock(bool)
	String() string
}
