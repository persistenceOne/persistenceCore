/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package types

import "cosmossdk.io/collections"

const (
	// ModuleName
	ModuleName = "halving"

	// DefaultParamspace params keeper
	DefaultParamspace = ModuleName

	// StoreKey is the default store key for halving
	StoreKey = ModuleName

	// RouterKey is the message route for halving
	RouterKey = ModuleName
)

var ParamsKey = collections.NewPrefix(0)
