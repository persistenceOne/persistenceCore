package hub

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"time"

	tendermintABCITypes "github.com/tendermint/tendermint/abci/types"
	tendermintCommon "github.com/tendermint/tendermint/libs/common"
	tendermintDB "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tendermintTypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

const applicationName = "CommitHubApplication"

var DefaultClientHome = os.ExpandEnv("$HOME/.hubClient")
var DefaultNodeHome = os.ExpandEnv("$HOME/.hubNode")

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	distribution.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	crisis.RegisterCodec(cdc)
	sdkTypes.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

type CommitHubApplication struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	keyMain          *sdkTypes.KVStoreKey
	keyAccount       *sdkTypes.KVStoreKey
	keyStaking       *sdkTypes.KVStoreKey
	tkeyStaking      *sdkTypes.TransientStoreKey
	keySlashing      *sdkTypes.KVStoreKey
	keyMint          *sdkTypes.KVStoreKey
	keyDistribution  *sdkTypes.KVStoreKey
	tkeyDistribution *sdkTypes.TransientStoreKey
	keyGov           *sdkTypes.KVStoreKey
	keyFeeCollection *sdkTypes.KVStoreKey
	keyParameter     *sdkTypes.KVStoreKey
	tkeyParameter    *sdkTypes.TransientStoreKey

	accountKeeper       auth.AccountKeeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	bankKeeper          bank.Keeper
	stakingKeeper       staking.Keeper
	slashingKeeper      slashing.Keeper
	mintKeeper          mint.Keeper
	distributionKeeper  distribution.Keeper
	govKeeper           gov.Keeper
	crisisKeeper        crisis.Keeper
	parameterKeeper     params.Keeper
}

