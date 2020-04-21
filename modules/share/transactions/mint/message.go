package mint

import (
	"github.com/asaskevich/govalidator"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/persistenceSDK/modules/share/constants"
)

type Message struct {
	From    sdkTypes.AccAddress `json:"from" yaml:"from" valid:"required~from"`
	To      sdkTypes.AccAddress `json:"to" yaml:"to" valid:"required~to"`
	Address string              `json:"address" yaml:"address" valid:"required~address"`
	Lock    bool                `json:"lock" yaml:"lock"`
}

var _ sdkTypes.Msg = Message{}

func (message Message) Route() string { return constants.ModuleName }
func (message Message) Type() string  { return constants.MintTransaction }
func (message Message) ValidateBasic() error {
	var _, error = govalidator.ValidateStruct(message)
	if error != nil {
		return errors.Wrap(constants.IncorrectMessageCode, error.Error())
	}
	return nil
}
func (message Message) GetSignBytes() []byte {
	return sdkTypes.MustSortJSON(packageCodec.MustMarshalJSON(message))
}
func (message Message) GetSigners() []sdkTypes.AccAddress {
	return []sdkTypes.AccAddress{message.From}
}
