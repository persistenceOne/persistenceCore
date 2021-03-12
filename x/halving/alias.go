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
)

var (
	// functions aliases

	NewKeeper       = keeper.NewKeeper
	NewGenesisState = types.NewGenesisState

	// variable aliases

	Factor = types.Factor
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params
)
