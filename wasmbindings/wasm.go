package wasmbindings

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	oraclekeeper "github.com/persistenceOne/persistence-sdk/v2/x/oracle/keeper"
)

func RegisterCustomPlugins(
	oracleKeeper *oraclekeeper.Keeper,
) []wasmkeeper.Option {
	wasmQueryPlugin := NewQueryPlugin(oracleKeeper)

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Custom: customQuerier(wasmQueryPlugin),
	})

	// TODO: Add more custom plugins based on the contract requirement
	return []wasm.Option{
		queryPluginOpt,
	}
}
