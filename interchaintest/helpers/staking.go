package helpers

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	"github.com/stretchr/testify/require"
)

// QueryAllValidators lists all validators
func QueryAllValidators(t *testing.T, ctx context.Context, chainNode *cosmos.ChainNode) []Validator {
	stdout, _, err := chainNode.ExecQuery(ctx, "staking", "validators")
	require.NoError(t, err)

	debugOutput(t, string(stdout))

	var resp queryValidatorsResponse
	err = json.Unmarshal(stdout, &resp)
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
	err = json.Unmarshal(stdout, &validator)
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

	var resp queryDelegationsResponseWrapper
	err = json.Unmarshal(stdout, &resp)
	require.NoError(t, err)

	return resp.QueryDelegationResponse.Delegation
}

// QueryUnbondingDelegations gets info about all unbonding delegations for a delegator
func QueryUnbondingDelegations(
	t *testing.T,
	ctx context.Context,
	chainNode *cosmos.ChainNode,
	delegatorAddr string,
) []UnbondingDelegation {
	stdout, _, err := chainNode.ExecQuery(ctx, "staking", "unbonding-delegations", delegatorAddr)
	require.NoError(t, err)

	debugOutput(t, string(stdout))

	var resp queryUnbondingDelegationsResponse
	err = json.Unmarshal(stdout, &resp)
	require.NoError(t, err)

	return resp.UnbondingResponses
}

// QueryUnbondingDelegation gets info about particular unbonding delegation
func QueryUnbondingDelegation(
	t *testing.T,
	ctx context.Context,
	chainNode *cosmos.ChainNode,
	delegatorAddr string,
	valoperAddr string,
) UnbondingDelegation {
	stdout, _, err := chainNode.ExecQuery(ctx, "staking", "unbonding-delegation", delegatorAddr, valoperAddr)
	require.NoError(t, err)

	debugOutput(t, string(stdout))

	var resp UnbondingDelegation
	err = json.Unmarshal([]byte(stdout), &resp)
	require.NoError(t, err)

	return resp
}

// QueryTotalLiquidStaked returns amount of tokens in liquid staking protocol globally (LSM, ICA, stkxprt)
func QueryTotalLiquidStaked(
	t *testing.T,
	ctx context.Context,
	chainNode *cosmos.ChainNode,
) math.Int {
	stdout, _, err := chainNode.ExecQuery(ctx, "staking", "total-liquid-staked")
	require.NoError(t, err)

	debugOutput(t, string(stdout))

	var resp queryTotalLiquidStaked
	err = json.Unmarshal(stdout, &resp)
	require.NoError(t, err)

	return resp.Tokens
}

type queryTotalLiquidStaked struct {
	Tokens math.Int `json:"tokens"`
}

type queryDelegationsResponseWrapper struct {
	QueryDelegationResponse queryDelegationResponse `json:"delegation_response"`
}

type queryDelegationResponse struct {
	Delegation Delegation `json:"delegation"`
}

type Delegation struct {
	DelegatorAddress string         `json:"delegator_address"`
	ValidatorAddress string         `json:"validator_address"`
	Shares           math.LegacyDec `json:"shares"`
	ValidatorBond    bool           `json:"validator_bond"`
}

type queryValidatorsResponse struct {
	Validators []Validator `json:"validators"`
}

type queryUnbondingDelegationsResponse struct {
	UnbondingResponses []UnbondingDelegation `json:"unbonding_responses"`
}

type Validator struct {
	OperatorAddress     string         `json:"operator_address"`
	Jailed              bool           `json:"jailed"`
	Status              string         `json:"status"`
	Tokens              math.Int       `json:"tokens"`
	DelegatorShares     math.LegacyDec `json:"delegator_shares"`
	UnbondingHeight     int64          `json:"unbonding_height,string"`
	UnbondingTime       time.Time      `json:"unbonding_time"`
	ValidatorBondShares math.LegacyDec `json:"validator_bond_shares"`
	LiquidShares        math.LegacyDec `json:"liquid_shares"`
}

type UnbondingDelegation struct {
	DelegatorAddress string `json:"delegator_address"`
	ValidatorAddress string `json:"validator_address"`

	Entries []UnbondingDelegationEntry `json:"entries"`
}

type UnbondingDelegationEntry struct {
	CreationHeight          string         `json:"creation_height"`
	CompletionTime          time.Time      `json:"completion_time"`
	InitialBalance          math.Int       `json:"initial_balance"`
	Balance                 math.Int       `json:"balance"`
	UnbondingID             string         `json:"unbonding_id"`
	UnbondingOnHoldRefCount string         `json:"unbonding_on_hold_ref_count"`
	ValidatorBondFactor     math.LegacyDec `json:"validator_bond_factor"`
	GlobalLiquidStakingCap  math.LegacyDec `json:"global_liquid_staking_cap"`
}