func NewCommitHubApplication(logger log.Logger, db tendermintDB.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, baseAppOptions ...func(*baseapp.BaseApp)) *CommitHubApplication {

	cdc := MakeCodec()

	baseApp := baseapp.NewBaseApp(applicationName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	baseApp.SetCommitMultiStoreTracer(traceStore)

	var application = &CommitHubApplication{
		BaseApp:          baseApp,
		cdc:              cdc,
		invCheckPeriod:   invCheckPeriod,
		keyMain:          sdkTypes.NewKVStoreKey(baseapp.MainStoreKey),
		keyAccount:       sdkTypes.NewKVStoreKey(auth.StoreKey),
		keyStaking:       sdkTypes.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:      sdkTypes.NewTransientStoreKey(staking.TStoreKey),
		keyMint:          sdkTypes.NewKVStoreKey(mint.StoreKey),
		keyDistribution:  sdkTypes.NewKVStoreKey(distribution.StoreKey),
		tkeyDistribution: sdkTypes.NewTransientStoreKey(distribution.TStoreKey),
		keySlashing:      sdkTypes.NewKVStoreKey(slashing.StoreKey),
		keyGov:           sdkTypes.NewKVStoreKey(gov.StoreKey),
		keyFeeCollection: sdkTypes.NewKVStoreKey(auth.FeeStoreKey),
		keyParameter:     sdkTypes.NewKVStoreKey(params.StoreKey),
		tkeyParameter:    sdkTypes.NewTransientStoreKey(params.TStoreKey),
	}

	application.parameterKeeper = params.NewKeeper(
		application.cdc,
		application.keyParameter,
		application.tkeyParameter,
	)

	application.accountKeeper = auth.NewAccountKeeper(
		application.cdc,
		application.keyAccount,
		application.parameterKeeper.Subspace(auth.DefaultParamspace),
		auth.ProtoBaseAccount,
	)

	application.bankKeeper = bank.NewBaseKeeper(
		application.accountKeeper,
		application.parameterKeeper.Subspace(bank.DefaultParamspace),
		bank.DefaultCodespace,
	)
	application.feeCollectionKeeper = auth.NewFeeCollectionKeeper(
		application.cdc,
		application.keyFeeCollection,
	)
	stakingKeeper := staking.NewKeeper(
		application.cdc,
		application.keyStaking,
		application.tkeyStaking,
		application.bankKeeper,
		application.parameterKeeper.Subspace(staking.DefaultParamspace),
		staking.DefaultCodespace,
	)
	application.mintKeeper = mint.NewKeeper(application.cdc, application.keyMint,
		application.parameterKeeper.Subspace(mint.DefaultParamspace),
		&stakingKeeper,
		application.feeCollectionKeeper,
	)
	application.distributionKeeper = distribution.NewKeeper(
		application.cdc,
		application.keyDistribution,
		application.parameterKeeper.Subspace(distribution.DefaultParamspace),
		application.bankKeeper, &stakingKeeper, application.feeCollectionKeeper,
		distribution.DefaultCodespace,
	)
	application.slashingKeeper = slashing.NewKeeper(
		application.cdc,
		application.keySlashing,
		&stakingKeeper,
		application.parameterKeeper.Subspace(slashing.DefaultParamspace),
		slashing.DefaultCodespace,
	)
	application.govKeeper = gov.NewKeeper(
		application.cdc,
		application.keyGov,
		application.parameterKeeper,
		application.parameterKeeper.Subspace(gov.DefaultParamspace),
		application.bankKeeper, &stakingKeeper,
		gov.DefaultCodespace,
	)
	application.crisisKeeper = crisis.NewKeeper(
		application.parameterKeeper.Subspace(crisis.DefaultParamspace),
		application.distributionKeeper,
		application.bankKeeper,
		application.feeCollectionKeeper,
	)
	application.stakingKeeper = *stakingKeeper.SetHooks(
		NewStakingHooks(application.distributionKeeper.Hooks(), application.slashingKeeper.Hooks()),
	)

	bank.RegisterInvariants(&application.crisisKeeper, application.accountKeeper)
	distribution.RegisterInvariants(&application.crisisKeeper, application.distributionKeeper, application.stakingKeeper)
	staking.RegisterInvariants(&application.crisisKeeper, application.stakingKeeper, application.feeCollectionKeeper, application.distributionKeeper, application.accountKeeper)

	application.Router().
		AddRoute(bank.RouterKey, bank.NewHandler(application.bankKeeper)).
		AddRoute(staking.RouterKey, staking.NewHandler(application.stakingKeeper)).
		AddRoute(distribution.RouterKey, distribution.NewHandler(application.distributionKeeper)).
		AddRoute(slashing.RouterKey, slashing.NewHandler(application.slashingKeeper)).
		AddRoute(gov.RouterKey, gov.NewHandler(application.govKeeper)).
		AddRoute(crisis.RouterKey, crisis.NewHandler(application.crisisKeeper))

	application.QueryRouter().
		AddRoute(auth.QuerierRoute, auth.NewQuerier(application.accountKeeper)).
		AddRoute(distribution.QuerierRoute, distribution.NewQuerier(application.distributionKeeper)).
		AddRoute(gov.QuerierRoute, gov.NewQuerier(application.govKeeper)).
		AddRoute(slashing.QuerierRoute, slashing.NewQuerier(application.slashingKeeper, application.cdc)).
		AddRoute(staking.QuerierRoute, staking.NewQuerier(application.stakingKeeper, application.cdc)).
		AddRoute(mint.QuerierRoute, mint.NewQuerier(application.mintKeeper))

	application.MountStores(
		application.keyMain,
		application.keyAccount,
		application.keyStaking,
		application.keyMint,
		application.keyDistribution,
		application.keySlashing,
		application.keyGov,
		application.keyFeeCollection,
		application.keyParameter,
		application.tkeyParameter,
		application.tkeyStaking,
		application.tkeyDistribution,
	)

	application.SetInitChainer(application.initChainer)
	application.SetBeginBlocker(application.BeginBlocker)
	application.SetAnteHandler(auth.NewAnteHandler(application.accountKeeper, application.feeCollectionKeeper))
	application.SetEndBlocker(application.EndBlocker)

	if loadLatest {
		err := application.LoadLatestVersion(application.keyMain)
		if err != nil {
			tendermintCommon.Exit(err.Error())
		}
	}

	return application
}
func (commitHubApplication *CommitHubApplication) assertRuntimeInvariants() {
	context := commitHubApplication.NewContext(false, tendermintABCITypes.Header{Height: commitHubApplication.LastBlockHeight() + 1})
	commitHubApplication.assertRuntimeInvariantsOnContext(context)
}
func (commitHubApplication *CommitHubApplication) assertRuntimeInvariantsOnContext(ctx sdkTypes.Context) {
	start := time.Now()
	invarRoutes := commitHubApplication.crisisKeeper.Routes()
	for _, ir := range invarRoutes {
		if err := ir.Invar(ctx); err != nil {
			panic(fmt.Errorf("invariant broken: %s\n"+
				"\tCRITICAL please submit the following transaction:\n"+
				"\t\t gaiacli tx crisis invariant-broken %v %v", err, ir.ModuleName, ir.Route))
		}
	}
	end := time.Now()
	diff := end.Sub(start)
	commitHubApplication.BaseApp.Logger().With("module", "invariants").Info(
		"Asserted all invariants", "duration", diff, "height", commitHubApplication.LastBlockHeight())
}
func (commitHubApplication *CommitHubApplication) BeginBlocker(ctx sdkTypes.Context, req tendermintABCITypes.RequestBeginBlock) tendermintABCITypes.ResponseBeginBlock {
	mint.BeginBlocker(ctx, commitHubApplication.mintKeeper)
	distribution.BeginBlocker(ctx, req, commitHubApplication.distributionKeeper)
	tags := slashing.BeginBlocker(ctx, req, commitHubApplication.slashingKeeper)
	return tendermintABCITypes.ResponseBeginBlock{
		Tags: tags.ToKVPairs(),
	}
}
func (commitHubApplication *CommitHubApplication) EndBlocker(ctx sdkTypes.Context, req tendermintABCITypes.RequestEndBlock) tendermintABCITypes.ResponseEndBlock {
	tags := gov.EndBlocker(ctx, commitHubApplication.govKeeper)
	validatorUpdates, endBlockerTags := staking.EndBlocker(ctx, commitHubApplication.stakingKeeper)
	tags = append(tags, endBlockerTags...)

	if commitHubApplication.invCheckPeriod != 0 && ctx.BlockHeight()%int64(commitHubApplication.invCheckPeriod) == 0 {
		commitHubApplication.assertRuntimeInvariants()
	}

	return tendermintABCITypes.ResponseEndBlock{
		ValidatorUpdates: validatorUpdates,
		Tags:             tags,
	}
}
func (commitHubApplication *CommitHubApplication) initializeFromGenesisState(ctx sdkTypes.Context, genesisState GenesisState) []tendermintABCITypes.ValidatorUpdate {
	genesisState.Sanitize()

	for _, genesisStateAccount := range genesisState.Accounts {
		account := genesisStateAccount.ToAccount()
		account = commitHubApplication.accountKeeper.NewAccount(ctx, account)
		commitHubApplication.accountKeeper.SetAccount(ctx, account)
	}

	distribution.InitGenesis(ctx, commitHubApplication.distributionKeeper, genesisState.DistributionData)

	validators, err := staking.InitGenesis(ctx, commitHubApplication.stakingKeeper, genesisState.StakingData)
	if err != nil {
		panic(err)
	}

	auth.InitGenesis(ctx, commitHubApplication.accountKeeper, commitHubApplication.feeCollectionKeeper, genesisState.AuthData)
	bank.InitGenesis(ctx, commitHubApplication.bankKeeper, genesisState.BankData)
	slashing.InitGenesis(ctx, commitHubApplication.slashingKeeper, genesisState.SlashingData, genesisState.StakingData.Validators.ToSDKValidators())
	gov.InitGenesis(ctx, commitHubApplication.govKeeper, genesisState.GovernmentData)
	crisis.InitGenesis(ctx, commitHubApplication.crisisKeeper, genesisState.CrisisData)
	mint.InitGenesis(ctx, commitHubApplication.mintKeeper, genesisState.MintData)

	if err := ValidateGenesisState(genesisState); err != nil {
		panic(err)
	}

	if len(genesisState.GenesisTransactions) > 0 {
		for _, genesisTransactions := range genesisState.GenesisTransactions {
			var tx auth.StdTx
			err = commitHubApplication.cdc.UnmarshalJSON(genesisTransactions, &tx)
			if err != nil {
				panic(err)
			}
			bz := commitHubApplication.cdc.MustMarshalBinaryLengthPrefixed(tx)
			res := commitHubApplication.BaseApp.DeliverTx(bz)
			if !res.IsOK() {
				panic(res.Log)
			}
		}
		validators = commitHubApplication.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}
	return validators
}
func (commitHubApplication *CommitHubApplication) initChainer(ctx sdkTypes.Context, req tendermintABCITypes.RequestInitChain) tendermintABCITypes.ResponseInitChain {
	stateJSON := req.AppStateBytes

	var genesisState GenesisState
	err := commitHubApplication.cdc.UnmarshalJSON(stateJSON, &genesisState)
	if err != nil {
		panic(err)
	}

	validators := commitHubApplication.initializeFromGenesisState(ctx, genesisState)

	if len(req.Validators) > 0 {
		if len(req.Validators) != len(validators) {
			panic(fmt.Errorf("len(RequestInitChain.Validators) != len(validators) (%d != %d)",
				len(req.Validators), len(validators)))
		}
		sort.Sort(tendermintABCITypes.ValidatorUpdates(req.Validators))
		sort.Sort(tendermintABCITypes.ValidatorUpdates(validators))
		for i, val := range validators {
			if !val.Equal(req.Validators[i]) {
				panic(fmt.Errorf("validators[%d] != req.Validators[%d] ", i, i))
			}
		}
	}

	commitHubApplication.assertRuntimeInvariants()

	return tendermintABCITypes.ResponseInitChain{
		Validators: validators,
	}
}
func (commitHubApplication *CommitHubApplication) LoadHeight(height int64) error {
	return commitHubApplication.LoadVersion(height, commitHubApplication.keyMain)
}
func (commitHubApplication *CommitHubApplication) ExportApplicationStateAndValidators(forZeroHeight bool, jailWhiteList []string) (
	appState json.RawMessage, validators []tendermintTypes.GenesisValidator, err error) {

	ctx := commitHubApplication.NewContext(true, tendermintABCITypes.Header{Height: commitHubApplication.LastBlockHeight()})

	if forZeroHeight {
		commitHubApplication.prepareForZeroHeightGenesis(ctx, jailWhiteList)
	}

	accounts := []GenesisAccount{}
	appendAccount := func(acc auth.Account) (stop bool) {
		account := NewGenesisAccountI(acc)
		accounts = append(accounts, account)
		return false
	}
	commitHubApplication.accountKeeper.IterateAccounts(ctx, appendAccount)

	genState := NewGenesisState(
		accounts,
		auth.ExportGenesis(ctx, commitHubApplication.accountKeeper, commitHubApplication.feeCollectionKeeper),
		bank.ExportGenesis(ctx, commitHubApplication.bankKeeper),
		staking.ExportGenesis(ctx, commitHubApplication.stakingKeeper),
		mint.ExportGenesis(ctx, commitHubApplication.mintKeeper),
		distribution.ExportGenesis(ctx, commitHubApplication.distributionKeeper),
		gov.ExportGenesis(ctx, commitHubApplication.govKeeper),
		crisis.ExportGenesis(ctx, commitHubApplication.crisisKeeper),
		slashing.ExportGenesis(ctx, commitHubApplication.slashingKeeper),
	)
	appState, err = codec.MarshalJSONIndent(commitHubApplication.cdc, genState)
	if err != nil {
		return nil, nil, err
	}
	validators = staking.WriteValidators(ctx, commitHubApplication.stakingKeeper)
	return appState, validators, nil
}
func (commitHubApplication *CommitHubApplication) prepareForZeroHeightGenesis(ctx sdkTypes.Context, jailWhiteList []string) {
	applyWhiteList := false

	if len(jailWhiteList) > 0 {
		applyWhiteList = true
	}

	whiteListMap := make(map[string]bool)

	for _, addr := range jailWhiteList {
		_, err := sdkTypes.ValAddressFromBech32(addr)
		if err != nil {
			panic(err)
		}
		whiteListMap[addr] = true
	}

	commitHubApplication.assertRuntimeInvariantsOnContext(ctx)

	commitHubApplication.stakingKeeper.IterateValidators(ctx, func(_ int64, val sdkTypes.Validator) (stop bool) {
		_, _ = commitHubApplication.distributionKeeper.WithdrawValidatorCommission(ctx, val.GetOperator())
		return false
	})

	dels := commitHubApplication.stakingKeeper.GetAllDelegations(ctx)
	for _, delegation := range dels {
		_, _ = commitHubApplication.distributionKeeper.WithdrawDelegationRewards(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	}

	commitHubApplication.distributionKeeper.DeleteAllValidatorSlashEvents(ctx)

	commitHubApplication.distributionKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	commitHubApplication.stakingKeeper.IterateValidators(ctx, func(_ int64, val sdkTypes.Validator) (stop bool) {

		scraps := commitHubApplication.distributionKeeper.GetValidatorOutstandingRewards(ctx, val.GetOperator())
		feePool := commitHubApplication.distributionKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps)
		commitHubApplication.distributionKeeper.SetFeePool(ctx, feePool)

		commitHubApplication.distributionKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
		return false
	})

	for _, del := range dels {
		commitHubApplication.distributionKeeper.Hooks().BeforeDelegationCreated(ctx, del.DelegatorAddress, del.ValidatorAddress)
		commitHubApplication.distributionKeeper.Hooks().AfterDelegationModified(ctx, del.DelegatorAddress, del.ValidatorAddress)
	}

	ctx = ctx.WithBlockHeight(height)

	commitHubApplication.stakingKeeper.IterateRedelegations(ctx, func(_ int64, red staking.Redelegation) (stop bool) {
		for i := range red.Entries {
			red.Entries[i].CreationHeight = 0
		}
		commitHubApplication.stakingKeeper.SetRedelegation(ctx, red)
		return false
	})

	commitHubApplication.stakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd staking.UnbondingDelegation) (stop bool) {
		for i := range ubd.Entries {
			ubd.Entries[i].CreationHeight = 0
		}
		commitHubApplication.stakingKeeper.SetUnbondingDelegation(ctx, ubd)
		return false
	})

	store := ctx.KVStore(commitHubApplication.keyStaking)
	iter := sdkTypes.KVStoreReversePrefixIterator(store, staking.ValidatorsKey)
	counter := int16(0)

	var valConsAddrs []sdkTypes.ConsAddress
	for ; iter.Valid(); iter.Next() {
		addr := sdkTypes.ValAddress(iter.Key()[1:])
		validator, found := commitHubApplication.stakingKeeper.GetValidator(ctx, addr)
		if !found {
			panic("expected validator, not found")
		}

		validator.UnbondingHeight = 0
		valConsAddrs = append(valConsAddrs, validator.ConsAddress())
		if applyWhiteList && !whiteListMap[addr.String()] {
			validator.Jailed = true
		}

		commitHubApplication.stakingKeeper.SetValidator(ctx, validator)
		counter++
	}

	iter.Close()

	_ = commitHubApplication.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)

	commitHubApplication.slashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(addr sdkTypes.ConsAddress, info slashing.ValidatorSigningInfo) (stop bool) {
			info.StartHeight = 0
			commitHubApplication.slashingKeeper.SetValidatorSigningInfo(ctx, addr, info)
			return false
		},
	)
}

