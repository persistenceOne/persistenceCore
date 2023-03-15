package keepers

import (
	"fmt"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/store/streaming"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icacontroller "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/host/types"
	ibcfee "github.com/cosmos/ibc-go/v6/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/types"
	"github.com/cosmos/ibc-go/v6/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v6/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibccoreclient "github.com/cosmos/ibc-go/v6/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	ibctypes "github.com/cosmos/ibc-go/v6/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v6/modules/core/keeper"
	epochskeeper "github.com/persistenceOne/persistence-sdk/v2/x/epochs/keeper"
	epochstypes "github.com/persistenceOne/persistence-sdk/v2/x/epochs/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/halving"
	"github.com/persistenceOne/persistence-sdk/v2/x/ibchooker"
	ibchookerkeeper "github.com/persistenceOne/persistence-sdk/v2/x/ibchooker/keeper"
	ibchookertypes "github.com/persistenceOne/persistence-sdk/v2/x/ibchooker/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/interchainquery"
	interchainquerykeeper "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/keeper"
	interchainquerytypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/lsnative/distribution"
	distributionkeeper "github.com/persistenceOne/persistence-sdk/v2/x/lsnative/distribution/keeper"
	distributiontypes "github.com/persistenceOne/persistence-sdk/v2/x/lsnative/distribution/types"
	slashingkeeper "github.com/persistenceOne/persistence-sdk/v2/x/lsnative/slashing/keeper"
	slashingtypes "github.com/persistenceOne/persistence-sdk/v2/x/lsnative/slashing/types"
	stakingkeeper "github.com/persistenceOne/persistence-sdk/v2/x/lsnative/staking/keeper"
	stakingtypes "github.com/persistenceOne/persistence-sdk/v2/x/lsnative/staking/types"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos"
	lscosmoskeeper "github.com/persistenceOne/pstake-native/v2/x/lscosmos/keeper"
	lscosmostypes "github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
	tmos "github.com/tendermint/tendermint/libs/os"

	oraclekeeper "github.com/persistenceOne/persistence-sdk/v2/x/oracle/keeper"
	oracletypes "github.com/persistenceOne/persistence-sdk/v2/x/oracle/types"
	"github.com/persistenceOne/persistenceCore/v7/wasmbindings"

	// unnamed import of statik for swagger UI support
	_ "github.com/cosmos/cosmos-sdk/client/docs/statik"
)

type AppKeepers struct {
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	AccountKeeper         *authkeeper.AccountKeeper
	BankKeeper            *bankkeeper.BaseKeeper
	CapabilityKeeper      *capabilitykeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        *slashingkeeper.Keeper
	MintKeeper            *mintkeeper.Keeper
	DistributionKeeper    *distributionkeeper.Keeper
	GovKeeper             *govkeeper.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	CrisisKeeper          *crisiskeeper.Keeper
	ParamsKeeper          *paramskeeper.Keeper
	IBCKeeper             *ibckeeper.Keeper
	IBCFeeKeeper          *ibcfeekeeper.Keeper
	ICAHostKeeper         *icahostkeeper.Keeper
	EvidenceKeeper        *evidencekeeper.Keeper
	TransferKeeper        *ibctransferkeeper.Keeper
	FeegrantKeeper        *feegrantkeeper.Keeper
	AuthzKeeper           *authzkeeper.Keeper
	HalvingKeeper         *halving.Keeper
	WasmKeeper            *wasm.Keeper
	EpochsKeeper          *epochskeeper.Keeper
	OracleKeeper          *oraclekeeper.Keeper
	ICAControllerKeeper   *icacontrollerkeeper.Keeper
	LSCosmosKeeper        *lscosmoskeeper.Keeper
	InterchainQueryKeeper *interchainquerykeeper.Keeper
	TransferHooksKeeper   *ibchookerkeeper.Keeper

	// Modules
	TransferModule             transfer.AppModule
	TransferIBCModule          ibctypes.IBCModule
	InterchainQueryModule      interchainquery.AppModule
	LSCosmosModule             lscosmos.AppModule
	IBCTransferHooksMiddleware ibchooker.AppModule

	// make scoped keepers public for test purposes
	ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper      capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper
	ScopedLSCosmosKeeper      capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper          capabilitykeeper.ScopedKeeper
}

