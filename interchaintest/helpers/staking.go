package helpers

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/stretchr/testify/require"
)

// QueryAllValidators lists all validators
func QueryAllValidators(t *testing.T, ctx context.Context, chainNode *cosmos.ChainNode) []Validator {
	stdout, _, err := chainNode.ExecQuery(ctx, "staking", "validators")
	require.NoError(t, err)

	debugOutput(t, string(stdout))

	var resp queryValidatorsResponse
	err = json.Unmarshal([]byte(stdout), &resp)
	require.NoError(t, err)

	return resp.Validators
}

// QueryValidator gets info about particular validator
func QueryValidator(
	t *testing.T,
	ctx context.Context,
	chainNode *cosmos.ChainNode,
	valoperAddr string,
) Validator {
	stdout, _, err := chainNode.ExecQuery(ctx, "staking", "validator", valoperAddr)
	require.NoError(t, err)

	debugOutput(t, string(stdout))

	var validator Validator
	err = json.Unmarshal([]byte(stdout), &validator)
	require.NoError(t, err)

	return validator
}

// QueryDelegation gets info about particular delegation
func QueryDelegation(
	t *testing.T,
	ctx context.Context,
	chainNode *cosmos.ChainNode,
	delegatorAddr string,
	valoperAddr string,
) Delegation {
	stdout, _, err := chainNode.ExecQuery(ctx, "staking", "delegation", delegatorAddr, valoperAddr)
	require.NoError(t, err)

	debugOutput(t, string(stdout))

	var resp queryDelegationResponse
	err = json.Unmarshal([]byte(stdout), &resp)
	require.NoError(t, err)

	return resp.Delegation
}

type queryDelegationResponse struct {
	Delegation Delegation `json:"delegation"`
}

type Delegation struct {
	DelegatorAddress string  `json:"delegator_address"`
	ValidatorAddress string  `json:"validator_address"`
	Shares           sdk.Dec `json:"shares"`
	ValidatorBond    bool    `json:"validator_bond"`
}

type queryValidatorsResponse struct {
	Validators []Validator `json:"validators"`
}

type Validator struct {
	OperatorAddress     string    `json:"operator_address"`
	Jailed              bool      `json:"jailed"`
	Status              string    `json:"status"`
	Tokens              sdk.Int   `json:"tokens"`
	DelegatorShares     sdk.Dec   `json:"delegator_shares"`
	UnbondingHeight     int64     `json:"unbonding_height,string"`
	UnbondingTime       time.Time `json:"unbonding_time"`
	ValidatorBondShares sdk.Dec   `json:"validator_bond_shares"`
	LiquidShares        sdk.Dec   `json:"liquid_shares"`
}
