package wasmbindings

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	oracletypes "github.com/persistenceOne/persistence-sdk/v3/x/oracle/types"

	liquidstaketypes "github.com/persistenceOne/pstake-native/v3/x/liquidstake/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v3/x/liquidstakeibc/types"
)

// RegisterStargateQueries returns wasm options for the stargate querier.
func RegisterStargateQueries(
	queryRouter *baseapp.GRPCQueryRouter, codec codec.Codec,
) []wasmkeeper.Option {
	acceptList := wasmkeeper.AcceptedStargateQueries{
		"/persistence.oracle.v1beta1.Query/ExchangeRate": &oracletypes.QueryExchangeRateResponse{},

		"/cosmos.gov.v1.Query/Proposal":  &govtypes.QueryProposalResponse{},
		"/cosmos.gov.v1.Query/Proposals": &govtypes.QueryProposalsResponse{},
		"/cosmos.gov.v1.Query/Deposit":   &govtypes.QueryDepositResponse{},
		"/cosmos.gov.v1.Query/Params":    &govtypes.QueryParamsResponse{},

		"/pstake.liquidstakeibc.v1beta1.Query/ExchangeRate": &liquidstakeibctypes.QueryExchangeRateResponse{},
		"/pstake.liquidstakeibc.v1beta1.Query/HostChains":   &liquidstakeibctypes.QueryHostChainsResponse{},
		"/pstake.liquidstake.v1beta1.Query/States":          &liquidstaketypes.QueryStatesResponse{},

		"/ibc.applications.transfer.v1.Query/DenomTrace": &ibctransfertypes.QueryDenomTraceResponse{},
	}

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Stargate: wasmkeeper.AcceptListStargateQuerier(acceptList, queryRouter, codec),
	})

	return []wasm.Option{
		queryPluginOpt,
	}
}
