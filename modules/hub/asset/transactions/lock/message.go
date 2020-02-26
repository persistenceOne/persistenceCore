package lock

import (
	"github.com/asaskevich/govalidator"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/constants"
)

type Message struct {
	From    sdkTypes.AccAddress `json:"from" yaml:"from" valid:"required~from"`
	Address string              `json:"address" yaml:"address" valid:"required~address"`
	Lock    bool                `json:"lock" yaml:"lock"`
}

var _ sdkTypes.Msg = Message{}

func (message Message) Route() string { return constants.ModuleName }
func (message Message) Type() string  { return constants.LockTransaction }
func (message Message) ValidateBasic() sdkTypes.Error {
	var _, error = govalidator.ValidateStruct(message)
	if error != nil {
		return incorrectMessageError(error.Error())
	}
	return nil
}
func (message Message) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(packageCodec.MustMarshalJSON(message))
}
func (message Message) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{message.From}
}
