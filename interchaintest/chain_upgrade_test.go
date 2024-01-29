package interchaintest

import (
	"context"
	"fmt"
	"testing"
	"time"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"

	"github.com/persistenceOne/persistenceCore/v10/interchaintest/helpers"
)

const (
	haltHeightDelta    = uint64(10) // will propose upgrade this many blocks in the future
	blocksAfterUpgrade = uint64(7)
)

func TestPersistenceUpgradeBasic(t *testing.T) {
	var (
		chainName            = "persistence"
		upgradeRepo          = PersistenceCoreImage.Repository
		initialVersion       = "v10.3.0"
		upgradeBranchVersion = PersistenceCoreImage.Version
		upgradeName          = "v10.4.0"
	)

	CosmosChainUpgradeTest(
		t,
		chainName,
		upgradeRepo,
		initialVersion,
		upgradeBranchVersion,
		upgradeName,
	)
}

func CosmosChainUpgradeTest(
	t *testing.T,
	chainName,
	upgradeRepo,
	initialVersion,
	upgradeBranchVersion,
	upgradeName string,
) {
	if testing.Short() {
		t.Skip("skipping in short mode")
	}

	t.Parallel()

	t.Log(chainName, initialVersion, upgradeBranchVersion, upgradeRepo, upgradeName)

	// override SDK beck prefixes with chain specific
	helpers.SetConfig()

	chainCfg := ibc.ChainConfig{
		Type:         "cosmos",
		Name:         "persistence",
		ChainID:      "ictest-core-1",
		Bin:          "persistenceCore",
		Bech32Prefix: "persistence",

		Images: []ibc.DockerImage{{
			Repository: PersistenceE2ERepo,
			Version:    initialVersion,
			UidGid:     PersistenceCoreImage.UidGid,
		}},

		UsingNewGenesisCommand: true,
		GasPrices:              fmt.Sprintf("0%s", helpers.PersistenceBondDenom),
		EncodingConfig:         persistenceEncoding(),
		ModifyGenesis:          cosmos.ModifyGenesis(defaultGenesisOverridesKV),
	}

	cf := interchaintest.NewBuiltinChainFactory(
		zaptest.NewLogger(t),
		[]*interchaintest.ChainSpec{{
			Name:        chainName,
			ChainName:   chainName,
			Version:     initialVersion,
			ChainConfig: chainCfg,
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

	const userFunds = int64(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, chain)
	chainUser := users[0]

	// TODO: pre-check state migrations

	height, err := chain.Height(ctx)
	require.NoError(t, err, "error fetching height before submit upgrade proposal")

	haltHeight := height + haltHeightDelta

	msgSoftwareUpgrade, err := codectypes.NewAnyWithValue(&upgradetypes.MsgSoftwareUpgrade{
		Authority: authtypes.NewModuleAddress("gov").String(),
		Plan: upgradetypes.Plan{
			Name:   upgradeName,
			Height: int64(haltHeight),
			Info:   upgradeName + " chain upgrade",
		},
	})

	require.NoError(t, err, "failed to pack upgradetypes.SoftwareUpgradeProposal")

	broadcaster := cosmos.NewBroadcaster(t, chain)
	txResp, err := cosmos.BroadcastTx(
		ctx,
		broadcaster,
		chainUser,
		&govv1.MsgSubmitProposal{
			InitialDeposit: []sdk.Coin{sdk.NewCoin(chain.Config().Denom, sdk.NewInt(500_000_000))},
			Proposer:       chainUser.FormattedAddress(),
			Title:          "Chain Upgrade 1",
			Summary:        "First chain software upgrade",
			Messages:       []*codectypes.Any{msgSoftwareUpgrade},
		},
	)
	require.NoError(t, err, "error submitting software upgrade tx")

	upgradeTx, err := helpers.QueryProposalTx(context.Background(), chain.Nodes()[0], txResp.TxHash)
	require.NoError(t, err, "error checking software upgrade tx")

	err = chain.VoteOnProposalAllValidators(ctx, upgradeTx.ProposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit votes")

	_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+haltHeightDelta, upgradeTx.ProposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	stdout, stderr, err := chain.Validators[0].ExecQuery(ctx, "upgrade", "plan")
	require.NoError(t, err, "error querying upgrade plan")

	t.Log("Upgrade", "plan_stdout:", string(stdout), "plan_stderr:", string(stderr))

	timeoutCtx, timeoutCtxCancel := context.WithTimeout(ctx, time.Second*45)
	defer timeoutCtxCancel()

	height, err = chain.Height(ctx)
	require.NoError(t, err, "error fetching height before upgrade")

	// this should timeout due to chain halt at upgrade height.
	_ = testutil.WaitForBlocks(timeoutCtx, int(haltHeight-height)+1, chain)

	height, err = chain.Height(ctx)
	require.NoError(t, err, "error fetching height after chain should have halted")

	// make sure that chain is halted
	require.Equal(t, haltHeight, height, "height is not equal to halt height")

	// bring down nodes to prepare for upgrade
	t.Log("stopping node(s)")
	err = chain.StopAllNodes(ctx)
	require.NoError(t, err, "error stopping node(s)")

	// upgrade version on all nodes
	t.Log("upgrading node(s)")
	chain.UpgradeVersion(ctx, client, upgradeRepo, upgradeBranchVersion)

	// start all nodes back up.
	// validators reach consensus on first block after upgrade height
	// and chain block production resumes.
	t.Log("starting node(s)")
	err = chain.StartAllNodes(ctx)
	require.NoError(t, err, "error starting upgraded node(s)")

	timeoutCtx, timeoutCtxCancel = context.WithTimeout(ctx, time.Second*60)
	defer timeoutCtxCancel()

	err = testutil.WaitForBlocks(timeoutCtx, int(blocksAfterUpgrade), chain)
	require.NoError(t, err, "chain did not produce blocks after upgrade")

	height, err = chain.Height(ctx)
	require.NoError(t, err, "error fetching height after upgrade")

	require.GreaterOrEqual(t, height, haltHeight+blocksAfterUpgrade, "height did not increment enough after upgrade")

	// TODO: post-check state migrations
}
