/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package halving

import (
	"github.com/persistenceOne/persistenceCore/x/halving/keeper"
	"github.com/persistenceOne/persistenceCore/x/halving/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
)

var (
	// functions aliases

	NewKeeper           = keeper.NewKeeper
	NewQuerier          = keeper.NewQuerier
	NewGenesisState     = types.NewGenesisState
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis

	// variable aliases

	ModuleCdc = types.ModuleCdc
	Factor    = types.Factor
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params
)
