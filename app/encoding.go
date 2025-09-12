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
	interchainquerytypes "github.com/persistenceOne/persistence-sdk/v4/x/interchainquery/types"
	oracletypes "github.com/persistenceOne/persistence-sdk/v4/x/oracle/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v4/x/liquidstakeibc/types"
	lscosmostypes "github.com/persistenceOne/pstake-native/v4/x/lscosmos/types"
	ratesynctypes "github.com/persistenceOne/pstake-native/v4/x/ratesync/types"

	"github.com/persistenceOne/persistenceCore/v14/app/params"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()

	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	//for icacontroller callback applications
	sdkstaking.AppModuleBasic{}.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	sdkslashing.AppModuleBasic{}.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	sdkdistr.AppModuleBasic{}.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	//deprecated modules types
	lscosmostypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	lscosmostypes.RegisterLegacyAminoCodec(encodingConfig.Amino)

	liquidstakeibctypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	liquidstakeibctypes.RegisterLegacyAminoCodec(encodingConfig.Amino)

	ratesynctypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ratesynctypes.RegisterCodec(encodingConfig.Amino)

	interchainquerytypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	interchainquerytypes.RegisterLegacyAminoCodec(encodingConfig.Amino)
	// oracle, but was never used
	oracletypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	oracletypes.RegisterLegacyAminoCodec(encodingConfig.Amino)
	//ibcfee, but was never used ...
	return encodingConfig
}
