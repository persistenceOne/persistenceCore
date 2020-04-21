package sign

import (
	"github.com/asaskevich/govalidator"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/persistenceSDK/modules/contract/constants"
)

type Message struct {
	From sdkTypes.AccAddress `json:"from" yaml:"from" valid:"required~From"`
}

var _ sdkTypes.Msg = Message{}

func (message Message) Route() string { return constants.ModuleName }

func (message Message) Type() string { return constants.SignTransaction }

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
