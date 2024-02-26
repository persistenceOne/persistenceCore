package interchaintest

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client/tx"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	liquidstaketypes "github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"

	"go.uber.org/zap/zaptest"

	"github.com/persistenceOne/persistenceCore/v11/interchaintest/helpers"
)

var (
	blocksAfterUpgradeFast = uint64(1)
)

// TestPersistenceUpgradeLiquidstake initializes some liquidstake delegations and then runs chain upgrade with address migrations.
func TestPersistenceUpgradeLiquidstake(t *testing.T) {
	var (
		chainName            = "persistence"
		upgradeRepo          = PersistenceCoreImage.Repository
		initialVersion       = "v11.6.0"
		upgradeBranchVersion = PersistenceCoreImage.Version
		upgradeName          = "v11.7.0"
	)

	LSChainUpgradeTest(
		t,
		chainName,
		upgradeRepo,
		initialVersion,
		upgradeBranchVersion,
		upgradeName,
	)
}

func LSChainUpgradeTest(
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
		ModifyGenesis:          cosmos.ModifyGenesis(fastVotingGenesisOverridesKV),

		GasAdjustment: 10,
	}

	// create a single chain instance with 2 validators
	validatorsCount := 2

	cf := interchaintest.NewBuiltinChainFactory(
		zaptest.NewLogger(t),
		[]*interchaintest.ChainSpec{{
			Name:          chainName,
			ChainName:     chainName,
			Version:       initialVersion,
			ChainConfig:   chainCfg,
			NumValidators: &validatorsCount,
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
		SkipPathCreation: false,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = ic.Close()
	})

	const userFunds = int64(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userFunds, chain)
	chainUser := users[0]

	chainNode := chain.Nodes()[0]
	testDenom := chain.Config().Denom

	// pre-check state migrations: setup liquidstake and stake some!
	broadcaster := cosmos.NewBroadcaster(t, chain)
	broadcaster.ConfigureFactoryOptions(func(factory tx.Factory) tx.Factory {
		return factory.WithSimulateAndExecute(true)
	})

	// Updating liquidstake params for a new chain
	height, err := chain.Height(ctx)
	require.NoError(t, err, "error fetching height before submitting a proposal")

	msgUpdateParams, err := codectypes.NewAnyWithValue(&liquidstaketypes.MsgUpdateParams{
		Authority: authtypes.NewModuleAddress("gov").String(),
		Params: liquidstaketypes.Params{
			LiquidBondDenom:       liquidstaketypes.DefaultLiquidBondDenom,
			LsmDisabled:           false,
			UnstakeFeeRate:        liquidstaketypes.DefaultUnstakeFeeRate,
			MinLiquidStakeAmount:  liquidstaketypes.DefaultMinLiquidStakeAmount,
			CwLockedPoolAddress:   "",
			FeeAccountAddress:     liquidstaketypes.DummyFeeAccountAcc.String(),
			AutocompoundFeeRate:   liquidstaketypes.DefaultAutocompoundFeeRate,
			WhitelistAdminAddress: chainUser.FormattedAddress(),
			ModulePaused:          false,
		},
	})

	require.NoError(t, err, "failed to pack liquidstaketypes.MsgUpdateParams")

	txResp, err := cosmos.BroadcastTx(
		ctx,
		broadcaster,
		chainUser,
		&govv1.MsgSubmitProposal{
			InitialDeposit: []sdk.Coin{sdk.NewCoin(chain.Config().Denom, sdk.NewInt(500_000_000))},
			Proposer:       chainUser.FormattedAddress(),
			Title:          "LiquidStake Params Update",
			Summary:        "Sets params for liquidstake",
			Messages:       []*codectypes.Any{msgUpdateParams},
		},
	)
	require.NoError(t, err, "error submitting liquidstake params update tx")

	paramsUpdateTx, err := helpers.QueryProposalTx(context.Background(), chain.Nodes()[0], txResp.TxHash)
	require.NoError(t, err, "error checking proposal tx")

	err = chain.VoteOnProposalAllValidators(ctx, paramsUpdateTx.ProposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit votes")

	_, err = cosmos.PollForProposalStatus(ctx, chain, height, height+15, paramsUpdateTx.ProposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	err = testutil.WaitForBlocks(ctx, 2, chain)
	require.NoError(t, err)

	// Get list of validators
	validators := helpers.QueryAllValidators(t, ctx, chainNode)
	require.Len(t, validators, validatorsCount, "validators returned must match count of validators created")

	// Update whitelisted validators list from the first user (just for convenience)
	txResp, err = cosmos.BroadcastTx(
		ctx,
		broadcaster,
		chainUser,
		&liquidstaketypes.MsgUpdateWhitelistedValidators{
			Authority: chainUser.FormattedAddress(),
			WhitelistedValidators: []liquidstaketypes.WhitelistedValidator{{
				ValidatorAddress: validators[0].OperatorAddress,
				TargetWeight:     math.NewInt(5000),
			}, {
				ValidatorAddress: validators[1].OperatorAddress,
				TargetWeight:     math.NewInt(5000),
			}},
		},
	)
	require.NoError(t, err, "error submitting liquidstake validators whitelist update tx")

	// Liquid stake XPRT from the first user (5 XPRT)

	chainUserLiquidStakeAmount := sdk.NewInt(5_000_000)
	chainUserLiquidStakeCoins := sdk.NewCoin(testDenom, chainUserLiquidStakeAmount)

	txResp, err = cosmos.BroadcastTx(
		ctx,
		broadcaster,
		chainUser,
		&liquidstaketypes.MsgLiquidStake{
			DelegatorAddress: chainUser.FormattedAddress(),
			Amount:           chainUserLiquidStakeCoins,
		},
	)
	require.NoError(t, err)

	_, err = helpers.QueryTx(ctx, chainNode, txResp.TxHash)
	require.NoError(t, err)

	stkXPRTBalance, err := chain.GetBalance(ctx, chainUser.FormattedAddress(), "stk/uxprt")
	require.NoError(t, err)
	require.Equal(t, chainUserLiquidStakeAmount, stkXPRTBalance, "stkXPRT balance must match the liquid-staked amount")

	// end of pre-check state migrations

	height, err = chain.Height(ctx)
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

	broadcaster = cosmos.NewBroadcaster(t, chain)
	txResp, err = cosmos.BroadcastTx(
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

	err = testutil.WaitForBlocks(timeoutCtx, int(blocksAfterUpgradeFast), chain)
	require.NoError(t, err, "chain did not produce blocks after upgrade")

	height, err = chain.Height(ctx)
	require.NoError(t, err, "error fetching height after upgrade")

	require.GreaterOrEqual(t, height, haltHeight+blocksAfterUpgradeFast, "height did not increment enough after upgrade")

	// TODO: post-check state migrations

	require.FailNow(t, "required a failure to dump the logs")
}
