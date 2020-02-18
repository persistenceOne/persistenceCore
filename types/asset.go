package types

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type Asset interface {
	GetAddress() AssetAddress
	GetOwner() sdkTypes.AccAddress
	String() string
}
