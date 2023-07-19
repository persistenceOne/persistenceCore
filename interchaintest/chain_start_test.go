package interchaintest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestBasicPersistenceStart is a basic test to assert that spinning up a Persistence network with 4 validators works properly.
// This test is not run part of CI, it's for checking and troubleshooting start-up locally.
func TestBasicPersistenceStart(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	ctx, cancelFn := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancelFn()
	})

	// create a single chain instance with 4 validators
	validatorsCount := 4

	ic, chain := CreateChain(t, ctx, validatorsCount, 0)
	require.NotNil(t, ic)
	require.NotNil(t, chain)

	t.Cleanup(func() {
		_ = ic.Close()
	})
}
