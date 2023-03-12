package wasmbindings

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterStargateQueries returns wasm options for the stargate querier.
func RegisterStargateQueries(
	queryRouter baseapp.GRPCQueryRouter, codec codec.Codec,
) []wasmkeeper.Option {
	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Stargate: stargateQuerier(queryRouter, codec),
	})

	return []wasm.Option{
		queryPluginOpt,
	}
}
