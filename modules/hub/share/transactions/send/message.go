package send

import (
	"github.com/asaskevich/govalidator"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/constants"
)

type Message struct {
	From    sdkTypes.AccAddress `json:"from" yaml:"from" valid:"required~from"`
	To      sdkTypes.AccAddress `json:"to" yaml:"to" valid:"required~to"`
	Address string              `json:"address" yaml:"address" valid:"required~address"`
}

var _ sdkTypes.Msg = Message{}

func (message Message) Route() string { return constants.ModuleName }
func (message Message) Type() string  { return constants.SendTransaction }
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
