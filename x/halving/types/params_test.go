/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package types

import (
	"fmt"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewParams(t *testing.T) {
	params := NewParams(100)
	require.Equal(t, Params{BlockHeight: uint64(100)}, params)
}

func TestDefaultParams(t *testing.T) {
	params := DefaultParams()
	require.Equal(t, Params{BlockHeight: uint64(2 * 60 * 60 * 8766 / 5)}, params)
}

func TestValidate(t *testing.T) {
	err := NewParams(100).Validate()
	require.Equal(t, nil, err)
}

func TestParams_String(t *testing.T) {
	s := NewParams(100).String()
	require.Equal(t, "blockHeight: 100\n", s)
}

func TestParams_ParamSetPairs(t *testing.T) {
	params := NewParams(100)
	got := params.ParamSetPairs()
	require.Equal(t, 1, len(got))

	expectedParamSetPair := paramsTypes.NewParamSetPair(KeyBlockHeight, &params.BlockHeight, validateBlockHeight)
	require.Equal(t, got[0].Key, expectedParamSetPair.Key)
	require.Equal(t, got[0].Value, expectedParamSetPair.Value)
}

func TestParams_validateBlockHeight(t *testing.T) {
	err := validateBlockHeight(uint64(100))
	require.Equal(t, nil, err)

	err = validateBlockHeight(-100)
	require.Equal(t, fmt.Errorf("invalid parameter type: int"), err)

}
