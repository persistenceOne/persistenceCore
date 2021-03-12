/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

// Parameter store keys
var (
	KeyBlockHeight = []byte("BlockHeight")
	Factor         = sdk.NewInt(2)
)

// ParamTable for halving module.
func ParamKeyTable() paramsTypes.KeyTable {
	return paramsTypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(blockHeight uint64) Params {
	return Params{
		BlockHeight: blockHeight,
	}
}

// default halving module parameters
func DefaultParams() Params {
	return Params{
		BlockHeight: uint64(2 * 60 * 60 * 8766 / 5), // 2 * blocksPerYear assuming 5s per block
	}
}

// validate params
func (p Params) Validate() error {
	if err := validateBlockHeight(p.BlockHeight); err != nil {
		return err
	}
	return nil
}

//stringer function
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Implements params.ParamSet
func (p *Params) ParamSetPairs() paramsTypes.ParamSetPairs {
	return paramsTypes.ParamSetPairs{
		paramsTypes.NewParamSetPair(KeyBlockHeight, &p.BlockHeight, validateBlockHeight),
	}
}

func validateBlockHeight(i interface{}) error {
	_, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
