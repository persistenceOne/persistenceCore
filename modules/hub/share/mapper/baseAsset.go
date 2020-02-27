package mapper

import (
	"encoding/json"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/types"
)

type baseShare struct {
	Address types.ShareAddress  `json:"address" yaml:"address" valid:"required~address"`
	Owner   sdkTypes.AccAddress `json:"owner" yaml:"owner" valid:"required~owner"`
	Lock    bool                `json:"lock" yaml:"lock"`
}

func newShare(address types.ShareAddress, owner sdkTypes.AccAddress, lock bool) types.Share {
	return &baseShare{
		Address: address,
		Owner:   owner,
		Lock:    lock,
	}
}

var _ types.Share = (*baseShare)(nil)

func (baseShare baseShare) GetAddress() types.ShareAddress      { return baseShare.Address }
func (baseShare baseShare) GetOwner() sdkTypes.AccAddress       { return baseShare.Owner }
func (baseShare *baseShare) SetOwner(Owner sdkTypes.AccAddress) { baseShare.Owner = Owner }
func (baseShare baseShare) GetLock() bool                       { return baseShare.Lock }
func (baseShare *baseShare) SetLock(Lock bool)                  { baseShare.Lock = Lock }
func (baseShare baseShare) String() string {
	bytes, error := json.Marshal(baseShare)
	if error != nil {
		panic(error)
	}
	return string(bytes)
}
