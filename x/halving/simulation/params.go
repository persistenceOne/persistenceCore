/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package simulation

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/persistenceOne/persistenceCore/x/halving/internal/types"
)

const (
	keyBlockHeight = "BlockHeight"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation
func ParamChanges(_ *rand.Rand) []simulation.ParamChange {
	return []simulation.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, keyBlockHeight,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%d\"", GetBlockHeight(r))
			},
		),
	}
}
