/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package app

import (
	"github.com/cosmos/cosmos-sdk/std"
	grouptypes "github.com/cosmos/cosmos-sdk/x/group"
	interchainquerytypes "github.com/persistenceOne/persistence-sdk/v5/x/interchainquery/types"
	"github.com/persistenceOne/persistence-sdk/v5/x/lsm/distribution"
	"github.com/persistenceOne/persistence-sdk/v5/x/lsm/staking"
	oracletypes "github.com/persistenceOne/persistence-sdk/v5/x/oracle/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v5/x/liquidstakeibc/types"
	lscosmostypes "github.com/persistenceOne/pstake-native/v5/x/lscosmos/types"
	ratesynctypes "github.com/persistenceOne/pstake-native/v5/x/ratesync/types"

	"github.com/persistenceOne/persistenceCore/v15/app/params"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()

	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	ModuleBasics.RegisterLegacyAminoCodec(encodingConfig.Amino)
	ModuleBasics.RegisterInterfaces(encodingConfig.InterfaceRegistry)

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

	// group module,
	grouptypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	grouptypes.RegisterLegacyAminoCodec(encodingConfig.Amino)
	//ibcfee, but was never used ...

	// cosmos-sdk-lsm staking msgs
	staking.RegisterLegacyAminoCodec(encodingConfig.Amino)
	staking.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	// cosmos-sdk-lsm distribution msgs
	distribution.RegisterLegacyAminoCodec(encodingConfig.Amino)
	distribution.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
