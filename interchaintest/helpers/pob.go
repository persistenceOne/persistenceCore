package helpers

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/stretchr/testify/require"
)

// BuilderAuctionBid creates an auction bid transaction with signed bundled transactions
func BuilderAuctionBid(
	t *testing.T,
	ctx context.Context,
	chain *cosmos.CosmosChain,
	user ibc.Wallet,
	bidder string,
	bid sdk.Coin,
	timeoutHeight uint64,
	transactions ...sdk.Tx,
) {
	txBytes := make([]string, 0, len(transactions))
	for _, tx := range transactions {
		bz, err := chain.Config().EncodingConfig.TxConfig.TxEncoder()(tx)
		if err != nil {
			require.NoError(t, err)
			return
		}

		txBytes = append(txBytes, fmt.Sprintf("%2X", bz))
	}

	//  persistenceCore tx builder auction-bid [bidder] [bid] [bundled_tx1,bundled_tx2,...,bundled_txN]
	cmd := append([]string{
		"builder", "auction-bid", bidder, bid.String(),
	}, txBytes...)

	// NOTE: --timeout-height is mandatory
	cmd = append(cmd, fmt.Sprintf("--timeout-height=%d", timeoutHeight))

	chainNode := chain.Nodes()[0]
	txHash, err := chainNode.ExecTx(ctx, user.KeyName(), cmd...)
	require.NoError(t, err)

	stdout, _, err := chainNode.ExecQuery(ctx, "tx", txHash)
	require.NoError(t, err)

	debugOutput(t, string(stdout))
}
