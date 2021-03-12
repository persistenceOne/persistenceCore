/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package types

const (
	// ModuleName
	ModuleName = "halving"

	// DefaultParamspace params keeper
	DefaultParamspace = ModuleName

	// StoreKey is the default store key for halving
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the halving store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the halving querier
	QueryParameters = "parameters"
)
