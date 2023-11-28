package interchaintest

import (
	"context"
	"testing"

	"cosmossdk.io/math"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"

	"github.com/persistenceOne/persistenceCore/v11/interchaintest/helpers"
)

// TestSkipMevAuction tests that x/builder corretly wired and allows to make auctions to prioritise txns
func TestSkipMevAuction(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	ctx, cancelFn := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancelFn()
	})

	// override SDK beck prefixes with chain specific
	helpers.SetConfig()

	// create a single chain instance with 1 validator
	validatorsCount := 1
	ic, chain := CreateChain(t, ctx, validatorsCount, 0)
	require.NotNil(t, ic)
	require.NotNil(t, chain)

	chainNode := chain.Nodes()[0]
	testDenom := chain.Config().Denom

	t.Cleanup(func() {
		_ = ic.Close()
	})

	const userFunds = int64(10_000_000_000)

	chainUserMnemonic := helpers.NewMnemonic()
	chainUser, err := interchaintest.GetAndFundTestUserWithMnemonic(ctx, t.Name(), chainUserMnemonic, userFunds, chain)
	require.NoError(t, err)

	paramsStdout, _, err := chainNode.ExecQuery(ctx, "builder", "params")
	require.NoError(t, err)
	t.Log("checking POB params", string(paramsStdout))

	kb, err := helpers.NewKeyringFromMnemonic(cosmos.DefaultEncoding().Codec, chainUser.KeyName(), chainUserMnemonic)
	require.NoError(t, err)

	clientContext := chainNode.CliContext()
	clientContext = clientContext.WithCodec(chain.Config().EncodingConfig.Codec)

	txFactory := helpers.NewTxFactory(clientContext)
	txFactory = txFactory.WithKeybase(kb)
	txFactory = txFactory.WithTxConfig(persistenceEncoding().TxConfig)

	accountRetriever := authtypes.AccountRetriever{}
	accountNum, currentSeq, err := accountRetriever.GetAccountNumberSequence(clientContext, chainUser.Address())
	require.NoError(t, err)

	txFactory = txFactory.WithAccountNumber(accountNum)
	// the tx that we put on auction will have the next sequence
	txFactory = txFactory.WithSequence(currentSeq + 1)

	txn, err := txFactory.BuildUnsignedTx(&banktypes.MsgSend{
		FromAddress: chainUser.FormattedAddress(),
		ToAddress:   helpers.SomeoneAddress.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin(testDenom, sdk.NewInt(100))),
	})
	require.NoError(t, err)

	currentHeight, err := chain.Height(ctx)
	require.NoError(t, err)

	// transaction simulation there is possible, but we skip it for now
	txn.SetGasLimit(100000)
	txn.SetTimeoutHeight(currentHeight + 5)

	err = tx.Sign(txFactory, chainUser.KeyName(), txn, true)
	require.NoError(t, err)

	auctionBid := sdk.NewCoin(testDenom, sdk.NewInt(100))
	helpers.BuilderAuctionBid(
		t, ctx, chain,
		chainUser,
		chainUser.FormattedAddress(),
		auctionBid,
		currentHeight+5,
		txn.GetTx(),
	)

	recipientBalance, err := chain.GetBalance(ctx, helpers.SomeoneAddress.String(), testDenom)
	require.NoError(t, err)

	require.Equal(t, math.NewInt(100), recipientBalance, "recipient must have balance")

	// TODO: verify that tx is actually prioritised over other txns in the block
	// The best way to do so it by using a wasm counter contract, but it requires some more orchestration
	// Send thress txns: [low bid, higher bid, normal tx] three times.
}