var _ sdkTypes.StakingHooks = StakingHooks{}

type StakingHooks struct {
	distributionHooks distribution.Hooks
	stakingHooks      slashing.Hooks
}

func NewStakingHooks(distributionHooks distribution.Hooks, stakingHooks slashing.Hooks) StakingHooks {
	return StakingHooks{distributionHooks, stakingHooks}
}
func (stakingHooks StakingHooks) AfterValidatorCreated(ctx sdkTypes.Context, valAddr sdkTypes.ValAddress) {
	stakingHooks.distributionHooks.AfterValidatorCreated(ctx, valAddr)
	stakingHooks.stakingHooks.AfterValidatorCreated(ctx, valAddr)
}
func (stakingHooks StakingHooks) BeforeValidatorModified(ctx sdkTypes.Context, valAddr sdkTypes.ValAddress) {
	stakingHooks.distributionHooks.BeforeValidatorModified(ctx, valAddr)
	stakingHooks.stakingHooks.BeforeValidatorModified(ctx, valAddr)
}
func (stakingHooks StakingHooks) AfterValidatorRemoved(ctx sdkTypes.Context, consAddr sdkTypes.ConsAddress, valAddr sdkTypes.ValAddress) {
	stakingHooks.distributionHooks.AfterValidatorRemoved(ctx, consAddr, valAddr)
	stakingHooks.stakingHooks.AfterValidatorRemoved(ctx, consAddr, valAddr)
}
func (stakingHooks StakingHooks) AfterValidatorBonded(ctx sdkTypes.Context, consAddr sdkTypes.ConsAddress, valAddr sdkTypes.ValAddress) {
	stakingHooks.distributionHooks.AfterValidatorBonded(ctx, consAddr, valAddr)
	stakingHooks.stakingHooks.AfterValidatorBonded(ctx, consAddr, valAddr)
}
func (stakingHooks StakingHooks) AfterValidatorBeginUnbonding(ctx sdkTypes.Context, consAddr sdkTypes.ConsAddress, valAddr sdkTypes.ValAddress) {
	stakingHooks.distributionHooks.AfterValidatorBeginUnbonding(ctx, consAddr, valAddr)
	stakingHooks.stakingHooks.AfterValidatorBeginUnbonding(ctx, consAddr, valAddr)
}
func (stakingHooks StakingHooks) BeforeDelegationCreated(ctx sdkTypes.Context, delAddr sdkTypes.AccAddress, valAddr sdkTypes.ValAddress) {
	stakingHooks.distributionHooks.BeforeDelegationCreated(ctx, delAddr, valAddr)
	stakingHooks.stakingHooks.BeforeDelegationCreated(ctx, delAddr, valAddr)
}
func (stakingHooks StakingHooks) BeforeDelegationSharesModified(ctx sdkTypes.Context, delAddr sdkTypes.AccAddress, valAddr sdkTypes.ValAddress) {
	stakingHooks.distributionHooks.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
	stakingHooks.stakingHooks.BeforeDelegationSharesModified(ctx, delAddr, valAddr)
}
func (stakingHooks StakingHooks) BeforeDelegationRemoved(ctx sdkTypes.Context, delAddr sdkTypes.AccAddress, valAddr sdkTypes.ValAddress) {
	stakingHooks.distributionHooks.BeforeDelegationRemoved(ctx, delAddr, valAddr)
	stakingHooks.stakingHooks.BeforeDelegationRemoved(ctx, delAddr, valAddr)
}
func (stakingHooks StakingHooks) AfterDelegationModified(ctx sdkTypes.Context, delAddr sdkTypes.AccAddress, valAddr sdkTypes.ValAddress) {
	stakingHooks.distributionHooks.AfterDelegationModified(ctx, delAddr, valAddr)
	stakingHooks.stakingHooks.AfterDelegationModified(ctx, delAddr, valAddr)
}
func (stakingHooks StakingHooks) BeforeValidatorSlashed(ctx sdkTypes.Context, valAddr sdkTypes.ValAddress, fraction sdkTypes.Dec) {
	stakingHooks.distributionHooks.BeforeValidatorSlashed(ctx, valAddr, fraction)
	stakingHooks.stakingHooks.BeforeValidatorSlashed(ctx, valAddr, fraction)
}
