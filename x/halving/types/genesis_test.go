/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package types

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewGenesisState(t *testing.T) {
	params := NewParams(uint64(100))
	genesisState := NewGenesisState(params)
	require.Equal(t, &GenesisState{Params: params}, genesisState)
}

func TestDefaultGenesisState(t *testing.T) {
	params := NewParams(uint64(2 * 60 * 60 * 8766 / 5))
	require.Equal(t, &GenesisState{Params: params}, DefaultGenesisState())
}

func TestValidateGenesis(t *testing.T) {
	params := NewParams(uint64(100))
	genesisState := NewGenesisState(params)
	err := ValidateGenesis(*genesisState)
	require.Equal(t, nil, err)
}