func NewAppKeeper(
	appCodec codec.Codec,
	bApp *baseapp.BaseApp,
	legacyAmino *codec.LegacyAmino,
	maccPerms map[string][]string,
	sendCoinBlockedAddrs map[string]bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	appOpts servertypes.AppOptions,
	wasmDir string,
	wasmEnabledProposals []wasm.ProposalType,
	wasmOpts []wasm.Option,
	bech32Prefix string,
) AppKeepers {
	appKeepers := AppKeepers{}

	// Set keys KVStoreKey
	appKeepers.GenerateKeys()

	// configure state listening capabilities using AppOptions
	// we are doing nothing with the returned streamingServices and waitGroup in this case
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, appKeepers.keys); err != nil {
		tmos.Exit(err.Error())
	}

	paramsKeeper := initParamsKeeper(
		appCodec,
		legacyAmino,
		appKeepers.keys[paramstypes.StoreKey],
		appKeepers.tkeys[paramstypes.TStoreKey],
	)
	appKeepers.ParamsKeeper = &paramsKeeper

	bApp.SetParamStore(appKeepers.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, appKeepers.keys[capabilitytypes.StoreKey], appKeepers.memKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedICAHostKeeper := appKeepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	scopedTransferKeeper := appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedWasmKeeper := appKeepers.CapabilityKeeper.ScopeToModule(wasm.ModuleName)
	scopedICAControllerKeeper := appKeepers.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedLSCosmosKeeper := appKeepers.CapabilityKeeper.ScopeToModule(lscosmostypes.ModuleName)
	appKeepers.CapabilityKeeper.Seal()

	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		appKeepers.keys[authtypes.StoreKey],
		appKeepers.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		maccPerms,
		bech32Prefix,
	)
	appKeepers.AccountKeeper = &accountKeeper

	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec,
		appKeepers.keys[banktypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.GetSubspace(banktypes.ModuleName),
		sendCoinBlockedAddrs, // these blocked address will be used in distribution keeper as well
	)
	appKeepers.BankKeeper = &bankKeeper

	authzKeeper := authzkeeper.NewKeeper(
		appKeepers.keys[authzkeeper.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
		*appKeepers.AccountKeeper,
	)
	appKeepers.AuthzKeeper = &authzKeeper

	feegrantKeeper := feegrantkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[feegrant.StoreKey],
		appKeepers.AccountKeeper,
	)
	appKeepers.FeegrantKeeper = &feegrantKeeper

	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[stakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.GetSubspace(stakingtypes.ModuleName),
	)

	mintKeeper := mintkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[minttypes.StoreKey],
		appKeepers.GetSubspace(minttypes.ModuleName),
		&stakingKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
	)
	appKeepers.MintKeeper = &mintKeeper

	distributionKeeper := distributionkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[distributiontypes.StoreKey],
		appKeepers.GetSubspace(distributiontypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		&stakingKeeper,
		authtypes.FeeCollectorName,
	)
	appKeepers.DistributionKeeper = &distributionKeeper

	slashingKeeper := slashingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[slashingtypes.StoreKey],
		&stakingKeeper,
		appKeepers.GetSubspace(slashingtypes.ModuleName),
	)
	appKeepers.SlashingKeeper = &slashingKeeper

	crisisKeeper := crisiskeeper.NewKeeper(
		appKeepers.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
	)
	appKeepers.CrisisKeeper = &crisisKeeper

	upgradeKeeper := upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		appKeepers.keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		bApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.UpgradeKeeper = &upgradeKeeper

	halvingKeeper := halving.NewKeeper(
		appKeepers.keys[halving.StoreKey],
		appKeepers.GetSubspace(halving.DefaultParamspace),
		appKeepers.MintKeeper,
	)
	appKeepers.HalvingKeeper = &halvingKeeper

	appKeepers.StakingKeeper = stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(appKeepers.DistributionKeeper.Hooks(), appKeepers.SlashingKeeper.Hooks()),
	)

	epochsKeeper := *epochskeeper.NewKeeper(appKeepers.keys[epochstypes.StoreKey])

	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibchost.StoreKey],
		appKeepers.GetSubspace(ibchost.ModuleName),
		appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		scopedIBCKeeper,
	)

	ibcFeeKeeper := ibcfeekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibcfeetypes.StoreKey],
		appKeepers.IBCKeeper.ChannelKeeper, // may be replaced with IBC middleware
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
	)
	appKeepers.IBCFeeKeeper = &ibcFeeKeeper

	transferKeeper := ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		appKeepers.IBCFeeKeeper, // ICS4 Wrapper: fee IBC middleware
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		scopedTransferKeeper,
	)
	appKeepers.TransferKeeper = &transferKeeper

	oracleKeeper := oraclekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[oracletypes.ModuleName],
		appKeepers.GetSubspace(oracletypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.DistributionKeeper,
		&stakingKeeper,
		distributiontypes.ModuleName,
	)
	appKeepers.OracleKeeper = &oracleKeeper

	appKeepers.TransferModule = transfer.NewAppModule(*appKeepers.TransferKeeper)
	appKeepers.TransferIBCModule = transfer.NewIBCModule(*appKeepers.TransferKeeper)

	icaHostKeeper := icahostkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icahosttypes.StoreKey],
		appKeepers.GetSubspace(icahosttypes.SubModuleName),
		appKeepers.IBCFeeKeeper, // use ics29 fee as ics4Wrapper in middleware stack
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		scopedICAHostKeeper,
		bApp.MsgServiceRouter(),
	)
	appKeepers.ICAHostKeeper = &icaHostKeeper

	icaControllerKeeper := icacontrollerkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icacontrollertypes.StoreKey],
		appKeepers.GetSubspace(icacontrollertypes.SubModuleName),
		appKeepers.IBCFeeKeeper, // use ics29 fee as ics4Wrapper in middleware stack
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		scopedICAControllerKeeper,
		bApp.MsgServiceRouter(),
	)
	appKeepers.ICAControllerKeeper = &icaControllerKeeper

	interchainQueryKeeper := interchainquerykeeper.NewKeeper(
		appCodec,
		appKeepers.keys[interchainquerytypes.StoreKey],
		appKeepers.IBCKeeper,
	)
	appKeepers.InterchainQueryKeeper = &interchainQueryKeeper
	appKeepers.InterchainQueryModule = interchainquery.NewAppModule(appCodec, *appKeepers.InterchainQueryKeeper)

	lsCosmosKeeper := lscosmoskeeper.NewKeeper(
		appCodec,
		appKeepers.keys[lscosmostypes.StoreKey],
		appKeepers.memKeys[lscosmostypes.MemStoreKey],
		appKeepers.GetSubspace(lscosmostypes.ModuleName),
		appKeepers.BankKeeper,
		appKeepers.AccountKeeper,
		epochsKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.TransferKeeper,
		appKeepers.ICAControllerKeeper,
		appKeepers.InterchainQueryKeeper,
		scopedLSCosmosKeeper,
		bApp.MsgServiceRouter(),
	)
	appKeepers.LSCosmosKeeper = &lsCosmosKeeper

	err := appKeepers.InterchainQueryKeeper.SetCallbackHandler(lscosmostypes.ModuleName, appKeepers.LSCosmosKeeper.CallbackHandler())
	if err != nil {
		panic(err)
	}

	appKeepers.EpochsKeeper = epochsKeeper.SetHooks(
		epochstypes.NewMultiEpochHooks(appKeepers.LSCosmosKeeper.NewEpochHooks()),
	)

	// Information will flow: ibc-port -> icaController -> lscosmos.
	appKeepers.LSCosmosModule = lscosmos.NewAppModule(appCodec, *appKeepers.LSCosmosKeeper, appKeepers.AccountKeeper, appKeepers.BankKeeper)
	icaControllerIBCModule := icacontroller.NewIBCMiddleware(appKeepers.LSCosmosModule, *appKeepers.ICAControllerKeeper)

	ibcTransferHooksKeeper := ibchookerkeeper.NewKeeper()
	appKeepers.TransferHooksKeeper = ibcTransferHooksKeeper.SetHooks(ibchookertypes.NewMultiStakingHooks(appKeepers.LSCosmosKeeper.NewIBCTransferHooks()))
	appKeepers.IBCTransferHooksMiddleware = ibchooker.NewAppModule(*appKeepers.TransferHooksKeeper, appKeepers.TransferIBCModule)

	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[evidencetypes.StoreKey],
		appKeepers.StakingKeeper,
		appKeepers.SlashingKeeper,
	)
	appKeepers.EvidenceKeeper = evidenceKeeper

	wasmConfig, err := wasm.ReadWasmConfig(appOpts)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	supportedFeatures := "iterator,stargate"
	wasmOpts = append(wasmbindings.RegisterStargateQueries(*bApp.GRPCQueryRouter(), appCodec), wasmOpts...)
	wasmKeeper := wasm.NewKeeper(
		appCodec,
		appKeepers.keys[wasm.StoreKey],
		appKeepers.GetSubspace(wasm.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.DistributionKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		wasmOpts...,
	)
	appKeepers.WasmKeeper = &wasmKeeper

	var icaHostStack ibctypes.IBCModule
	icaHostStack = icahost.NewIBCModule(*appKeepers.ICAHostKeeper)
	icaHostStack = ibcfee.NewIBCMiddleware(icaHostStack, *appKeepers.IBCFeeKeeper)

	// RecvPacket, message that originates from core IBC and goes down to app, the flow is
	// channel.RecvPacket -> fee.OnRecvPacket -> ibchooker.OnRecvPacket -> transfer.OnRecvPacket
	transferStack := ibcfee.NewIBCMiddleware(appKeepers.IBCTransferHooksMiddleware, *appKeepers.IBCFeeKeeper)

	var wasmStack ibctypes.IBCModule
	wasmStack = wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCFeeKeeper)
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, *appKeepers.IBCFeeKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibctypes.NewRouter()
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostStack).
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(icacontrollertypes.SubModuleName, icaControllerIBCModule).
		AddRoute(lscosmostypes.ModuleName, icaControllerIBCModule).
		AddRoute(wasm.ModuleName, wasmStack)
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	govRouter := govv1beta1.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramsproposal.RouterKey, params.NewParamChangeProposalHandler(*appKeepers.ParamsKeeper)).
		AddRoute(distributiontypes.RouterKey, distribution.NewCommunityPoolSpendProposalHandler(*appKeepers.DistributionKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(*appKeepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibccoreclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper)).
		AddRoute(lscosmostypes.RouterKey, lscosmos.NewLSCosmosProposalHandler(*appKeepers.LSCosmosKeeper))

	if len(wasmEnabledProposals) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(appKeepers.WasmKeeper, wasmEnabledProposals))
	}

	govConfig := govtypes.DefaultConfig()
	// Example of setting gov params:
	// govConfig.MaxMetadataLen = 10000

	govKeeper := govkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[govtypes.StoreKey],
		appKeepers.GetSubspace(govtypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		&stakingKeeper,
		govRouter,
		bApp.MsgServiceRouter(),
		govConfig,
	)
	appKeepers.GovKeeper = &govKeeper

	// set scoped keeper from variables in appKeepers
	appKeepers.ScopedIBCKeeper = scopedIBCKeeper
	appKeepers.ScopedTransferKeeper = scopedTransferKeeper
	appKeepers.ScopedICAHostKeeper = scopedICAHostKeeper
	appKeepers.ScopedWasmKeeper = scopedWasmKeeper
	appKeepers.ScopedICAControllerKeeper = scopedICAControllerKeeper
	appKeepers.ScopedLSCosmosKeeper = scopedLSCosmosKeeper

	return appKeepers
}

// GetSubspace returns a param subspace for a given module name.
func (appKeepers *AppKeepers) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := appKeepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// initParamsKeeper init params keeper and its subspaces.
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distributiontypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(halving.DefaultParamspace)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govv1.ParamKeyTable())
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(wasm.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(lscosmostypes.ModuleName)
	paramsKeeper.Subspace(interchainquerytypes.ModuleName)
	paramsKeeper.Subspace(oracletypes.ModuleName)

	return paramsKeeper
}
