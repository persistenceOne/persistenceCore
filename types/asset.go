package types

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type AssetAddress interface {
	Bytes() []byte
}

type Asset interface {
	GetAddress() AssetAddress
	GetOwner() sdkTypes.AccAddress
}
