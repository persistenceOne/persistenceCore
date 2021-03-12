/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package simulation

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/cosmos/cosmos-sdk/types/kv"

	"github.com/persistenceOne/persistenceCore/x/halving/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding halving type
func DecodeStore(cdc codec.Marshaler) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key, types.HalvingKey):
			return fmt.Sprintf("%v\n%v", kvA, kvB)
		default:
			panic(fmt.Sprintf("invalid halving key %X", kvA.Key))
		}
	}
}
