package application

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	ibcclient "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
	port "github.com/cosmos/cosmos-sdk/x/ibc/05-port"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/persistenceOne/persistenceSDK/modules/assets"
	"github.com/persistenceOne/persistenceSDK/modules/classifications"
	"github.com/persistenceOne/persistenceSDK/modules/exchanges"
	"github.com/persistenceOne/persistenceSDK/modules/exchanges/auxiliaries/custody"
	"github.com/persistenceOne/persistenceSDK/modules/exchanges/auxiliaries/reverse"
	"github.com/persistenceOne/persistenceSDK/modules/exchanges/auxiliaries/swap"
	"github.com/persistenceOne/persistenceSDK/modules/identities"
	"github.com/persistenceOne/persistenceSDK/modules/identities/auxiliaries/verify"
	"github.com/persistenceOne/persistenceSDK/modules/metas"
	"github.com/persistenceOne/persistenceSDK/modules/metas/auxiliaries/initialize"
	"github.com/persistenceOne/persistenceSDK/modules/orders"
	"github.com/persistenceOne/persistenceSDK/modules/splits"
	"github.com/persistenceOne/persistenceSDK/modules/splits/auxiliaries/burn"
	auxiliariesMint "github.com/persistenceOne/persistenceSDK/modules/splits/auxiliaries/mint"
	"github.com/persistenceOne/persistenceSDK/schema"
	"github.com/spf13/viper"
	"io"
	"os"
	"path/filepath"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tendermintOS "github.com/tendermint/tendermint/libs/os"
	tendermintTypes "github.com/tendermint/tendermint/types"
	tendermintDB "github.com/tendermint/tm-db"
	"honnef.co/go/tools/version"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsClient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	upgradeClient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
)

const applicationName = "AssetMantle"

var DefaultClientHome = os.ExpandEnv("$HOME/.assetClient")
var DefaultNodeHome = os.ExpandEnv("$HOME/.assetNode")
var moduleAccountPermissions = map[string][]string{
	auth.FeeCollectorName:           nil,
	distribution.ModuleName:         nil,
	mint.ModuleName:                 {auth.Minter},
	staking.BondedPoolName:          {auth.Burner, auth.Staking},
	staking.NotBondedPoolName:       {auth.Burner, auth.Staking},
	gov.ModuleName:                  {auth.Burner},
	transfer.GetModuleAccountName(): {auth.Minter, auth.Burner},
	splits.Module.Name():            nil,
}
var tokenReceiveAllowedModules = map[string]bool{
	distribution.ModuleName: true,
}
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.NewAppModuleBasic(paramsClient.ProposalHandler, distribution.ProposalHandler, upgradeClient.ProposalHandler),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	wasm.AppModuleBasic{},
	slashing.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},

	assets.Module,
	classifications.Module,
	exchanges.Module,
	identities.Module,
	metas.Module,
	orders.Module,
	splits.Module,
)

type GenesisState map[string]json.RawMessage

func MakeCodecs() (*std.Codec, *codec.Codec) {
	cdc := std.MakeCodec(ModuleBasics)
	schema.RegisterCodec(cdc)
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	appCodec := std.NewAppCodec(cdc, interfaceRegistry)

	sdkTypes.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterInterfaceModules(interfaceRegistry)

	return appCodec, cdc
}

type Application struct {
	*baseapp.BaseApp
	codec *codec.Codec

	invCheckPeriod uint

	keys               map[string]*sdkTypes.KVStoreKey
	transientStoreKeys map[string]*sdkTypes.TransientStoreKey
	memoryStoreKeys    map[string]*sdkTypes.MemoryStoreKey

	subspaces map[string]params.Subspace

	accountKeeper      auth.AccountKeeper
	bankKeeper         bank.Keeper
	capabilityKeeper   *capability.Keeper
	stakingKeeper      staking.Keeper
	slashingKeeper     slashing.Keeper
	mintKeeper         mint.Keeper
	distributionKeeper distribution.Keeper
	govKeeper          gov.Keeper
	crisisKeeper       crisis.Keeper
	upgradeKeeper      upgrade.Keeper
	paramsKeeper       params.Keeper
	ibcKeeper          *ibc.Keeper
	evidenceKeeper     evidence.Keeper
	transferKeeper     transfer.Keeper

	scopedIBCKeeper      capability.ScopedKeeper
	scopedTransferKeeper capability.ScopedKeeper
	wasmKeeper           wasm.Keeper

	moduleManager *module.Manager

	simulationManager *module.SimulationManager
}

