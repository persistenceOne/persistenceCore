/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package simulation

import (
	"bytes"
	"fmt"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceCore/x/halving/internal/types"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding halving type
func DecodeStore(cdc *codec.Codec, kvA, kvB tmkv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key, types.HalvingKey):
		var paramA, paramB types.Params
		cdc.MustUnmarshalBinaryLengthPrefixed(kvA.Value, &paramA)
		cdc.MustUnmarshalBinaryLengthPrefixed(kvB.Value, &paramB)
		return fmt.Sprintf("%v\n%v", paramA, paramB)
	default:
		panic(fmt.Sprintf("invalid halving key %X", kvA.Key))
	}
}
