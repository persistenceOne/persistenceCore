package types

import sdkTypes "github.com/cosmos/cosmos-sdk/types"

type Share interface {
	GetAddress() ShareAddress
	GetOwner() sdkTypes.AccAddress
}