// WasmWrapper allows us to use namespacing in the config file
// This is only used for parsing in the app, x/wasm expects WasmConfig
type WasmWrapper struct {
	Wasm wasm.WasmConfig `mapstructure:"wasm"`
}

func NewApplication(
	logger log.Logger,
	db tendermintDB.DB,
	traceStore io.Writer,
	loadLatest bool,
	invCheckPeriod uint,
	skipUpgradeHeights map[int64]bool,
	home string,
	baseAppOptions ...func(*baseapp.BaseApp),
) *Application {

	appCodec, Codec := MakeCodecs()
	baseApp := baseapp.NewBaseApp(
		applicationName,
		logger,
		db,
		auth.DefaultTxDecoder(Codec),
		baseAppOptions...,
	)
	baseApp.SetCommitMultiStoreTracer(traceStore)
	baseApp.SetAppVersion(version.Version)

	keys := sdkTypes.NewKVStoreKeys(
		auth.StoreKey,
		bank.ModuleName,
		staking.StoreKey,
		mint.StoreKey,
		distribution.StoreKey,
		slashing.StoreKey,
		gov.StoreKey,
		params.StoreKey,
		ibc.StoreKey,
		upgrade.StoreKey,
		evidence.StoreKey,
		transfer.StoreKey,
		capability.StoreKey,
		wasm.StoreKey,
	)
	keys[assets.Module.Name()] = assets.Module.GetKVStoreKey()
	keys[classifications.Module.Name()] = classifications.Module.GetKVStoreKey()
	keys[exchanges.Module.Name()] = exchanges.Module.GetKVStoreKey()
	keys[identities.Module.Name()] = identities.Module.GetKVStoreKey()
	keys[metas.Module.Name()] = metas.Module.GetKVStoreKey()
	keys[orders.Module.Name()] = orders.Module.GetKVStoreKey()
	keys[splits.Module.Name()] = splits.Module.GetKVStoreKey()

	transientStoreKeys := sdkTypes.NewTransientStoreKeys(params.TStoreKey)
	memoryStoreKeys := sdkTypes.NewMemoryStoreKeys(capability.MemStoreKey)

	var application = &Application{
		BaseApp: baseApp,
		codec:   Codec,

		invCheckPeriod: invCheckPeriod,

		keys:               keys,
		transientStoreKeys: transientStoreKeys,
		memoryStoreKeys:    memoryStoreKeys,

		subspaces: make(map[string]params.Subspace),
	}

	application.paramsKeeper = params.NewKeeper(
		appCodec,
		keys[params.StoreKey],
		transientStoreKeys[params.TStoreKey],
	)
	application.subspaces[auth.ModuleName] = application.paramsKeeper.Subspace(auth.DefaultParamspace)
	application.subspaces[bank.ModuleName] = application.paramsKeeper.Subspace(bank.DefaultParamspace)
	application.subspaces[staking.ModuleName] = application.paramsKeeper.Subspace(staking.DefaultParamspace)
	application.subspaces[mint.ModuleName] = application.paramsKeeper.Subspace(mint.DefaultParamspace)
	application.subspaces[distribution.ModuleName] = application.paramsKeeper.Subspace(distribution.DefaultParamspace)
	application.subspaces[slashing.ModuleName] = application.paramsKeeper.Subspace(slashing.DefaultParamspace)
	application.subspaces[gov.ModuleName] = application.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	application.subspaces[crisis.ModuleName] = application.paramsKeeper.Subspace(crisis.DefaultParamspace)

	application.subspaces[assets.Module.Name()] = application.paramsKeeper.Subspace(assets.Module.GetDefaultParamspace())
	application.subspaces[classifications.Module.Name()] = application.paramsKeeper.Subspace(classifications.Module.GetDefaultParamspace())
	application.subspaces[exchanges.Module.Name()] = application.paramsKeeper.Subspace(exchanges.Module.GetDefaultParamspace())
	application.subspaces[identities.Module.Name()] = application.paramsKeeper.Subspace(identities.Module.GetDefaultParamspace())
	application.subspaces[metas.Module.Name()] = application.paramsKeeper.Subspace(metas.Module.GetDefaultParamspace())
	application.subspaces[orders.Module.Name()] = application.paramsKeeper.Subspace(orders.Module.GetDefaultParamspace())
	application.subspaces[splits.Module.Name()] = application.paramsKeeper.Subspace(splits.Module.GetDefaultParamspace())

	baseApp.SetParamStore(application.paramsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(std.ConsensusParamsKeyTable()))

	application.capabilityKeeper = capability.NewKeeper(appCodec, keys[capability.StoreKey], memoryStoreKeys[capability.MemStoreKey])
	scopedIBCKeeper := application.capabilityKeeper.ScopeToModule(ibc.ModuleName)
	scopedTransferKeeper := application.capabilityKeeper.ScopeToModule(transfer.ModuleName)

	application.accountKeeper = auth.NewAccountKeeper(
		appCodec,
		keys[auth.StoreKey],
		application.subspaces[auth.ModuleName],
		auth.ProtoBaseAccount,
		moduleAccountPermissions,
	)

	application.bankKeeper = bank.NewBaseKeeper(
		appCodec,
		keys[bank.StoreKey],
		application.accountKeeper,
		application.subspaces[bank.ModuleName],
		application.BlacklistedAccAddrs(),
	)

	stakingKeeper := staking.NewKeeper(
		appCodec,
		keys[staking.StoreKey],
		application.accountKeeper,
		application.bankKeeper,
		application.subspaces[staking.ModuleName],
	)
	application.mintKeeper = mint.NewKeeper(
		appCodec,
		keys[mint.StoreKey],
		application.subspaces[mint.ModuleName],
		&stakingKeeper,
		application.accountKeeper,
		application.bankKeeper,
		auth.FeeCollectorName,
	)
	application.distributionKeeper = distribution.NewKeeper(
		appCodec,
		keys[distribution.StoreKey],
		application.subspaces[distribution.ModuleName],
		application.accountKeeper,
		application.bankKeeper,
		&stakingKeeper,
		auth.FeeCollectorName,
		application.ModuleAccountAddress(),
	)
	application.slashingKeeper = slashing.NewKeeper(
		appCodec,
		keys[slashing.StoreKey],
		&stakingKeeper,
		application.subspaces[slashing.ModuleName],
	)
	application.crisisKeeper = crisis.NewKeeper(
		application.subspaces[crisis.ModuleName],
		invCheckPeriod,
		application.bankKeeper,
		auth.FeeCollectorName,
	)
	application.upgradeKeeper = upgrade.NewKeeper(
		skipUpgradeHeights,
		keys[upgrade.StoreKey],
		appCodec,
		home,
	)

	govRouter := gov.NewRouter()
	govRouter.AddRoute(
		gov.RouterKey,
		gov.ProposalHandler,
	).AddRoute(
		paramproposal.RouterKey,
		params.NewParamChangeProposalHandler(application.paramsKeeper),
	).AddRoute(
		distribution.RouterKey,
		distribution.NewCommunityPoolSpendProposalHandler(application.distributionKeeper),
	).AddRoute(
		upgrade.RouterKey,
		upgrade.NewSoftwareUpgradeProposalHandler(application.upgradeKeeper),
	)
	application.govKeeper = gov.NewKeeper(
		appCodec,
		keys[gov.StoreKey],
		application.subspaces[gov.ModuleName],
		application.accountKeeper,
		application.bankKeeper,
		&stakingKeeper,
		govRouter,
	)

	application.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(application.distributionKeeper.Hooks(), application.slashingKeeper.Hooks()),
	)

	application.ibcKeeper = ibc.NewKeeper(
		application.codec, appCodec, keys[ibc.StoreKey], application.stakingKeeper, scopedIBCKeeper,
	)

	application.transferKeeper = transfer.NewKeeper(
		appCodec, keys[transfer.StoreKey],
		application.ibcKeeper.ChannelKeeper, &application.ibcKeeper.PortKeeper,
		application.accountKeeper, application.bankKeeper,
		scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(application.transferKeeper)

	ibcRouter := port.NewRouter()
	ibcRouter.AddRoute(transfer.ModuleName, transferModule)
	application.ibcKeeper.SetRouter(ibcRouter)

	evidenceKeeper := evidence.NewKeeper(
		appCodec,
		keys[evidence.StoreKey],
		&stakingKeeper,
		application.slashingKeeper,
	)
	evidenceRouter := evidence.NewRouter().AddRoute(ibcclient.RouterKey, ibcclient.HandlerClientMisbehaviour(application.ibcKeeper.ClientKeeper))
	evidenceKeeper.SetRouter(evidenceRouter)
	application.evidenceKeeper = *evidenceKeeper

	classifications.Module.InitializeKeepers()
	identities.Module.InitializeKeepers()
	metas.Module.InitializeKeepers()
	splits.Module.InitializeKeepers(
		application.bankKeeper,
		identities.Module.GetAuxiliary(verify.AuxiliaryName),
	)
	assets.Module.InitializeKeepers(
		identities.Module.GetAuxiliary(verify.AuxiliaryName),
		metas.Module.GetAuxiliary(initialize.AuxiliaryName),
		splits.Module.GetAuxiliary(auxiliariesMint.AuxiliaryName),
		splits.Module.GetAuxiliary(burn.AuxiliaryName),
	)
	exchanges.Module.InitializeKeepers(
		splits.Module.GetAuxiliary(auxiliariesMint.AuxiliaryName),
		splits.Module.GetAuxiliary(burn.AuxiliaryName),
	)
	orders.Module.InitializeKeepers(
		application.bankKeeper,
		exchanges.Module.GetAuxiliary(swap.AuxiliaryName),
		exchanges.Module.GetAuxiliary(custody.AuxiliaryName),
		exchanges.Module.GetAuxiliary(reverse.AuxiliaryName),
		identities.Module.GetAuxiliary(verify.AuxiliaryName),
	)

	// just re-use the full router - do we want to limit this more?
	var wasmRouter = baseApp.Router()
	wasmDir := filepath.Join(home, wasm.ModuleName)

	wasmWrap := WasmWrapper{Wasm: wasm.DefaultWasmConfig()}
	err := viper.Unmarshal(&wasmWrap)
	if err != nil {
		panic("error while reading wasm config: " + err.Error())
	}
	wasmConfig := wasmWrap.Wasm

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	supportedFeatures := "staking"
	application.wasmKeeper = wasm.NewKeeper(appCodec, keys[wasm.StoreKey],
		application.accountKeeper, application.bankKeeper, application.stakingKeeper,
		wasmRouter, wasmDir, wasmConfig, supportedFeatures, WasmCustomMessageEncoder(Codec), nil)

	application.moduleManager = module.NewManager(
		genutil.NewAppModule(application.accountKeeper, application.stakingKeeper, application.BaseApp.DeliverTx),
		auth.NewAppModule(appCodec, application.accountKeeper),
		bank.NewAppModule(appCodec, application.bankKeeper, application.accountKeeper),
		capability.NewAppModule(appCodec, *application.capabilityKeeper),
		crisis.NewAppModule(&application.crisisKeeper),
		gov.NewAppModule(appCodec, application.govKeeper, application.accountKeeper, application.bankKeeper),
		mint.NewAppModule(appCodec, application.mintKeeper, application.accountKeeper),
		slashing.NewAppModule(appCodec, application.slashingKeeper, application.accountKeeper, application.bankKeeper, application.stakingKeeper),
		distribution.NewAppModule(appCodec, application.distributionKeeper, application.accountKeeper, application.bankKeeper, application.stakingKeeper),
		staking.NewAppModule(appCodec, application.stakingKeeper, application.accountKeeper, application.bankKeeper),
		upgrade.NewAppModule(application.upgradeKeeper),
		wasm.NewAppModule(application.wasmKeeper),
		evidence.NewAppModule(application.evidenceKeeper),
		ibc.NewAppModule(application.ibcKeeper),
		params.NewAppModule(application.paramsKeeper),
		transferModule,

		assets.Module,
		classifications.Module,
		exchanges.Module,
		identities.Module,
		metas.Module,
		orders.Module,
		splits.Module,
	)

	application.moduleManager.SetOrderBeginBlockers(
		upgrade.ModuleName,
		mint.ModuleName,
		distribution.ModuleName,
		slashing.ModuleName,
	)
	application.moduleManager.SetOrderEndBlockers(
		crisis.ModuleName,
		gov.ModuleName,
		staking.ModuleName,
	)
	application.moduleManager.SetOrderInitGenesis(
		capability.ModuleName,
		auth.ModuleName,
		distribution.ModuleName,
		staking.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		gov.ModuleName,
		mint.ModuleName,
		crisis.ModuleName,
		ibc.ModuleName,
		genutil.ModuleName,
		evidence.ModuleName,
		transfer.ModuleName,
		wasm.ModuleName,
		assets.Module.Name(),
		classifications.Module.Name(),
		exchanges.Module.Name(),
		identities.Module.Name(),
		metas.Module.Name(),
		orders.Module.Name(),
		splits.Module.Name(),
	)
	application.moduleManager.RegisterInvariants(&application.crisisKeeper)
	application.moduleManager.RegisterRoutes(application.Router(), application.QueryRouter())

	//TODO add persistenceSDK modules to simulation
	application.simulationManager = module.NewSimulationManager(
		auth.NewAppModule(appCodec, application.accountKeeper),
		bank.NewAppModule(appCodec, application.bankKeeper, application.accountKeeper),
		capability.NewAppModule(appCodec, *application.capabilityKeeper),
		gov.NewAppModule(appCodec, application.govKeeper, application.accountKeeper, application.bankKeeper),
		mint.NewAppModule(appCodec, application.mintKeeper, application.accountKeeper),
		staking.NewAppModule(appCodec, application.stakingKeeper, application.accountKeeper, application.bankKeeper),
		distribution.NewAppModule(appCodec, application.distributionKeeper, application.accountKeeper, application.bankKeeper, application.stakingKeeper),
		slashing.NewAppModule(appCodec, application.slashingKeeper, application.accountKeeper, application.bankKeeper, application.stakingKeeper),
		params.NewAppModule(application.paramsKeeper),
		evidence.NewAppModule(application.evidenceKeeper),
	)

	application.simulationManager.RegisterStoreDecoders()

	application.MountKVStores(keys)
	application.MountTransientStores(transientStoreKeys)
	application.MountMemoryStores(memoryStoreKeys)

	application.SetInitChainer(application.InitChainer)
	application.SetBeginBlocker(application.BeginBlocker)
	application.SetAnteHandler(auth.NewAnteHandler(application.accountKeeper, application.bankKeeper, *application.ibcKeeper, ante.DefaultSigVerificationGasConsumer))
	application.SetEndBlocker(application.EndBlocker)

	if loadLatest {
		err := application.LoadLatestVersion()
		if err != nil {
			tendermintOS.Exit(err.Error())
		}
	}

	ctx := application.BaseApp.NewUncachedContext(true, abciTypes.Header{})
	application.capabilityKeeper.InitializeAndSeal(ctx)

	application.scopedIBCKeeper = scopedIBCKeeper
	application.scopedTransferKeeper = scopedTransferKeeper

	return application
}
func (application *Application) BeginBlocker(ctx sdkTypes.Context, req abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	return application.moduleManager.BeginBlock(ctx, req)
}
func (application *Application) EndBlocker(ctx sdkTypes.Context, req abciTypes.RequestEndBlock) abciTypes.ResponseEndBlock {
	return application.moduleManager.EndBlock(ctx, req)
}
func (application *Application) InitChainer(ctx sdkTypes.Context, req abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
	var genesisState GenesisState
	application.codec.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return application.moduleManager.InitGenesis(ctx, application.codec, genesisState)
}
func (application *Application) LoadHeight(height int64) error {
	return application.LoadVersion(height)
}
func (application *Application) ModuleAccountAddress() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range moduleAccountPermissions {
		modAccAddrs[auth.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}
func (application *Application) ExportApplicationStateAndValidators(forZeroHeight bool, jailWhiteList []string,
) (applicationState json.RawMessage, validators []tendermintTypes.GenesisValidator, cp *abciTypes.ConsensusParams, err error) {
	ctx := application.NewContext(true, abciTypes.Header{Height: application.LastBlockHeight()})

	if forZeroHeight {
		application.prepareForZeroHeightGenesis(ctx, jailWhiteList)
	}

	genesisState := application.moduleManager.ExportGenesis(ctx, application.codec)
	applicationState, err = codec.MarshalJSONIndent(application.codec, genesisState)
	if err != nil {
		return nil, nil, nil, err
	}
	validators = staking.WriteValidators(ctx, application.stakingKeeper)
	return applicationState, validators, application.BaseApp.GetConsensusParams(ctx), nil
}
func (application *Application) BlacklistedAccAddrs() map[string]bool {
	blacklistedAddresses := make(map[string]bool)
	for account := range moduleAccountPermissions {
		blacklistedAddresses[auth.NewModuleAddress(account).String()] = !tokenReceiveAllowedModules[account]
	}

	return blacklistedAddresses
}

func (application *Application) prepareForZeroHeightGenesis(ctx sdkTypes.Context, jailWhiteList []string) {
	applyWhiteList := false

	if len(jailWhiteList) > 0 {
		applyWhiteList = true
	}

	whiteListMap := make(map[string]bool)

	for _, address := range jailWhiteList {
		_, err := sdkTypes.ValAddressFromBech32(address)
		if err != nil {
			//log.Fatal(err) //todo
		}
		whiteListMap[address] = true
	}

	application.crisisKeeper.AssertInvariants(ctx)

	application.stakingKeeper.IterateValidators(ctx, func(_ int64, val staking.ValidatorI) (stop bool) {
		_, _ = application.distributionKeeper.WithdrawValidatorCommission(ctx, val.GetOperator())
		return false
	})

	delegations := application.stakingKeeper.GetAllDelegations(ctx)
	for _, delegation := range delegations {
		_, _ = application.distributionKeeper.WithdrawDelegationRewards(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	}

	application.distributionKeeper.DeleteAllValidatorSlashEvents(ctx)

	application.distributionKeeper.DeleteAllValidatorHistoricalRewards(ctx)

	height := ctx.BlockHeight()
	ctx = ctx.WithBlockHeight(0)

	application.stakingKeeper.IterateValidators(ctx, func(_ int64, val staking.ValidatorI) (stop bool) {

		scraps := application.distributionKeeper.GetValidatorOutstandingRewards(ctx, val.GetOperator()).Rewards
		feePool := application.distributionKeeper.GetFeePool(ctx)
		feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
		application.distributionKeeper.SetFeePool(ctx, feePool)

		application.distributionKeeper.Hooks().AfterValidatorCreated(ctx, val.GetOperator())
		return false
	})

	for _, delegation := range delegations {
		application.distributionKeeper.Hooks().BeforeDelegationCreated(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
		application.distributionKeeper.Hooks().AfterDelegationModified(ctx, delegation.DelegatorAddress, delegation.ValidatorAddress)
	}

	ctx = ctx.WithBlockHeight(height)

	application.stakingKeeper.IterateRedelegations(ctx, func(_ int64, redelegation staking.Redelegation) (stop bool) {
		for i := range redelegation.Entries {
			redelegation.Entries[i].CreationHeight = 0
		}
		application.stakingKeeper.SetRedelegation(ctx, redelegation)
		return false
	})

	application.stakingKeeper.IterateUnbondingDelegations(ctx, func(_ int64, unbondingDelegation staking.UnbondingDelegation) (stop bool) {
		for i := range unbondingDelegation.Entries {
			unbondingDelegation.Entries[i].CreationHeight = 0
		}
		application.stakingKeeper.SetUnbondingDelegation(ctx, unbondingDelegation)
		return false
	})

	store := ctx.KVStore(application.keys[staking.StoreKey])
	kvStoreReversePrefixIterator := sdkTypes.KVStoreReversePrefixIterator(store, staking.ValidatorsKey)
	counter := int16(0)

	for ; kvStoreReversePrefixIterator.Valid(); kvStoreReversePrefixIterator.Next() {
		addr := sdkTypes.ValAddress(kvStoreReversePrefixIterator.Key()[1:])
		validator, found := application.stakingKeeper.GetValidator(ctx, addr)
		if !found {
			panic("Validator not found!")
		}

		validator.UnbondingHeight = 0

		if applyWhiteList && !whiteListMap[addr.String()] {
			validator.Jailed = true
		}

		application.stakingKeeper.SetValidator(ctx, validator)
		counter++
	}

	kvStoreReversePrefixIterator.Close()

	_ = application.stakingKeeper.ApplyAndReturnValidatorSetUpdates(ctx)

	application.slashingKeeper.IterateValidatorSigningInfos(
		ctx,
		func(validatorConsAddress sdkTypes.ConsAddress, validatorSigningInfo slashing.ValidatorSigningInfo) (stop bool) {
			validatorSigningInfo.StartHeight = 0
			application.slashingKeeper.SetValidatorSigningInfo(ctx, validatorConsAddress, validatorSigningInfo)
			return false
		},
	)
}
