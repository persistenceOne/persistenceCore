/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

const moduleName = "queuing"

func RegisterCodec(codec *codec.LegacyAmino) {
	codec.RegisterConcrete(KafkaMsg{}, moduleName+"KafkaMsg", nil)
}

var ModuleCdc *codec.LegacyAmino

func init() {
	ModuleCdc = codec.NewLegacyAmino()
	RegisterCodec(ModuleCdc)
	ModuleCdc.Seal()
}
