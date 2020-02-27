package asset

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct{}

func DefaultGenesisState() GenesisState { return GenesisState{} }

func ValidateGenesis(genesisState GenesisState) error { return nil }

func InitializeGenesisState(context sdkTypes.Context, keeper Keeper, genesisState GenesisState) {
}
func ExportGenesis(context sdkTypes.Context, keeper Keeper) GenesisState {
	return GenesisState{}
}
