package application

import (
	"encoding/json"
	"fmt"
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/persistenceSDK/modules/assetFactory/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/types"
	"strings"
)

// this is for adding raw messages to wasm //

type CustomMsg struct {
	//	Debug string `json:"debug,omitempty"`
	MsgType string          `json:"msgtype,required"`
	Raw     json.RawMessage `json:"raw,omitempty"`
}

// Type will be assetFactory/mint , assetFactory/burn, assetFactory/Mmtate , like codec register types

func wasmCustomMessageEncoder(codec *codec.Codec) *wasm.MessageEncoders {

	return &wasm.MessageEncoders{
		Custom: customEncoder(codec),
	}
}

func customEncoder(codec *codec.Codec) wasm.CustomEncoder {
	return func(sender sdkTypes.AccAddress, msg json.RawMessage) ([]sdkTypes.Msg, error) {
		var customMessage CustomMsg
		err := json.Unmarshal(msg, &customMessage)
		if err != nil {
			return nil, sdkErrors.Wrap(sdkErrors.ErrJSONUnmarshal, err.Error())
		}

		fmt.Println("customMessage-MessageType: ", customMessage.MsgType, customMessage, string(msg))
		switch customMessage.MsgType {
		case "assetFactory/mint":
			return assetFactoryMintEncoder(codec, sender, customMessage.Raw)
		case "assetFactory/mutate":
			return assetFactoryMutateEncoder(codec, sender, customMessage.Raw)
		case "assetFactory/burn":
			return assetFactoryBurnEncoder(codec, sender, customMessage.Raw)
		default:
			return nil, sdkErrors.Wrap(wasm.ErrInvalidMsg, "Custom variant not supported in SDK- default")
		}
		return nil, sdkErrors.Wrap(wasm.ErrInvalidMsg, "Custom variant not supported in SDK")
	}
}

func assetFactoryMintEncoder(codec *codec.Codec, sender sdkTypes.AccAddress, rawMessage json.RawMessage) ([]sdkTypes.Msg, error) {
	if rawMessage != nil {
		var assetMessage AssetMintMessage
		fmt.Println(rawMessage)
		err := json.Unmarshal(rawMessage, &assetMessage)
		if err != nil {
			return nil, sdkErrors.Wrap(sdkErrors.ErrJSONUnmarshal, err.Error())
		}
		EncodeAssestmintMsg(sender, assetMessage)
	}
	return nil, sdkErrors.Wrap(wasm.ErrInvalidMsg, "Custom variant assetMint not supported")
}

func EncodeAssestmintMsg(sender sdkTypes.AccAddress, ast AssetMintMessage) ([]sdkTypes.Msg, error) {

	properties := strings.Split(ast.Properties, ",")
	basePropertyList := make([]types.BaseProperty, 0)
	for _, property := range properties {
		traitIDProperty := strings.Split(property, ":")
		if len(traitIDProperty) == 2 && traitIDProperty[0] != "" {
			basePropertyList = append(basePropertyList,
				types.BaseProperty{
					BaseID:   types.BaseID{IDString: traitIDProperty[0]},
					BaseFact: types.BaseFact{BaseString: traitIDProperty[1]},
				})
		}
	}

	chainid := types.BaseID{IDString: ast.ChainID}
	fromAddr, stderr := sdkTypes.AccAddressFromBech32(ast.From)
	if stderr != nil {
		return nil, sdkErrors.Wrap(sdkErrors.ErrInvalidAddress, ast.From)
	}
	maintainersID := types.BaseID{IDString: ast.MaintainersID}
	burn := types.BaseHeight{Height: ast.Burn}
	lock := types.BaseHeight{Height: ast.Lock}
	classificationID := types.BaseID{IDString: ast.ClassificationID}

	newmg := mint.Message{ChainID: chainid, From: fromAddr, Burn: burn, MaintainersID: maintainersID, Properties: &types.BaseProperties{BasePropertyList: basePropertyList}, ClassificationID: classificationID, Lock: lock}

	return []sdkTypes.Msg{newmg}, nil
}

func assetFactoryMutateEncoder(codec *codec.Codec, sender sdkTypes.AccAddress, rawMessage json.RawMessage) ([]sdkTypes.Msg, error) {
	return nil, sdkErrors.Wrap(wasm.ErrInvalidMsg, "Custom variant assetmutate not supported")
}

func assetFactoryBurnEncoder(codec *codec.Codec, sender sdkTypes.AccAddress, rawMessage json.RawMessage) ([]sdkTypes.Msg, error) {
	return nil, sdkErrors.Wrap(wasm.ErrInvalidMsg, "Custom variant assetburn not supported")
}

type AssetMintMessage struct {
	From             string `json:"from"`
	ChainID          string `json:"chainID"`
	MaintainersID    string `json:"maintainersID"`
	ClassificationID string `json:"classificationID"`
	Properties       string `json:"properties"`
	Lock             int    `json:"lock"`
	Burn             int    `json:"burn"`
}

type Property struct {
	ID             string     `json:"id"`
	baseFact       string     `json:"factBytes"`
	baseSignatures Signatures `json:"factSignatures"`
}

type Properties []Property

func (p Properties) MarshalJSON() ([]byte, error) {
	if len(p) == 0 {
		return []byte("[]"), nil
	}
	var d []Property = p
	return json.Marshal(d)
}

func (p *Properties) UnmarshalJSON(data []byte) error {
	// make sure we deserialize [] back to null
	if string(data) == "[]" || string(data) == "null" {
		return nil
	}
	var d []Property
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	*p = d
	return nil
}

type Signature struct {
	SignatureID    string `json:"signatureID"`
	SignatureBytes []byte `json:"signatureBytes"`
	ValidityHeight int    `json:"validityHeight"`
}

type Signatures []Signature

func (s Signatures) MarshalJSON() ([]byte, error) {
	if len(s) == 0 {
		return []byte("[]"), nil
	}
	var d []Signature = s
	return json.Marshal(d)
}

func (s *Signatures) UnmarshalJSON(data []byte) error {
	// make sure we deserialize [] back to null
	if string(data) == "[]" || string(data) == "null" {
		return nil
	}
	var d []Signature
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}
	*s = d
	return nil
}

//till here
