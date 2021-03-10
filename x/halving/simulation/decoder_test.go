/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package simulation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceCore/x/halving/types"
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	return
}

func TestDecodeStore(t *testing.T) {
	cdc := makeTestCodec()
	params := types.NewParams(uint64(10000000))

	kvPairs := tmkv.Pairs{
		tmkv.Pair{Key: types.HalvingKey, Value: cdc.MustMarshalBinaryLengthPrefixed(params)},
		tmkv.Pair{Key: []byte{0x99}, Value: []byte{0x99}},
	}
	tests := []struct {
		name        string
		expectedLog string
	}{
		{"BlockHeight", fmt.Sprintf("%v\n%v", params, params)},
		{"other", ""},
	}

	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { DecodeStore(cdc, kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
