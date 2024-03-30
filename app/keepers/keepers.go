package keepers

import (
	"fmt"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	tmos "github.com/cometbft/cometbft/libs/os"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
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
	consensusparamskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	consensusparamstypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	group "github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	packetforward "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward"
	packetforwardkeeper "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/keeper"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7"
	ibchookskeeper "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/keeper"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
	icacontroller "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/host/types"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/keeper"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibctransfer "github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v7/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibccoreclient "github.com/cosmos/ibc-go/v7/modules/core/02-client"
	ibcclienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	ibctypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	epochskeeper "github.com/persistenceOne/persistence-sdk/v2/x/epochs/keeper"
	epochstypes "github.com/persistenceOne/persistence-sdk/v2/x/epochs/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/halving"
	"github.com/persistenceOne/persistence-sdk/v2/x/ibchooker"
	ibchookerkeeper "github.com/persistenceOne/persistence-sdk/v2/x/ibchooker/keeper"
	ibchookertypes "github.com/persistenceOne/persistence-sdk/v2/x/ibchooker/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/interchainquery"
	interchainquerykeeper "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/keeper"
	interchainquerytypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"
	oraclekeeper "github.com/persistenceOne/persistence-sdk/v2/x/oracle/keeper"
	oracletypes "github.com/persistenceOne/persistence-sdk/v2/x/oracle/types"
	liquidstakekeeper "github.com/persistenceOne/pstake-native/v2/x/liquidstake/keeper"
	liquidstaketypes "github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc"
	liquidstakeibckeeper "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/keeper"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	"github.com/persistenceOne/pstake-native/v2/x/ratesync"
	ratesynckeeper "github.com/persistenceOne/pstake-native/v2/x/ratesync/keeper"
	ratesynctypes "github.com/persistenceOne/pstake-native/v2/x/ratesync/types"
	builderkeeper "github.com/skip-mev/pob/x/builder/keeper"
	buildertypes "github.com/skip-mev/pob/x/builder/types"
	"github.com/spf13/cast"

	"github.com/persistenceOne/persistenceCore/v11/wasmbindings"

	// unnamed import of statik for swagger UI support
	_ "github.com/cosmos/cosmos-sdk/client/docs/statik"
	"github.com/cosmos/cosmos-sdk/client/flags"
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
	WasmKeeper            *wasmkeeper.Keeper
	EpochsKeeper          *epochskeeper.Keeper
	OracleKeeper          *oraclekeeper.Keeper
	ICAControllerKeeper   *icacontrollerkeeper.Keeper
	InterchainQueryKeeper *interchainquerykeeper.Keeper
	TransferHooksKeeper   *ibchookerkeeper.Keeper
	LiquidStakeIBCKeeper  *liquidstakeibckeeper.Keeper
	LiquidStakeKeeper     *liquidstakekeeper.Keeper
	RateSyncKeeper        *ratesynckeeper.Keeper
	ConsensusParamsKeeper *consensusparamskeeper.Keeper
	GroupKeeper           *groupkeeper.Keeper
	PacketForwardKeeper   *packetforwardkeeper.Keeper
	BuilderKeeper         *builderkeeper.Keeper

	// Modules
	TransferModule             ibctransfer.AppModule
	InterchainQueryModule      interchainquery.AppModule
	IBCTransferHooksMiddleware ibchooker.AppModule

	// IBC hooks
	IBCHooksKeeper   *ibchookskeeper.Keeper
	ICS20WasmHooks   *ibchooks.WasmHooks
	HooksICS4Wrapper *ibchooks.ICS4Middleware

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
	appOpts servertypes.AppOptions,
	wasmDir string,
	wasmOpts []wasmkeeper.Option,
	bech32Prefix string,
) *AppKeepers {
	appKeepers := &AppKeepers{}

	// Set keys KVStoreKey
	appKeepers.GenerateKeys()

	// configure state listening capabilities using AppOptions
	// we are doing nothing with the returned streamingServices and waitGroup in this case
	if _, _, err := streaming.LoadStreamingServices(bApp, appOpts, appCodec, bApp.Logger(), appKeepers.keys); err != nil {
		tmos.Exit(err.Error())
	}

	paramsKeeper := initParamsKeeper(
		appCodec,
		legacyAmino,
		appKeepers.keys[paramstypes.StoreKey],
		appKeepers.tkeys[paramstypes.TStoreKey],
	)
	appKeepers.ParamsKeeper = &paramsKeeper

	consensusKeeper := consensusparamskeeper.NewKeeper(
		appCodec,
		appKeepers.keys[consensusparamstypes.StoreKey],
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.ConsensusParamsKeeper = &consensusKeeper

	bApp.SetParamStore(appKeepers.ConsensusParamsKeeper)

	// add capability keeper and ScopeToModule for ibc module
	appKeepers.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, appKeepers.keys[capabilitytypes.StoreKey], appKeepers.memKeys[capabilitytypes.MemStoreKey])
	appKeepers.ScopedIBCKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	appKeepers.ScopedICAHostKeeper = appKeepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	appKeepers.ScopedTransferKeeper = appKeepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	appKeepers.ScopedWasmKeeper = appKeepers.CapabilityKeeper.ScopeToModule(wasmtypes.ModuleName)
	appKeepers.ScopedICAControllerKeeper = appKeepers.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	appKeepers.CapabilityKeeper.Seal()

	accountKeeper := authkeeper.NewAccountKeeper(
		appCodec,
		appKeepers.keys[authtypes.StoreKey],
		authtypes.ProtoBaseAccount,
		maccPerms,
		bech32Prefix,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.AccountKeeper = &accountKeeper

	bankKeeper := bankkeeper.NewBaseKeeper(
		appCodec,
		appKeepers.keys[banktypes.StoreKey],
		appKeepers.AccountKeeper,
		sendCoinBlockedAddrs, // these blocked address will be used in distribution keeper as well
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.BankKeeper = &bankKeeper

	authzKeeper := authzkeeper.NewKeeper(
		appKeepers.keys[authzkeeper.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
		*appKeepers.AccountKeeper,
	)
	appKeepers.AuthzKeeper = &authzKeeper

	groupConfig := group.DefaultConfig()
	/*
		Example of setting group params:
		groupConfig.MaxMetadataLen = 1000
	*/
	groupKeeper := groupkeeper.NewKeeper(
		appKeepers.keys[group.StoreKey],
		appCodec,
		bApp.MsgServiceRouter(),
		appKeepers.AccountKeeper,
		groupConfig,
	)
	appKeepers.GroupKeeper = &groupKeeper

	feegrantKeeper := feegrantkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[feegrant.StoreKey],
		appKeepers.AccountKeeper,
	)
	appKeepers.FeegrantKeeper = &feegrantKeeper

	appKeepers.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[stakingtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	mintKeeper := mintkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[minttypes.StoreKey],
		appKeepers.StakingKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.MintKeeper = &mintKeeper

	distributionKeeper := distributionkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[distributiontypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.DistributionKeeper = &distributionKeeper

	slashingKeeper := slashingkeeper.NewKeeper(
		appCodec,
		legacyAmino,
		appKeepers.keys[slashingtypes.StoreKey],
		appKeepers.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.SlashingKeeper = &slashingKeeper

	invCheckPeriod := cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod))
	appKeepers.CrisisKeeper = crisiskeeper.NewKeeper(
		appCodec,
		appKeepers.keys[crisistypes.StoreKey],
		invCheckPeriod,
		appKeepers.BankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// get skipUpgradeHeights from the app options
	skipUpgradeHeights := map[int64]bool{}
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}
	homePath := cast.ToString(appOpts.Get(flags.FlagHome))
	appKeepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		appKeepers.keys[upgradetypes.StoreKey],
		appCodec,
		homePath,
		bApp,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	halvingKeeper := halving.NewKeeper(
		appKeepers.keys[halving.StoreKey],
		appKeepers.GetSubspace(halving.DefaultParamspace),
		appKeepers.MintKeeper,
	)
	appKeepers.HalvingKeeper = &halvingKeeper

	appKeepers.StakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(appKeepers.DistributionKeeper.Hooks(), appKeepers.SlashingKeeper.Hooks()),
	)

	appKeepers.EpochsKeeper = epochskeeper.NewKeeper(appKeepers.keys[epochstypes.StoreKey])

	appKeepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibcexported.StoreKey],
		appKeepers.GetSubspace(ibcexported.ModuleName),
		appKeepers.StakingKeeper,
		appKeepers.UpgradeKeeper,
		appKeepers.ScopedIBCKeeper,
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

	appKeepers.PacketForwardKeeper = packetforwardkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[packetforwardtypes.StoreKey],
		appKeepers.TransferKeeper, // Will be zero-value here. Reference is set later on with SetTransferKeeper.
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.DistributionKeeper,
		appKeepers.BankKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	transferKeeper := ibctransferkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[ibctransfertypes.StoreKey],
		appKeepers.GetSubspace(ibctransfertypes.ModuleName),
		// The ICS4Wrapper is replaced by the PacketForwardKeeper
		// so that sending can be overridden by the middleware
		appKeepers.PacketForwardKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.ScopedTransferKeeper,
	)
	appKeepers.TransferKeeper = &transferKeeper
	appKeepers.TransferModule = ibctransfer.NewAppModule(*appKeepers.TransferKeeper)
	appKeepers.PacketForwardKeeper.SetTransferKeeper(*appKeepers.TransferKeeper)

	oracleKeeper := oraclekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[oracletypes.ModuleName],
		appKeepers.GetSubspace(oracletypes.ModuleName),
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.DistributionKeeper,
		appKeepers.StakingKeeper,
		distributiontypes.ModuleName,
	)
	appKeepers.OracleKeeper = &oracleKeeper

	icaHostKeeper := icahostkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[icahosttypes.StoreKey],
		appKeepers.GetSubspace(icahosttypes.SubModuleName),
		appKeepers.IBCFeeKeeper, // use ics29 fee as ics4Wrapper in middleware stack
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.AccountKeeper,
		appKeepers.ScopedICAHostKeeper,
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
		appKeepers.ScopedICAControllerKeeper,
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

	liquidStakeIBCKeeper := liquidstakeibckeeper.NewKeeper(
		appCodec,
		appKeepers.keys[liquidstakeibctypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.EpochsKeeper,
		appKeepers.ICAControllerKeeper,
		appKeepers.IBCKeeper,
		appKeepers.TransferKeeper,
		appKeepers.InterchainQueryKeeper,
		appKeepers.GetSubspace(liquidstakeibctypes.ModuleName),
		bApp.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.LiquidStakeIBCKeeper = &liquidStakeIBCKeeper

	liquidStakeKeeper := liquidstakekeeper.NewKeeper(
		appCodec,
		appKeepers.keys[liquidstaketypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		appKeepers.MintKeeper,
		appKeepers.DistributionKeeper,
		appKeepers.SlashingKeeper,
		bApp.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.LiquidStakeKeeper = &liquidStakeKeeper

	err := appKeepers.InterchainQueryKeeper.SetCallbackHandler(liquidstakeibctypes.ModuleName, appKeepers.LiquidStakeIBCKeeper.CallbackHandler())
	if err != nil {
		panic(err)
	}

	liquidStakeIBCModule := liquidstakeibc.NewIBCModule(*appKeepers.LiquidStakeIBCKeeper)
	appKeepers.RateSyncKeeper = ratesynckeeper.NewKeeper(appCodec,
		appKeepers.keys[ratesynctypes.StoreKey],
		appKeepers.EpochsKeeper,
		appKeepers.LiquidStakeKeeper,
		appKeepers.ICAControllerKeeper,
		appKeepers.IBCKeeper, bApp.MsgServiceRouter(),
		authtypes.NewModuleAddress(govtypes.ModuleName).String())

	appKeepers.LiquidStakeIBCKeeper.SetHooks(liquidstakeibctypes.NewMultiLiquidStakeIBCHooks(
		appKeepers.RateSyncKeeper.LiquidStakeIBCHooks()),
	)
	appKeepers.EpochsKeeper.SetHooks(
		epochstypes.NewMultiEpochHooks(
			appKeepers.LiquidStakeIBCKeeper.NewEpochHooks(), appKeepers.RateSyncKeeper.EpochHooks(),
		),
	)

	ibcTransferHooksKeeper := ibchookerkeeper.NewKeeper()
	appKeepers.TransferHooksKeeper = ibcTransferHooksKeeper.SetHooks(
		ibchookertypes.NewMultiStakingHooks(
			appKeepers.LiquidStakeIBCKeeper.NewIBCTransferHooks(),
		),
	)

	transferIBCModule := ibctransfer.NewIBCModule(*appKeepers.TransferKeeper)
	appKeepers.IBCTransferHooksMiddleware = ibchooker.NewAppModule(*appKeepers.TransferHooksKeeper, transferIBCModule)

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
	// checkout: https://github.com/CosmWasm/cosmwasm/blob/main/docs/CAPABILITIES-BUILT-IN.md
	availableCapabilities := "iterator,staking,stargate,cosmwasm_1_1,cosmwasm_1_2,cosmwasm_1_3,cosmwasm_1_4"
	wasmOpts = append(wasmbindings.RegisterStargateQueries(bApp.GRPCQueryRouter(), appCodec), wasmOpts...)
	wasmKeeper := wasmkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[wasmtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		distributionkeeper.NewQuerier(*appKeepers.DistributionKeeper),
		appKeepers.IBCFeeKeeper,
		appKeepers.IBCKeeper.ChannelKeeper,
		&appKeepers.IBCKeeper.PortKeeper,
		appKeepers.ScopedWasmKeeper,
		appKeepers.TransferKeeper,
		bApp.MsgServiceRouter(),
		bApp.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		availableCapabilities,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		wasmOpts...,
	)
	appKeepers.WasmKeeper = &wasmKeeper

	// Configure the hooks keeper
	hooksKeeper := ibchookskeeper.NewKeeper(appKeepers.keys[ibchookstypes.StoreKey])
	appKeepers.IBCHooksKeeper = &hooksKeeper

	wasmHooks := ibchooks.NewWasmHooks(&hooksKeeper, appKeepers.WasmKeeper, bech32Prefix)
	appKeepers.ICS20WasmHooks = &wasmHooks

	hooksICS4Wrapper := ibchooks.NewICS4Middleware(
		appKeepers.IBCKeeper.ChannelKeeper,
		appKeepers.ICS20WasmHooks,
	)
	appKeepers.HooksICS4Wrapper = &hooksICS4Wrapper

	builderKeeper := builderkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[buildertypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.DistributionKeeper,
		appKeepers.StakingKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	appKeepers.BuilderKeeper = &builderKeeper

	var icaHostStack ibctypes.IBCModule
	icaHostStack = icahost.NewIBCModule(*appKeepers.ICAHostKeeper)
	icaHostStack = ibcfee.NewIBCMiddleware(icaHostStack, *appKeepers.IBCFeeKeeper)

	var transferStack ibctypes.IBCModule = appKeepers.IBCTransferHooksMiddleware
	transferStack = ibcfee.NewIBCMiddleware(transferStack, *appKeepers.IBCFeeKeeper)
	transferStack = ibchooks.NewIBCMiddleware(transferStack, appKeepers.HooksICS4Wrapper)
	transferStack = packetforward.NewIBCMiddleware(
		transferStack,
		appKeepers.PacketForwardKeeper,
		0, // no retries on timeout
		packetforwardkeeper.DefaultForwardTransferPacketTimeoutTimestamp,
		packetforwardkeeper.DefaultRefundTransferPacketTimeoutTimestamp,
	)

	// Information will flow: ibc-port -> icaController -> lscosmos.
	var icaControllerStack ibctypes.IBCModule = liquidStakeIBCModule
	icaControllerStack = ratesync.NewIBCModule(icaControllerStack, *appKeepers.RateSyncKeeper)
	icaControllerStack = icacontroller.NewIBCMiddleware(icaControllerStack, *appKeepers.ICAControllerKeeper)

	var wasmStack ibctypes.IBCModule
	wasmStack = wasm.NewIBCHandler(appKeepers.WasmKeeper, appKeepers.IBCKeeper.ChannelKeeper, appKeepers.IBCFeeKeeper)
	wasmStack = ibcfee.NewIBCMiddleware(wasmStack, *appKeepers.IBCFeeKeeper)

	ibcRouter := ibctypes.NewRouter().
		AddRoute(ibctransfertypes.ModuleName, transferStack).
		AddRoute(wasmtypes.ModuleName, wasmStack).
		AddRoute(icacontrollertypes.SubModuleName, icaControllerStack).
		AddRoute(icahosttypes.SubModuleName, icaHostStack).
		AddRoute(liquidstakeibctypes.ModuleName, icaControllerStack).
		AddRoute(ratesynctypes.ModuleName, icaControllerStack)
	appKeepers.IBCKeeper.SetRouter(ibcRouter)

	govRouter := govv1beta1.NewRouter()
	govRouter.
		AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).
		AddRoute(paramsproposal.RouterKey, params.NewParamChangeProposalHandler(*appKeepers.ParamsKeeper)). // this is kept for the modules that are yet to migrate from legacy x/params.
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(appKeepers.UpgradeKeeper)).
		AddRoute(ibcclienttypes.RouterKey, ibccoreclient.NewClientProposalHandler(appKeepers.IBCKeeper.ClientKeeper))

	govConfig := govtypes.DefaultConfig()
	govConfig.MaxMetadataLen = 5000

	appKeepers.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		appKeepers.keys[govtypes.StoreKey],
		appKeepers.AccountKeeper,
		appKeepers.BankKeeper,
		appKeepers.StakingKeeper,
		bApp.MsgServiceRouter(),
		govConfig,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	// Set legacy router for backwards compatibility with gov v1beta1
	appKeepers.GovKeeper.SetLegacyRouter(govRouter)

	appKeepers.GovKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
		// register the governance hooks
		),
	)

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
	paramsKeeper.Subspace(govtypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(wasmtypes.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(oracletypes.ModuleName)
	paramsKeeper.Subspace(packetforwardtypes.ModuleName)

	return paramsKeeper
}
