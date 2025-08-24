package interchaintest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/persistenceOne/persistenceCore/v13/interchaintest/helpers"
)

func TestPersistenceLSMHaltFork(t *testing.T) {
	var (
		chainName            = "persistence"
		upgradeRepo          = PersistenceCoreImage.Repository
		initialVersion       = "v9.1.1"
		upgradeBranchVersion = PersistenceCoreImage.Version
	)

	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	t.Log(chainName, initialVersion, upgradeBranchVersion, upgradeRepo, "LSM FORK TEST")

	// override SDK bech prefixes with chain specific
	helpers.SetConfig()

	chainCfg := ibc.ChainConfig{
		Type:         "cosmos",
		Name:         "persistence",
		ChainID:      "ictest-core-1", // this must trigger Fork() in persistenceCore
		Bin:          "persistenceCore",
		Bech32Prefix: "persistence",

		Images: []ibc.DockerImage{{
			Repository: PersistenceE2ERepo,
			Version:    initialVersion,
			UidGid:     PersistenceCoreImage.UidGid,
		}},

		GasPrices:      fmt.Sprintf("0%s", helpers.PersistenceBondDenom),
		EncodingConfig: persistenceEncoding(),
		ModifyGenesis:  cosmos.ModifyGenesis(defaultGenesisOverridesKV),
	}

	validatorsCount := 4
	fullNodesCount := 0
	cf := interchaintest.NewBuiltinChainFactory(
		zaptest.NewLogger(t),
		[]*interchaintest.ChainSpec{{
			Name:          chainName,
			ChainName:     chainName,
			Version:       initialVersion,
			ChainConfig:   chainCfg,
			NumValidators: &validatorsCount,
			NumFullNodes:  &fullNodesCount,
		}},
	)

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chain := chains[0].(*cosmos.CosmosChain)
	ic := interchaintest.NewInterchain().
		AddChain(chain)

	ctx := context.Background()
	client, network := interchaintest.DockerSetup(t)

	err = ic.Build(ctx, nil, interchaintest.InterchainBuildOptions{
		TestName:         t.Name(),
		Client:           client,
		NetworkID:        network,
		SkipPathCreation: true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = ic.Close()
	})

	chainNode := chain.Nodes()[0]
	testDenom := chain.Config().Denom

	// Allocate two chain users with funds
	firstUserFunds := math.NewInt(10_000_000_000_000)
	firstUser := interchaintest.GetAndFundTestUsers(t, ctx, firstUserName(t.Name()), firstUserFunds, chain)[0]
	secondUserFunds := math.NewInt(10_000_000_000_000)
	secondUser := interchaintest.GetAndFundTestUsers(t, ctx, secondUserName(t.Name()), secondUserFunds, chain)[0]

	validatorInitialDelegation := math.NewInt(5_000_000_000_000)

	// Ensure chain has started
	checkCtx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	err = testutil.WaitForBlocks(checkCtx, 1, chain)
	require.NoError(t, err, "error waiting for blocks produced by the chain")
	cancelFn()

	// Get list of validators
	validators := helpers.QueryAllValidators(t, ctx, chainNode)
	require.Len(t, validators, validatorsCount, "validator returned must match count of validators created")

	// Bond first user
	firstUserBondAmount := math.NewInt(100_000_000_000)
	firstUserBondCoins := sdk.NewCoin(testDenom, firstUserBondAmount)
	txHash, err := chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "delegate", validators[0].OperatorAddress, firstUserBondCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err, "failed to execute delegate tx (user A)")

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err, "failed to execute delegate tx (user A)")

	txHash, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "validator-bond", validators[0].OperatorAddress,
		"--gas=auto",
	)
	require.NoError(t, err, "failed to execute validator-bond tx (user A)")

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err, "failed to execute validator-bond tx (user A)")

	// Query list of validators
	validators = helpers.QueryAllValidators(t, ctx, chainNode)
	require.Len(t, validators, validatorsCount, "validator returned must match count of validators created")
	require.Equal(t, validators[0].ValidatorBondShares.TruncateInt(), firstUserBondAmount, "validator bond shares must match first user's bond")

	// Increase first user delegation (x2)
	txHash, err = chainNode.ExecTx(ctx, firstUser.KeyName(),
		"staking", "delegate", validators[0].OperatorAddress, firstUserBondCoins.String(),
		"--gas=auto",
	)
	require.NoError(t, err, "failed to execute delegate tx (user A)")

	_, err = helpers.QueryTx(ctx, chainNode, txHash)
	require.NoError(t, err, "failed to execute delegate tx (user A)")

	// Query list of validators
	validators = helpers.QueryAllValidators(t, ctx, chainNode)
	require.Len(t, validators, validatorsCount, "validator returned must match count of validators created")
	require.Equal(t,
		validators[0].ValidatorBondShares.TruncateInt(),
		firstUserBondAmount.Mul(math.NewInt(2)),
		"validator bond shares must match first user's bond * 2",
	)

	height, err := chain.Height(ctx)
	require.NoError(t, err, "error fetching height before triggering chain halt")

	t.Logf("on height %d, trying to halt chain", height)

	// Ensure chain is halted in ~15s
	checkCtx, cancelFn = context.WithTimeout(context.Background(), 15*time.Second)

	// Create second user delegation
	secondUserBondAmount := math.NewInt(88_000_000_000) // 88k magic number to tell from 100k that comes from A
	secondUserBondCoins := sdk.NewCoin(testDenom, secondUserBondAmount)
	_, err = chainNode.ExecTx(checkCtx, secondUser.KeyName(),
		"staking", "delegate", validators[0].OperatorAddress, secondUserBondCoins.String(),
		"--gas=auto",
	)
	require.ErrorContains(t, err, "context deadline exceeded", "expected tx to fail broadcasting due to timeout")

	// this is just in case, should confirm quickly since checkCtx was expired
	err = testutil.WaitForBlocks(checkCtx, 2, chain)
	require.Error(t, err, "chain didn't halt! the blocks were produced")
	cancelFn()

	haltHeight, err := chain.Height(ctx)
	require.NoError(t, err, "error fetching haltHeight after triggering chain halt")

	t.Logf("chain halted at height %d after 15 seconds of timeout", haltHeight)

	validators = helpers.QueryAllValidators(t, ctx, chainNode)
	require.Len(t, validators, validatorsCount, "validator returned must match count of validators created")

	// Assert the broken state

	// a) bond shares have two delegations from user A
	require.Equal(t,
		firstUserBondAmount.Mul(math.NewInt(2)),
		validators[0].ValidatorBondShares.TruncateInt(),
		"validator bond shares must match first user's bond * 2",
	)

	// b) delegator shares missing second delegation from user A
	require.Equal(t,
		validatorInitialDelegation.Add(firstUserBondAmount),
		validators[0].DelegatorShares.TruncateInt(),
		"validator delegator shares must have only first delegation of the first user",
	)

	// c) obviously, both missing delegation from user B which is stuck

	t.Log("preparing chain fork")

	// bring down nodes to prepare for upgrade
	t.Log("stopping node(s)")
	err = chain.StopAllNodes(ctx)
	require.NoError(t, err, "error stopping node(s)")

	// upgrade version on all nodes -> 0.9.2
	t.Log("upgrading node(s)")
	chain.UpgradeVersion(ctx, client, upgradeRepo, upgradeBranchVersion)

	// start all nodes back up, if migration runs correctly, block production must resume
	t.Log("starting node(s)")
	err = chain.StartAllNodes(ctx)
	require.NoError(t, err, "error starting upgraded node(s)")

	timeoutCtx, timeoutCtxCancel := context.WithTimeout(ctx, time.Second*60)
	defer timeoutCtxCancel()

	err = testutil.WaitForBlocks(timeoutCtx, int(blocksAfterUpgrade), chain)
	require.NoError(t, err, "chain did not produce blocks after upgrade")

	height, err = chain.Height(ctx)
	require.NoError(t, err, "error fetching height after upgrade")

	require.GreaterOrEqual(t, height, haltHeight+blocksAfterUpgrade, "height did not increment enough after upgrade")

	// Post-check state migrations, i.e. delegation discrepancy

	validators = helpers.QueryAllValidators(t, ctx, chainNode)
	require.Len(t, validators, validatorsCount, "validator returned must match count of validators created")

	// Assert the fixed state

	// a) bond shares have two delegations from user A
	require.Equal(t,
		firstUserBondAmount.Mul(math.NewInt(2)),
		validators[0].ValidatorBondShares.TruncateInt(),
		"validator bond shares must match first user's bond * 2",
	)

	// b) delegator shares have two delegations from user A, as well as delegation from user B
	require.Equal(t,
		validatorInitialDelegation.Add(firstUserBondAmount.Mul(math.NewInt(2)).Add(secondUserBondAmount)),
		validators[0].DelegatorShares.TruncateInt(),
		"validator delegator shares must have 2*A + B delegations",
	)
}
