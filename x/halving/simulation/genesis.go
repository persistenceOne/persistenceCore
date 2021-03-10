/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/persistenceOne/persistenceCore/x/halving/types"
)

// Simulation parameter constants
const (
	BlockHeight = "blockHeight"
)

// GetBlockHeight randomized BlockHeight
func GetBlockHeight(r *rand.Rand) uint64 {
	return uint64(r.Intn(1000))
}

// RandomizedGenState generates a random GenesisState for halving
func RandomizedGenState(simState *module.SimulationState) {

	// params
	var blockHeight uint64
	simState.AppParams.GetOrGenerate(
		simState.Cdc, BlockHeight, &blockHeight, simState.Rand,
		func(r *rand.Rand) { blockHeight = GetBlockHeight(r) },
	)

	halvingGenesis := types.NewGenesisState(types.NewParams(blockHeight))

	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", codec.MustMarshalJSONIndent(simState.Cdc, halvingGenesis))
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(halvingGenesis)
}
