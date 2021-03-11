/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"

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

	bz, err := json.MarshalIndent(&halvingGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated halving parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(halvingGenesis)
}
