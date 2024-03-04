package interchaintest

import (
	"context"
	"fmt"
	"testing"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	interchaintestrelayer "github.com/strangelove-ventures/interchaintest/v7/relayer"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/strangelove-ventures/interchaintest/v7/testutil"
	"github.com/stretchr/testify/require"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
	"go.uber.org/zap/zaptest"

	"github.com/persistenceOne/persistenceCore/v11/interchaintest/helpers"
)

// TestAllianceBasic spins up a Persistence and Gaia network, initializes an IBC connection between them,
// and sends an ICS20 token transfer from Gaia->Persistence. Creates an alliance, stakes ICS20 token and collects rewards.
func TestAllianceBasic(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	t.Parallel()

	// override SDK bech prefixes with chain specific
	helpers.SetConfig()

	ctx, cancelFn := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cancelFn()
	})

	// Create chain factory with Persistence and Gaia
	numVals := 1
	numFullNodes := 1

	// important overrides: fast voting for quick proposal passing, also
	// x/alliance: set rewards delay time to 0 for faster testing (default is 7 days).
	overridesKV := append([]cosmos.GenesisKV{}, fastVotingGenesisOverridesKV...)
	overridesKV = append(overridesKV, cosmos.GenesisKV{
		Key:   "app_state.alliance.params.reward_delay_time",
		Value: "0",
	})

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name: "persistence",
			ChainConfig: persistenceChainConfig(
				overridesKV...,
			),
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
		{
			Name:          "gaia",
			Version:       "v14.1.0",
			NumValidators: &numVals,
			NumFullNodes:  &numFullNodes,
		},
	})

	const (
		path = "ibc-path"
	)

	// Get chains from the chain factory
	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	client, network := interchaintest.DockerSetup(t)

	persistenceChain, gaiaChain := chains[0].(*cosmos.CosmosChain), chains[1].(*cosmos.CosmosChain)

	relayerType, relayerName := ibc.CosmosRly, "relay"

	// Get a relayer instance
	rf := interchaintest.NewBuiltinRelayerFactory(
		relayerType,
		zaptest.NewLogger(t),
		interchaintestrelayer.CustomDockerImage(IBCRelayerImage, IBCRelayerVersion, "100:1000"),
		interchaintestrelayer.StartupFlags("--processor", "events", "--block-history", "100"),
	)

	r := rf.Build(t, client, network)

	ic := interchaintest.NewInterchain().
		AddChain(persistenceChain).
		AddChain(gaiaChain).
		AddRelayer(r, relayerName).
		AddLink(interchaintest.InterchainLink{
			Chain1:  persistenceChain,
			Chain2:  gaiaChain,
			Relayer: r,
			Path:    path,
		})

	rep := testreporter.NewNopReporter()
	eRep := rep.RelayerExecReporter(t)

	require.NoError(t, ic.Build(ctx, eRep, interchaintest.InterchainBuildOptions{
		TestName:          t.Name(),
		Client:            client,
		NetworkID:         network,
		BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation:  false,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	// Create some user accounts on both chains
	userNativeTokensAmount := int64(10_000_000_000)
	users := interchaintest.GetAndFundTestUsers(t, ctx, t.Name(), userNativeTokensAmount, persistenceChain, gaiaChain)

	// Wait a few blocks for relayer to start and for user accounts to be created
	err = testutil.WaitForBlocks(ctx, 5, persistenceChain, gaiaChain)
	require.NoError(t, err)

	// Get our Bech32 encoded user addresses
	persistenceUser, gaiaUser := users[0], users[1]

	persistenceUserAddr := persistenceUser.FormattedAddress()
	gaiaUserAddr := gaiaUser.FormattedAddress()

	// Get original account balances
	persistenceOrigBal, err := persistenceChain.GetBalance(ctx, persistenceUserAddr, persistenceChain.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(userNativeTokensAmount), persistenceOrigBal)

	gaiaOrigBal, err := gaiaChain.GetBalance(ctx, gaiaUserAddr, gaiaChain.Config().Denom)
	require.NoError(t, err)
	require.Equal(t, math.NewInt(userNativeTokensAmount), gaiaOrigBal)

	// Compose an IBC transfer and send from Gaia -> Persistence
	var transferAmount = math.NewInt(1_000_000_000)
	transfer := ibc.WalletAmount{
		Address: persistenceUserAddr,
		Denom:   gaiaChain.Config().Denom,
		Amount:  transferAmount,
	}

	// [Gaia -> Persistence] channel
	channel, err := ibc.GetTransferChannel(ctx, r, eRep, gaiaChain.Config().ChainID, persistenceChain.Config().ChainID)
	require.NoError(t, err)

	gaiaHeight, err := gaiaChain.Height(ctx)
	require.NoError(t, err)

	transferTx, err := gaiaChain.SendIBCTransfer(ctx, channel.ChannelID, gaiaUserAddr, transfer, ibc.TransferOptions{})
	require.NoError(t, err)

	err = r.StartRelayer(ctx, eRep, path)
	require.NoError(t, err)

	t.Cleanup(
		func() {
			err := r.StopRelayer(ctx, eRep)
			if err != nil {
				t.Logf("an error occurred while stopping the relayer: %s", err)
			}
		},
	)

	// Poll for the ack to know the transfer was successful
	_, err = testutil.PollForAck(ctx, gaiaChain, gaiaHeight, gaiaHeight+50, transferTx.Packet)
	require.NoError(t, err)

	err = testutil.WaitForBlocks(ctx, 10, gaiaChain)
	require.NoError(t, err)

	// Get the IBC denom for uatom on Persistence
	gaiaTokenDenom := transfertypes.GetPrefixedDenom(channel.Counterparty.PortID, channel.Counterparty.ChannelID, gaiaChain.Config().Denom)
	gaiaIBCDenom := transfertypes.ParseDenomTrace(gaiaTokenDenom).IBCDenom()

	t.Logf("[1] Expecting some balance of %s (uatom over IBC) on Persistence chain", gaiaIBCDenom)
	gaiaIBCDenomBalanceOnPersistence, err := persistenceChain.GetBalance(ctx, persistenceUserAddr, gaiaIBCDenom)
	require.NoError(t, err)
	require.Equal(t, transferAmount, gaiaIBCDenomBalanceOnPersistence,
		"uatom balance on Persistence chain must match amount transferred")

	persistenceChainNode := persistenceChain.Nodes()[0]

	// Get list of validators
	validators := helpers.QueryAllValidators(t, ctx, persistenceChainNode)
	require.Len(t, validators, numVals, "validators returned must match count of validators created")

	var totalTokensDelegated = math.NewInt(0)
	for _, val := range validators {
		totalTokensDelegated = totalTokensDelegated.Add(val.Tokens)
	}
	t.Logf("[2] Total Tokens Delegated: %s uxprt", totalTokensDelegated.String())

	// create an alliance using a governance proposal
	height, err := persistenceChain.Height(ctx)
	require.NoError(t, err, "error fetching height before submitting a proposal")

	msgCreateAlliance, err := codectypes.NewAnyWithValue(&alliancetypes.MsgCreateAlliance{
		Authority:            authtypes.NewModuleAddress("gov").String(),
		Denom:                gaiaIBCDenom,
		RewardWeight:         math.LegacyOneDec(),
		TakeRate:             math.LegacyZeroDec(),
		RewardChangeRate:     math.LegacyOneDec(),
		RewardChangeInterval: 0,
		RewardWeightRange: alliancetypes.RewardWeightRange{
			Min: math.LegacyOneDec(),
			Max: math.LegacyOneDec(),
		},
	})

	require.NoError(t, err, "failed to pack alliancetypes.MsgCreateAlliance")

	broadcaster := cosmos.NewBroadcaster(t, persistenceChain)
	txResp, err := cosmos.BroadcastTx(
		ctx,
		broadcaster,
		persistenceUser,
		&govv1.MsgSubmitProposal{
			InitialDeposit: []sdk.Coin{sdk.NewCoin(persistenceChain.Config().Denom, sdk.NewInt(500_000_000))},
			Proposer:       persistenceUserAddr,
			Title:          "Create IBC-ATOM alliance",
			Summary:        fmt.Sprintf("Creates %s alliance via proposal msg", gaiaIBCDenom),
			Messages:       []*codectypes.Any{msgCreateAlliance},
		},
	)
	require.NoError(t, err, "error submitting alliance creation proposal tx")
	if txResp.Code != 0 {
		require.FailNowf(t, "alliance creation proposal err", "proposal transaction error: code %d %s", txResp.Code, txResp.RawLog)
	}

	upgradeTx, err := helpers.QueryProposalTx(context.Background(), persistenceChainNode, txResp.TxHash)
	require.NoError(t, err, "error checking proposal tx")

	err = persistenceChain.VoteOnProposalAllValidators(ctx, upgradeTx.ProposalID, cosmos.ProposalVoteYes)
	require.NoError(t, err, "failed to submit votes")

	_, err = cosmos.PollForProposalStatus(ctx, persistenceChain, height, height+15, upgradeTx.ProposalID, cosmos.ProposalStatusPassed)
	require.NoError(t, err, "proposal status did not change to passed in expected number of blocks")

	err = testutil.WaitForBlocks(ctx, 2, persistenceChain)
	require.NoError(t, err)

	alliancesList, _, err := persistenceChainNode.ExecQuery(ctx, "alliance", "alliances")
	require.NoError(t, err)
	t.Logf("[3] Alliances available: %s", string(alliancesList))

	// Delegate alliance-enabled tokens to a validator - all user's balance goes there
	allianceDelegationCoins := sdk.NewCoin(gaiaIBCDenom, gaiaIBCDenomBalanceOnPersistence)
	txResp, err = cosmos.BroadcastTx(
		ctx,
		broadcaster,
		persistenceUser,
		&alliancetypes.MsgDelegate{
			DelegatorAddress: persistenceUserAddr,
			ValidatorAddress: validators[0].OperatorAddress,
			Amount:           allianceDelegationCoins,
		},
	)
	require.NoError(t, err, "error submitting alliance delegation tx")
	if txResp.Code != 0 {
		require.FailNowf(t, "alliance delegation err", "transaction error: code %d %s", txResp.Code, txResp.RawLog)
	}

	uatomBalanceAfterDelegation, err := persistenceChain.GetBalance(ctx, persistenceUserAddr, gaiaIBCDenom)
	require.NoError(t, err)
	require.Equal(t, math.ZeroInt(), uatomBalanceAfterDelegation,
		"uatom balance on Persistence chain must be 0 after delegation")

	validators = helpers.QueryAllValidators(t, ctx, persistenceChainNode)
	require.Len(t, validators, numVals, "validators returned must match count of validators created")

	totalTokensDelegated2 := math.NewInt(0)
	for _, val := range validators {
		totalTokensDelegated2 = totalTokensDelegated2.Add(val.Tokens)
	}
	t.Logf("[4] Total Tokens Delegated (including Alliance-minted): %s uxprt", totalTokensDelegated2.String())

	require.Equal(t, totalTokensDelegated.Mul(math.NewInt(2)), totalTokensDelegated2, "total uxprt delegated to validator must be x2 since weight of alliance asset is 1:1")

	// Wait 5 blocks
	err = testutil.WaitForBlocks(ctx, 5, persistenceChain)
	require.NoError(t, err)

	// Query rewards!
	allianceRewards, _, err := persistenceChainNode.ExecQuery(
		ctx, "alliance", "rewards",
		persistenceUserAddr, validators[0].OperatorAddress, gaiaIBCDenom,
	)
	require.NoError(t, err)
	t.Logf("[5] Alliance Staking rewards: %s", string(allianceRewards))

	uatomAlliance := helpers.QueryAlliance(t, ctx, persistenceChainNode, gaiaIBCDenom)
	require.Equal(t, allianceDelegationCoins.Amount, uatomAlliance.TotalTokens, "alliance staked total tokens must match initial delegation amount")

	uxprtBalanceBeforeRewardsClaim, err := persistenceChain.GetBalance(ctx, persistenceUserAddr, persistenceChain.Config().Denom)
	require.NoError(t, err)

	// Claim all rewards from a delegation
	txResp, err = cosmos.BroadcastTx(
		ctx,
		broadcaster,
		persistenceUser,
		&alliancetypes.MsgClaimDelegationRewards{
			DelegatorAddress: persistenceUserAddr,
			ValidatorAddress: validators[0].OperatorAddress,
			Denom:            gaiaIBCDenom,
		},
	)
	require.NoError(t, err, "error submitting alliance rewards claim tx")
	if txResp.Code != 0 {
		require.FailNowf(t, "alliance rewards claim err", "transaction error: code %d %s", txResp.Code, txResp.RawLog)
	}

	uxprtBalanceAfterRewardsClaim, err := persistenceChain.GetBalance(ctx, persistenceUserAddr, persistenceChain.Config().Denom)
	require.Greater(t, uxprtBalanceAfterRewardsClaim.Int64(), uxprtBalanceBeforeRewardsClaim.Int64(),
		"uxprt balance on Persistence chain must increase after rewards claim")
}
