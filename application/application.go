/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package application

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmClient "github.com/CosmWasm/wasmd/x/wasm/client"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsClient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeClient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	"github.com/persistenceOne/persistenceSDK/modules/assets"
	"github.com/persistenceOne/persistenceSDK/modules/classifications"
	"github.com/persistenceOne/persistenceSDK/modules/identities"
	"github.com/persistenceOne/persistenceSDK/modules/maintainers"
	"github.com/persistenceOne/persistenceSDK/modules/metas"
	"github.com/persistenceOne/persistenceSDK/modules/orders"
	"github.com/persistenceOne/persistenceSDK/modules/splits"
	"github.com/persistenceOne/persistenceSDK/schema/applications/base"
)

var moduleAccountPermissions = map[string][]string{
	auth.FeeCollectorName:     nil,
	distribution.ModuleName:   nil,
	mint.ModuleName:           {supply.Minter},
	staking.BondedPoolName:    {supply.Burner, supply.Staking},
	staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	gov.ModuleName:            {supply.Burner},
	splits.Prototype().Name(): nil,
}
var tokenReceiveAllowedModules = map[string]bool{
	distribution.ModuleName: true,
}
var ModuleBasics = module.NewBasicManager(
	genutil.AppModuleBasic{},
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.NewAppModuleBasic(append(wasmClient.ProposalHandlers, paramsClient.ProposalHandler, distribution.ProposalHandler, upgradeClient.ProposalHandler)...),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	wasm.AppModuleBasic{},
	slashing.AppModuleBasic{},
	supply.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},

	assets.Prototype(),
	classifications.Prototype(),
	identities.Prototype(),
	maintainers.Prototype(),
	metas.Prototype(),
	orders.Prototype(),
	splits.Prototype(),
)
var NewApplication = base.Prototype(
	Name,
	Codec,
	wasm.EnableAllProposals,
	moduleAccountPermissions,
	tokenReceiveAllowedModules,
)
