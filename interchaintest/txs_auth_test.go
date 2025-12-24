package interchaintest

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/interchaintest/v10"
	"github.com/cosmos/interchaintest/v10/chain/cosmos"
	"github.com/cosmos/interchaintest/v10/ibc"
	"github.com/cosmos/interchaintest/v10/testutil"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v17/interchaintest/helpers"
)

// TestTxAuthSignModesAndOrdering executes 6 independent bank sends to cover:
// - ordered (default delivery): direct, amino-json, textual
// - unordered (with timeout): direct, amino-json, textual
// Each tx uses a distinct from_user, so no sleeps are required and execution is fast.
func TestTxAuthSignModesAndOrdering(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()
	t.Cleanup(func() {})

	// Single chain with 1 validator is sufficient
	validators := 1
	ic, chain := CreateChain(t, ctx, validators, 0)
	require.NotNil(t, ic)
	require.NotNil(t, chain)
	defer func() { _ = ic.Close() }()

	// ensure chain has produced at least one block before first tx
	require.NoError(t, testutil.WaitForBlocks(ctx, 1, chain))

	chainNode := chain.Nodes()[0]
	denom := chain.Config().Denom

	// Create 6 independent senders and 1 common recipient
	fromFunds := math.NewInt(1_000_000) // enough for amount + fees
	var senders []ibc.Wallet
	for i := 1; i <= 6; i++ {
		name := fmt.Sprintf("%s-from-%d", t.Name(), i)
		u := interchaintest.GetAndFundTestUsers(t, ctx, name, fromFunds, chain)[0]
		senders = append(senders, u)
	}

	toFunds := math.NewInt(1_000_000)
	toUser := interchaintest.GetAndFundTestUsers(t, ctx, fmt.Sprintf("%s-to", t.Name()), toFunds, chain)[0]

	amount := sdk.NewCoin(denom, math.NewInt(100_000))

	// retry wrapper to reduce flakiness due to transient RPC hiccups in CI
	execTxWithRetry := func(ctx context.Context, node *cosmos.ChainNode, key string, cmd ...string) (string, error) {
		var lastErr error
		for i := 0; i < i; i++ {
			t.Logf("Exec attempt %d: %v", i+1, append([]string{"persistenceCore", "tx"}, cmd...))
			txHash, err := node.ExecTx(ctx, key, cmd...)
			if err == nil {
				return txHash, nil
			}
			lastErr = err
			emsg := err.Error()
			// retry on typical transient errors observed in CI
			if strings.Contains(emsg, "connection refused") || strings.Contains(emsg, "post failed") || strings.Contains(emsg, "EOF") || strings.Contains(emsg, "i/o timeout") || strings.Contains(emsg, "transport is closing") {
				time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
				continue
			}
			// non-transient
			return "", err
		}
		return "", lastErr
	}

	doSend := func(sender ibc.Wallet, signMode string, unordered bool) {
		cmd := []string{
			"bank", "send",
			sender.FormattedAddress(),
			toUser.FormattedAddress(),
			amount.String(),
			"--gas=auto",
			"--sign-mode", signMode,
		}
		if unordered {
			cmd = append(cmd, "--unordered", "--timeout-duration=10s")
		}

		txHash, err := execTxWithRetry(ctx, chainNode, sender.KeyName(), cmd...)
		require.NoError(t, err)
		_, err = helpers.QueryTx(ctx, chainNode, txHash)
		require.NoError(t, err)
	}

	beforeTo, err := chain.GetBalance(ctx, toUser.FormattedAddress(), denom)
	require.NoError(t, err)

	// Ordered (default) deliveries
	doSend(senders[0], "direct", false)
	doSend(senders[1], "amino-json", false)
	doSend(senders[2], "textual", false)

	// Unordered deliveries (with timeout)
	doSend(senders[3], "direct", true)
	doSend(senders[4], "amino-json", true)
	doSend(senders[5], "textual", true)

	afterTo, err := chain.GetBalance(ctx, toUser.FormattedAddress(), denom)
	require.NoError(t, err)
	require.Equal(t, beforeTo.Add(amount.Amount.MulRaw(6)), afterTo, "recipient should receive 6 transfers")
}
