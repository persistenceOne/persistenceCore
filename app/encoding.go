/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package app

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	grouptypes "github.com/cosmos/cosmos-sdk/x/group"
	interchainquerytypes "github.com/persistenceOne/persistence-sdk/v7/x/interchainquery/types"
	liquidstakeibctypes "github.com/persistenceOne/persistence-sdk/v7/x/liquidstakeibc/types"
	lscosmostypes "github.com/persistenceOne/persistence-sdk/v7/x/lscosmos/types"
	"github.com/persistenceOne/persistence-sdk/v7/x/lsm/distribution"
	"github.com/persistenceOne/persistence-sdk/v7/x/lsm/staking"
	oracletypes "github.com/persistenceOne/persistence-sdk/v7/x/oracle/types"
	pobtypes "github.com/persistenceOne/persistence-sdk/v7/x/pob/types"
	ratesynctypes "github.com/persistenceOne/persistence-sdk/v7/x/ratesync/types"

	"github.com/persistenceOne/persistenceCore/v16/app/params"
)

// MakeEncodingConfig creates an EncodingConfig for testing
func MakeEncodingConfig() params.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()

	std.RegisterLegacyAminoCodec(encodingConfig.Amino)
	std.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return encodingConfig
}

func AppendModuleCodecs(amino *codec.LegacyAmino, interfaceRegistry codectypes.InterfaceRegistry, mbm module.BasicManager) {
	if mbm == nil {
		panic("no modules provided")
	}

	AppendModuleInterfaces(interfaceRegistry, mbm)
	AppendModuleLegacyCodecs(amino, mbm)
}

func AppendModuleInterfaces(interfaceRegistry codectypes.InterfaceRegistry, mbm module.BasicManager) {
	if mbm == nil {
		panic("no modules provided")
	}

	mbm.RegisterInterfaces(interfaceRegistry)

	//deprecated modules types
	lscosmostypes.RegisterInterfaces(interfaceRegistry)

	liquidstakeibctypes.RegisterInterfaces(interfaceRegistry)

	ratesynctypes.RegisterInterfaces(interfaceRegistry)

	interchainquerytypes.RegisterInterfaces(interfaceRegistry)
	// oracle, but was never used
	oracletypes.RegisterInterfaces(interfaceRegistry)

	// group module,
	grouptypes.RegisterInterfaces(interfaceRegistry)
	//ibcfee, but was never used ...

	// cosmos-sdk-lsm staking msgs
	staking.RegisterInterfaces(interfaceRegistry)

	// cosmos-sdk-lsm distribution msgs
	distribution.RegisterInterfaces(interfaceRegistry)

	// skip-pob - testnet
	pobtypes.RegisterInterfaces(interfaceRegistry)
}

func AppendModuleLegacyCodecs(amino *codec.LegacyAmino, mbm module.BasicManager) {
	if mbm == nil {
		panic("no modules provided")
	}

	mbm.RegisterLegacyAminoCodec(amino)

	//deprecated modules types
	lscosmostypes.RegisterLegacyAminoCodec(amino)

	liquidstakeibctypes.RegisterLegacyAminoCodec(amino)

	ratesynctypes.RegisterCodec(amino)

	interchainquerytypes.RegisterLegacyAminoCodec(amino)
	// oracle, but was never used
	oracletypes.RegisterLegacyAminoCodec(amino)

	// group module,
	grouptypes.RegisterLegacyAminoCodec(amino)
	//ibcfee, but was never used ...

	// cosmos-sdk-lsm staking msgs
	staking.RegisterLegacyAminoCodec(amino)

	// cosmos-sdk-lsm distribution msgs
	distribution.RegisterLegacyAminoCodec(amino)

	// skip-pob - testnet
	pobtypes.RegisterLegacyAminoCodec(amino)
}
