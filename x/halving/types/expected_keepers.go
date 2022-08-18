/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

type MintKeeper interface {
	GetParams(ctx sdk.Context) (params minttypes.Params)
	SetParams(ctx sdk.Context, params minttypes.Params)
}
