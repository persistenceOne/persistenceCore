package helpers

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/stretchr/testify/require"
)

// QueryAlliance gets info about particular alliance
func QueryAlliance(
	t *testing.T,
	ctx context.Context,
	chainNode *cosmos.ChainNode,
	allianceDenom string,
) Alliance {
	stdout, _, err := chainNode.ExecQuery(ctx, "alliance", "alliance", allianceDenom)
	require.NoError(t, err)

	debugOutput(t, string(stdout))

	var allianceResp queryAllianceResponse
	err = json.Unmarshal([]byte(stdout), &allianceResp)
	require.NoError(t, err)

	return allianceResp.Alliance
}

type queryAllianceResponse struct {
	Alliance Alliance `json:"alliance"`
}
