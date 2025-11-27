package wasmbindings

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	"github.com/cosmos/gogoproto/proto"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	liquidstaketypes "github.com/persistenceOne/persistenceCore/v16/x/liquidstake/types"
)

// RegisterStargateQueries returns wasm options for the stargate querier.
func RegisterStargateQueries(
	queryRouter *baseapp.GRPCQueryRouter, codec codec.Codec,
) []wasmkeeper.Option {
	acceptList := wasmkeeper.AcceptedQueries{

		// used by dex
		"/cosmos.gov.v1.Query/Proposal":  func() proto.Message { return &govtypes.QueryProposalResponse{} },
		"/cosmos.gov.v1.Query/Proposals": func() proto.Message { return &govtypes.QueryProposalsResponse{} },
		"/cosmos.gov.v1.Query/Deposit":   func() proto.Message { return &govtypes.QueryDepositResponse{} },
		"/cosmos.gov.v1.Query/Params":    func() proto.Message { return &govtypes.QueryParamsResponse{} },

		//used by dex
		"/pstake.liquidstake.v1beta1.Query/States": func() proto.Message { return &liquidstaketypes.QueryStatesResponse{} },

		//"/ibc.applications.transfer.v1.Query/DenomTrace": func() proto.Message { return &ibctransfertypes.QueryDenomTraceResponse{} },
		"/ibc.applications.transfer.v1.Query/Denoms":    func() proto.Message { return &ibctransfertypes.QueryDenomsResponse{} },
		"/ibc.applications.transfer.v1.Query/Denom":     func() proto.Message { return &ibctransfertypes.QueryDenomResponse{} },
		"/ibc.applications.transfer.v1.Query/DenomHash": func() proto.Message { return &ibctransfertypes.QueryDenomHashResponse{} },
	}

	queryPluginOpt := wasmkeeper.WithQueryPlugins(&wasmkeeper.QueryPlugins{
		Stargate: wasmkeeper.AcceptListStargateQuerier(acceptList, queryRouter, codec),
	})

	return []wasm.Option{
		queryPluginOpt,
	}
}
