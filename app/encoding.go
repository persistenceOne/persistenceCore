/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package app

import (
	"github.com/cosmos/cosmos-sdk/std"
	sdkdistr "github.com/cosmos/cosmos-sdk/x/distribution"
	sdkslashing "github.com/cosmos/cosmos-sdk/x/slashing"
	sdkstaking "github.com/cosmos/cosmos-sdk/x/staking"

	"github.com/persistenceOne/persistenceCore/v12/app/params"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()
	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	sdkstaking.AppModuleBasic{}.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	sdkslashing.AppModuleBasic{}.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	sdkdistr.AppModuleBasic{}.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}
