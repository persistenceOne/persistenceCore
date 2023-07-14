package interchaintest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"

	"github.com/persistenceOne/persistenceCore/v8/interchaintest/helpers"
)

// TestSkipMevAuction tests that x/builder corretly wired and allows to make auctions to prioritise txns
func TestSkipMevAuction(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// override SDK beck prefixes with chain specific
	helpers.SetConfig()

	// Base setup
	chains := CreateThisBranchChain(t, 1, 0)
	ic, ctx, _, _ := BuildInitialChain(t, chains)
	chain := chains[0].(*cosmos.CosmosChain)
	testDenom := chain.Config().Denom

	require.NotNil(t, ic)
	require.NotNil(t, ctx)

	t.Cleanup(func() {
		_ = ic.Close()
	})

	const userFunds = int64(10_000_000_000)

	chainUserMnemonic := helpers.NewMnemonic()
	chainUser, err := interchaintest.GetAndFundTestUserWithMnemonic(ctx, t.Name(), chainUserMnemonic, userFunds, chain)
	require.NoError(t, err)

	chainNode := chain.Nodes()[0]

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

	require.Equal(t, int64(100), recipientBalance, "recipient must have balance")

	// TODO: verify that tx is actually prioritised over other txns in the block
	// The best way to do so it by using a wasm counter contract, but it requires some more orchestration
	// Send thress txns: [low bid, higher bid, normal tx] three times.
}
