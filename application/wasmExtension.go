package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/persistenceSDK/modules/assets/constants"
	"github.com/persistenceOne/persistenceSDK/modules/assets/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/types"
	"strings"
)

// this is for adding raw messages to wasm //

type customMessage struct {
	MsgType string          `json:"msgtype,required"`
	Raw     json.RawMessage `json:"raw,omitempty"`
}

// Type will be assets/mint , assets/burn, assets/Mmtate , (like codec register types)

func WasmCustomMessageEncoder(codec *codec.Codec) *wasm.MessageEncoders {
	return &wasm.MessageEncoders{
		Custom: customEncoder(codec),
	}
}

func customEncoder(codec *codec.Codec) wasm.CustomEncoder {
	return func(sender sdkTypes.AccAddress, rawMessage json.RawMessage) ([]sdkTypes.Msg, error) {
		var customMessage customMessage
		err := json.Unmarshal(rawMessage, &customMessage)
		if err != nil {
			return nil, sdkErrors.Wrap(sdkErrors.ErrJSONUnmarshal, err.Error())
		}

		switch customMessage.MsgType {
		case "assets/mint":
			return assetsMintEncoder(codec, sender, customMessage.Raw)
		case "assets/mutate":
			return assetsMutateEncoder(codec, sender, customMessage.Raw)
		case "assets/burn":
			return assetsBurnEncoder(codec, sender, customMessage.Raw)
		default:
			return nil, sdkErrors.Wrap(wasm.ErrInvalidMsg, "Custom variant not supported in SDK")
		}
	}
}

func assetsMintEncoder(codec *codec.Codec, sender sdkTypes.AccAddress, rawMessage json.RawMessage) ([]sdkTypes.Msg, error) {
	if rawMessage != nil {
		var assetMintMessage AssetMintMessage
		err := json.Unmarshal(rawMessage, &assetMintMessage)
		if err != nil {
			return nil, sdkErrors.Wrap(sdkErrors.ErrJSONUnmarshal, err.Error())
		}
		return encodeAssetMintMessage(sender, assetMintMessage)
	}
	return nil, sdkErrors.Wrap(wasm.ErrInvalidMsg, "rawMessage cannot be nil or empty")
}

func encodeAssetMintMessage(sender sdkTypes.AccAddress, assetMintMessage AssetMintMessage) ([]sdkTypes.Msg, error) {

	properties := strings.Split(assetMintMessage.Properties, constants.PropertiesSeparator)
	if len(properties) > constants.MaxTraitCount {
		panic(errors.New(fmt.Sprintf("")))
	}

	var propertyList []types.Property
	for _, property := range properties {
		traitIDAndProperty := strings.Split(property, constants.TraitIDAndPropertySeparator)
		if len(traitIDAndProperty) == 2 && traitIDAndProperty[0] != "" {
			propertyList = append(propertyList, types.NewProperty(types.NewID(traitIDAndProperty[0]), types.NewFact(traitIDAndProperty[1], types.NewSignatures(nil))))
		}
	}

	mintMessage := mint.Message{
		ChainID:          types.NewID(assetMintMessage.ChainID),
		From:             sender,
		Burn:             types.NewHeight(assetMintMessage.Burn),
		MaintainersID:    types.NewID(assetMintMessage.MaintainersID),
		Properties:       types.NewProperties(propertyList),
		ClassificationID: types.NewID(assetMintMessage.ClassificationID),
		Lock:             types.NewHeight(assetMintMessage.Lock),
	}
	return []sdkTypes.Msg{mintMessage}, nil
}

func assetsMutateEncoder(codec *codec.Codec, sender sdkTypes.AccAddress, rawMessage json.RawMessage) ([]sdkTypes.Msg, error) {
	return nil, sdkErrors.Wrap(wasm.ErrInvalidMsg, "Custom variant assetMutate not supported")
}

func assetsBurnEncoder(codec *codec.Codec, sender sdkTypes.AccAddress, rawMessage json.RawMessage) ([]sdkTypes.Msg, error) {
	return nil, sdkErrors.Wrap(wasm.ErrInvalidMsg, "Custom variant assetBurn not supported")
}

// AssetMintMessage should look like rest request, or similar and should be convertible to sdk message
type AssetMintMessage struct {
	From             string `json:"from"`
	ChainID          string `json:"chainID"`
	MaintainersID    string `json:"maintainersID"`
	ClassificationID string `json:"classificationID"`
	Properties       string `json:"properties"`
	Lock             int64  `json:"lock"`
	Burn             int64  `json:"burn"`
}
