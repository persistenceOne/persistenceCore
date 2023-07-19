package helpers

import (
	"context"
	"encoding/json"
	"testing"

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

type queryValidatorsResponse struct {
	Validators []Validator `json:"validators"`
}

type Validator struct {
	OperatorAddress string `json:"operator_address"`
}
