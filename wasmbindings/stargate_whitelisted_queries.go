package wasmbindings

import (
	"fmt"
	"sync"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	oracletypes "github.com/persistenceOne/persistence-sdk/v2/x/oracle/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
)

// StargateQueries is a map of stargate queries registered for the contract
var stargateWhitelistQueries sync.Map

func init() {
	// stargate queries available for the contract
	setStargateWhitelistQuery("/persistence.oracle.v1beta1.Query/ExchangeRate", &oracletypes.QueryExchangeRateResponse{})

	// governance module
	setStargateWhitelistQuery("/cosmos.gov.v1.Query/Proposal", &govtypes.QueryProposalResponse{})
	setStargateWhitelistQuery("/cosmos.gov.v1.Query/Proposals", &govtypes.QueryProposalsResponse{})
	setStargateWhitelistQuery("/cosmos.gov.v1.Query/Deposit", &govtypes.QueryDepositResponse{})
	setStargateWhitelistQuery("/cosmos.gov.v1.Query/Params", &govtypes.QueryParamsResponse{})

	// liquid staking module for exchange rate query from a contract
	setStargateWhitelistQuery("/pstake.liquidstakeibc.v1beta1.Query/ExchangeRate", &liquidstakeibctypes.QueryExchangeRateResponse{})
}

// setStargateWhitelistQuery stores the stargate queries.
func setStargateWhitelistQuery(path string, responseType codec.ProtoMarshaler) {
	stargateWhitelistQueries.Store(path, responseType)
}

// GetStargateWhitelistedQuery returns the stargate query based on the query path.
func GetStargateWhitelistedQuery(path string) (codec.ProtoMarshaler, error) {
	responseType, ok := stargateWhitelistQueries.Load(path)
	if !ok {
		return nil, wasmvmtypes.UnsupportedRequest{Kind: fmt.Sprintf("'%s' path is not allowed from the contract", path)}
	}

	resp, ok := responseType.(codec.ProtoMarshaler)
	if !ok {
		return nil, wasmvmtypes.Unknown{}
	}

	return resp, nil
}
