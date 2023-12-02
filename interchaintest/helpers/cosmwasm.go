package helpers

import (
	"context"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/stretchr/testify/require"
)

func SetupContract(
	t *testing.T,
	ctx context.Context,
	chain *cosmos.CosmosChain,
	keyname string,
	fileLoc string,
	message string,
) (codeId, contract string) {
	codeId, err := chain.StoreContract(ctx, keyname, fileLoc)
	if err != nil {
		t.Fatal(err)
	}

	contractAddr, err := chain.InstantiateContract(ctx, keyname, codeId, message, true)
	if err != nil {
		t.Fatal(err)
	}

	return codeId, contractAddr
}

func ExecuteMsgWithAmount(
	t *testing.T,
	ctx context.Context,
	chain *cosmos.CosmosChain,
	user ibc.Wallet,
	contractAddr, amount, message string,
) (txHash string) {
	cmd := []string{
		"wasm", "execute", contractAddr, message,
		"--gas", "500000",
		"--amount", amount,
	}

	chainNode := chain.Nodes()[0]
	txHash, err := chainNode.ExecTx(ctx, user.KeyName(), cmd...)
	require.NoError(t, err)

	stdout, _, err := chainNode.ExecQuery(ctx, "tx", txHash)
	require.NoError(t, err)

	debugOutput(t, string(stdout))

	return txHash
}

func ExecuteMsgWithFee(
	t *testing.T,
	ctx context.Context,
	chain *cosmos.CosmosChain,
	user ibc.Wallet,
	contractAddr, amount, feeCoin, message string,
) (txHash string) {
	cmd := []string{
		"wasm", "execute", contractAddr, message,
		"--fees", feeCoin,
		"--gas", "500000",
	}

	if amount != "" {
		cmd = append(cmd, "--amount", amount)
	}

	chainNode := chain.Nodes()[0]
	txHash, err := chainNode.ExecTx(ctx, user.KeyName(), cmd...)
	require.NoError(t, err)

	stdout, _, err := chainNode.ExecQuery(ctx, "tx", txHash)
	require.NoError(t, err)

	debugOutput(t, string(stdout))

	return txHash
}
