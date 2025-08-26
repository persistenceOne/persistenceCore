package helpers

import (
	"context"
	"encoding/json"
	"testing"

	"cosmossdk.io/math"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	"github.com/stretchr/testify/require"
)

func GetTotalAmountLocked(
	t *testing.T,
	ctx context.Context,
	chainNode *cosmos.ChainNode,
	contract,
	uaddr string,
) math.Int {
	query, err := json.Marshal(QueryMsg{
		GetTotalAmountLocked: &GetTotalAmountLockedQuery{
			AssetInfo: AssetInfo{
				NativeToken: NativeTokenInfo{
					Denom: "stk/uxprt",
				},
			},
			User: uaddr,
		},
	})
	require.NoError(t, err)

	stdout, _, err := chainNode.ExecQuery(ctx, "wasm", "contract-state", "smart", contract, string(query))
	require.NoError(t, err, "error querying superfluid LP contract for locked LST amount")

	debugOutput(t, string(stdout))

	var res GetTotalAmountLockedResponse

	err = json.Unmarshal([]byte(stdout), &res)
	require.NoError(t, err)

	return math.Int(res.Data)
}
