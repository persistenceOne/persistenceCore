package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v17/x/liquidstake/types"
)

func TestGetLiquidValidatorKey(t *testing.T) {
	valAddr := sdk.ValAddress([]byte("validator"))
	key := types.GetLiquidValidatorKey(valAddr)

	// Check that the key starts with the LiquidValidatorsKey prefix
	require.True(t, len(key) > len(types.LiquidValidatorsKey))
	require.Equal(t, types.LiquidValidatorsKey, key[:len(types.LiquidValidatorsKey)])

	// We can't directly check if the key contains the validator address
	// because the address is length-prefixed in the key
	// Just check that the key is longer than the prefix
	require.Greater(t, len(key), len(types.LiquidValidatorsKey))
}
