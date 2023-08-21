package wasmbindings

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	liquidstakeibckeeper "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/keeper"
)

func RegisterCustomPlugins(
	bank *bankkeeper.BaseKeeper,
	liquidStakeIBCKeeper *liquidstakeibckeeper.Keeper,
) []wasmkeeper.Option {
	//wasmQueryPlugin := NewQueryPlugin(tokenFactory)
	//
	//queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
	//	Custom: CustomQuerier(wasmQueryPlugin),
	//})
	messengerDecoratorOpt := wasmkeeper.WithMessageHandlerDecorator(
		CustomMessageDecorator(bank, liquidStakeIBCKeeper),
	)

	return []wasm.Option{
		//queryPluginOpt,
		messengerDecoratorOpt,
	}
}

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
