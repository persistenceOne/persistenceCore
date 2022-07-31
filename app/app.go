/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package app

import (
	"fmt"
	"io"
	stdlog "log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/CosmWasm/wasmd/x/wasm"
	wasmClient "github.com/CosmWasm/wasmd/x/wasm/client"
	wasmKeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmTypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	serverTypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	storeTypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authRest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authSimulation "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingTypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzKeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzModule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilityKeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilityTypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisisKeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionClient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distributionKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidenceKeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidenceTypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantKeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantModule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilTypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsClient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsProposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingKeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeClient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts"
	icaControllerTypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	icaHost "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host"
	icaHostKeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/keeper"
	icaHostTypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	icaTypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	ibcTransferKeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibcTransferTypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v3/modules/core"
	ibcCoreClient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcClient "github.com/cosmos/ibc-go/v3/modules/core/02-client/client"
	ibcClientTypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibcTypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibcHost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibcKeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	"github.com/gogo/protobuf/grpc"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	tendermintJSON "github.com/tendermint/tendermint/libs/json"
	tendermintLog "github.com/tendermint/tendermint/libs/log"
	tendermintOS "github.com/tendermint/tendermint/libs/os"
	tendermintProto "github.com/tendermint/tendermint/proto/tendermint/types"
	tendermintDB "github.com/tendermint/tm-db"

	appParams "github.com/persistenceOne/persistenceCore/v3/app/params"
	"github.com/persistenceOne/persistenceCore/v3/x/halving"
)

var DefaultNodeHome string

var (
	// ProposalsEnabled is "true" and EnabledSpecificProposals is "", then enable all x/wasm proposals.
	// ProposalsEnabled is not "true" and EnabledSpecificProposals is "", then disable all x/wasm proposals.
	ProposalsEnabled = "true"
	// EnableSpecificProposals if set to non-empty string it must be comma-separated list of values that are all a subset
	// of "EnableAllProposals" (takes precedence over ProposalsEnabled)
	// https://github.com/CosmWasm/wasmd/blob/02a54d33ff2c064f3539ae12d75d027d9c665f05/x/wasm/internal/types/proposal.go#L28-L34
	EnableSpecificProposals = ""
)

// GetEnabledProposals parses the ProposalsEnabled / EnableSpecificProposals values to
// produce a list of enabled proposals to pass into wasmd app.
func GetEnabledProposals() []wasm.ProposalType {
	if EnableSpecificProposals == "" {
		if ProposalsEnabled == "true" {
			return wasm.EnableAllProposals
		}
		return wasm.DisableAllProposals
	}
	chunks := strings.Split(EnableSpecificProposals, ",")
	proposals, err := wasm.ConvertToProposals(chunks)
	if err != nil {
		panic(err)
	}
	return proposals
}

var ModuleAccountPermissions = map[string][]string{
	authTypes.FeeCollectorName:     nil,
	distributionTypes.ModuleName:   nil,
	icaTypes.ModuleName:            nil,
	mintTypes.ModuleName:           {authTypes.Minter},
	stakingTypes.BondedPoolName:    {authTypes.Burner, authTypes.Staking},
	stakingTypes.NotBondedPoolName: {authTypes.Burner, authTypes.Staking},
	govTypes.ModuleName:            {authTypes.Burner},
	ibcTransferTypes.ModuleName:    {authTypes.Minter, authTypes.Burner},
	wasm.ModuleName:                {authTypes.Burner},
}

var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.NewAppModuleBasic(
		append(
			wasmClient.ProposalHandlers,
			paramsClient.ProposalHandler,
			distributionClient.ProposalHandler,
			upgradeClient.ProposalHandler,
			upgradeClient.CancelProposalHandler,
			ibcClient.UpdateClientProposalHandler,
			ibcClient.UpgradeProposalHandler,
		)...,
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feegrantModule.AppModuleBasic{},
	authzModule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	wasm.AppModuleBasic{},
	halving.AppModuleBasic{},
	ica.AppModuleBasic{},
)

var (
	_ simapp.App              = (*Application)(nil)
	_ serverTypes.Application = (*Application)(nil)
)

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		stdlog.Println("Failed to get home dir %2", err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".persistenceCore")
}

type Application struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	applicationCodec  codec.Codec
	interfaceRegistry types.InterfaceRegistry

	keys map[string]*sdk.KVStoreKey

	AccountKeeper      authKeeper.AccountKeeper
	BankKeeper         bankKeeper.Keeper
	CapabilityKeeper   *capabilityKeeper.Keeper
	StakingKeeper      stakingKeeper.Keeper
	SlashingKeeper     slashingKeeper.Keeper
	MintKeeper         mintKeeper.Keeper
	DistributionKeeper distributionKeeper.Keeper
	GovKeeper          govKeeper.Keeper
	UpgradeKeeper      upgradeKeeper.Keeper
	CrisisKeeper       crisisKeeper.Keeper
	ParamsKeeper       paramsKeeper.Keeper
	IBCKeeper          *ibcKeeper.Keeper
	ICAHostKeeper      icaHostKeeper.Keeper
	EvidenceKeeper     evidenceKeeper.Keeper
	TransferKeeper     ibcTransferKeeper.Keeper
	FeegrantKeeper     feegrantKeeper.Keeper
	AuthzKeeper        authzKeeper.Keeper
	HalvingKeeper      halving.Keeper
	WasmKeeper         wasm.Keeper

	moduleManager     *module.Manager
	configurator      module.Configurator
	simulationManager *module.SimulationManager

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilityKeeper.ScopedKeeper
	ScopedTransferKeeper capabilityKeeper.ScopedKeeper
	ScopedICAHostKeeper  capabilityKeeper.ScopedKeeper
	ScopedWasmKeeper     capabilityKeeper.ScopedKeeper
}

func NewApplication(
	applicationName string,
	encodingConfiguration appParams.EncodingConfig,
	moduleAccountPermissions map[string][]string,
	logger tendermintLog.Logger,
	db tendermintDB.DB,
	traceStore io.Writer,
	loadLatest bool,
	invCheckPeriod uint,
	skipUpgradeHeights map[int64]bool,
	home string,
	enabledProposals []wasm.ProposalType,
	applicationOptions serverTypes.AppOptions,
	wasmOpts []wasm.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *Application {

	applicationCodec := encodingConfiguration.Marshaler
	legacyAmino := encodingConfiguration.Amino
	interfaceRegistry := encodingConfiguration.InterfaceRegistry

	baseApp := baseapp.NewBaseApp(
		applicationName,
		logger,
		db,
		encodingConfiguration.TransactionConfig.TxDecoder(),
		baseAppOptions...,
	)
	baseApp.SetCommitMultiStoreTracer(traceStore)
	baseApp.SetVersion(version.Version)
	baseApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authTypes.StoreKey, bankTypes.StoreKey, stakingTypes.StoreKey,
		mintTypes.StoreKey, distributionTypes.StoreKey, slashingTypes.StoreKey,
		govTypes.StoreKey, paramsTypes.StoreKey, ibcHost.StoreKey, upgradeTypes.StoreKey,
		evidenceTypes.StoreKey, ibcTransferTypes.StoreKey, capabilityTypes.StoreKey,
		feegrant.StoreKey, authzKeeper.StoreKey, icaHostTypes.StoreKey, halving.StoreKey, wasm.StoreKey,
	)

	transientStoreKeys := sdk.NewTransientStoreKeys(paramsTypes.TStoreKey)
	memoryKeys := sdk.NewMemoryStoreKeys(capabilityTypes.MemStoreKey)

	app := &Application{
		BaseApp:           baseApp,
		legacyAmino:       legacyAmino,
		applicationCodec:  applicationCodec,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
	}

	app.ParamsKeeper = initParamsKeeper(
		applicationCodec,
		legacyAmino,
		keys[paramsTypes.StoreKey],
		transientStoreKeys[paramsTypes.TStoreKey],
	)
	app.BaseApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramsKeeper.ConsensusParamsKeyTable()))

	app.CapabilityKeeper = capabilityKeeper.NewKeeper(applicationCodec, keys[capabilityTypes.StoreKey], memoryKeys[capabilityTypes.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcHost.ModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icaHostTypes.SubModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibcTransferTypes.ModuleName)
	scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)
	app.CapabilityKeeper.Seal()

	app.AccountKeeper = authKeeper.NewAccountKeeper(
		applicationCodec,
		keys[authTypes.StoreKey],
		app.GetSubspace(authTypes.ModuleName),
		authTypes.ProtoBaseAccount,
		moduleAccountPermissions,
	)

	blacklistedAddresses := make(map[string]bool)
	for account := range moduleAccountPermissions {
		blacklistedAddresses[authTypes.NewModuleAddress(account).String()] = true
	}
	blackListedModuleAddresses := make(map[string]bool)
	for moduleAccount := range moduleAccountPermissions {
		blackListedModuleAddresses[authTypes.NewModuleAddress(moduleAccount).String()] = true
	}

	app.BankKeeper = bankKeeper.NewBaseKeeper(
		applicationCodec,
		keys[bankTypes.StoreKey],
		app.AccountKeeper,
		app.GetSubspace(bankTypes.ModuleName),
		blacklistedAddresses,
	)

	app.AuthzKeeper = authzKeeper.NewKeeper(
		keys[authzKeeper.StoreKey],
		applicationCodec,
		app.BaseApp.MsgServiceRouter(),
	)

	app.FeegrantKeeper = feegrantKeeper.NewKeeper(
		applicationCodec,
		keys[feegrant.StoreKey],
		app.AccountKeeper,
	)

	stakingKeeper := stakingKeeper.NewKeeper(
		applicationCodec,
		keys[stakingTypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(stakingTypes.ModuleName),
	)

	app.MintKeeper = mintKeeper.NewKeeper(
		applicationCodec,
		keys[mintTypes.StoreKey],
		app.GetSubspace(mintTypes.ModuleName),
		&stakingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		authTypes.FeeCollectorName,
	)

	app.DistributionKeeper = distributionKeeper.NewKeeper(
		applicationCodec,
		keys[distributionTypes.StoreKey],
		app.GetSubspace(distributionTypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		authTypes.FeeCollectorName,
		blackListedModuleAddresses,
	)
	app.SlashingKeeper = slashingKeeper.NewKeeper(
		applicationCodec,
		keys[slashingTypes.StoreKey],
		&stakingKeeper,
		app.GetSubspace(slashingTypes.ModuleName),
	)
	app.CrisisKeeper = crisisKeeper.NewKeeper(
		app.GetSubspace(crisisTypes.ModuleName),
		invCheckPeriod,
		app.BankKeeper,
		authTypes.FeeCollectorName,
	)
	app.UpgradeKeeper = upgradeKeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradeTypes.StoreKey],
		applicationCodec,
		home,
		app.BaseApp,
	)

	app.HalvingKeeper = halving.NewKeeper(
		keys[halving.StoreKey],
		app.GetSubspace(halving.DefaultParamspace),
		app.MintKeeper,
	)

	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingTypes.NewMultiStakingHooks(app.DistributionKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.IBCKeeper = ibcKeeper.NewKeeper(
		applicationCodec,
		keys[ibcHost.StoreKey],
		app.GetSubspace(ibcHost.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	app.TransferKeeper = ibcTransferKeeper.NewKeeper(
		applicationCodec,
		keys[ibcTransferTypes.StoreKey],
		app.GetSubspace(ibcTransferTypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(app.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(app.TransferKeeper)

	app.ICAHostKeeper = icaHostKeeper.NewKeeper(
		applicationCodec,
		keys[icaHostTypes.StoreKey],
		app.GetSubspace(icaHostTypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
	)

	icaModule := ica.NewAppModule(nil, &app.ICAHostKeeper)
	icaHostIBCModule := icaHost.NewIBCModule(app.ICAHostKeeper)

	evidenceKeeper := evidenceKeeper.NewKeeper(
		applicationCodec,
		keys[evidenceTypes.StoreKey],
		&app.StakingKeeper,
		app.SlashingKeeper,
	)
	app.EvidenceKeeper = *evidenceKeeper

	wasmDir := filepath.Join(home, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(applicationOptions)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	supportedFeatures := "iterator,staking,stargate"
	app.WasmKeeper = wasm.NewKeeper(
		applicationCodec,
		keys[wasm.StoreKey],
		app.GetSubspace(wasm.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		app.StakingKeeper,
		app.DistributionKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		scopedWasmKeeper,
		app.TransferKeeper,
		app.MsgServiceRouter(),
		app.GRPCQueryRouter(),
		wasmDir,
		wasmConfig,
		supportedFeatures,
		wasmOpts...,
	)
	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcTypes.NewRouter()
	ibcRouter.AddRoute(icaHostTypes.SubModuleName, icaHostIBCModule).
		AddRoute(ibcTransferTypes.ModuleName, transferIBCModule).
		AddRoute(wasm.ModuleName, wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper))
	app.IBCKeeper.SetRouter(ibcRouter)

	govRouter := govTypes.NewRouter()
	govRouter.AddRoute(
		govTypes.RouterKey,
		govTypes.ProposalHandler,
	).AddRoute(
		paramsProposal.RouterKey,
		params.NewParamChangeProposalHandler(app.ParamsKeeper),
	).AddRoute(
		distributionTypes.RouterKey,
		distribution.NewCommunityPoolSpendProposalHandler(app.DistributionKeeper),
	).AddRoute(
		upgradeTypes.RouterKey,
		upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper),
	).AddRoute(ibcClientTypes.RouterKey, ibcCoreClient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))
	if len(enabledProposals) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.WasmKeeper, enabledProposals))
	}
	app.GovKeeper = govKeeper.NewKeeper(
		applicationCodec,
		keys[govTypes.StoreKey],
		app.GetSubspace(govTypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		govRouter,
	)

	/****  Module Options ****/
	var skipGenesisInvariants = false

	opt := applicationOptions.Get(crisis.FlagSkipGenesisInvariants)
	if opt, ok := opt.(bool); ok {
		skipGenesisInvariants = opt
	}

	app.moduleManager = module.NewManager(
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfiguration.TransactionConfig,
		),
		auth.NewAppModule(applicationCodec, app.AccountKeeper, nil),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(applicationCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(applicationCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(applicationCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(applicationCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(applicationCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distribution.NewAppModule(applicationCodec, app.DistributionKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(applicationCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		feegrantModule.NewAppModule(applicationCodec, app.AccountKeeper, app.BankKeeper, app.FeegrantKeeper, app.interfaceRegistry),
		authzModule.NewAppModule(applicationCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		halving.NewAppModule(applicationCodec, app.HalvingKeeper),
		transferModule,
		icaModule,
		wasm.NewAppModule(applicationCodec, &app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
	)

	app.moduleManager.SetOrderBeginBlockers(
		upgradeTypes.ModuleName,
		capabilityTypes.ModuleName,
		crisisTypes.ModuleName,
		govTypes.ModuleName,
		stakingTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcHost.ModuleName,
		icaTypes.ModuleName,
		authTypes.ModuleName,
		bankTypes.ModuleName,
		distributionTypes.ModuleName,
		slashingTypes.ModuleName,
		mintTypes.ModuleName,
		genutilTypes.ModuleName,
		evidenceTypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramsTypes.ModuleName,
		vestingTypes.ModuleName,
		halving.ModuleName,
		wasm.ModuleName,
	)
	app.moduleManager.SetOrderEndBlockers(
		crisisTypes.ModuleName,
		govTypes.ModuleName,
		stakingTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcHost.ModuleName,
		icaTypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		capabilityTypes.ModuleName,
		authTypes.ModuleName,
		bankTypes.ModuleName,
		distributionTypes.ModuleName,
		slashingTypes.ModuleName,
		mintTypes.ModuleName,
		genutilTypes.ModuleName,
		evidenceTypes.ModuleName,
		paramsTypes.ModuleName,
		upgradeTypes.ModuleName,
		vestingTypes.ModuleName,
		halving.ModuleName,
		wasm.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.moduleManager.SetOrderInitGenesis(
		capabilityTypes.ModuleName,
		bankTypes.ModuleName,
		distributionTypes.ModuleName,
		stakingTypes.ModuleName,
		slashingTypes.ModuleName,
		govTypes.ModuleName,
		mintTypes.ModuleName,
		crisisTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcHost.ModuleName,
		icaTypes.ModuleName,
		evidenceTypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		authTypes.ModuleName,
		genutilTypes.ModuleName,
		paramsTypes.ModuleName,
		upgradeTypes.ModuleName,
		vestingTypes.ModuleName,
		halving.ModuleName,
		wasm.ModuleName,
	)

	app.moduleManager.RegisterInvariants(&app.CrisisKeeper)
	app.moduleManager.RegisterRoutes(app.BaseApp.Router(), app.BaseApp.QueryRouter(), encodingConfiguration.Amino)
	app.configurator = module.NewConfigurator(app.applicationCodec, app.BaseApp.MsgServiceRouter(), app.BaseApp.GRPCQueryRouter())
	app.moduleManager.RegisterServices(app.configurator)

	simulationManager := module.NewSimulationManager(
		auth.NewAppModule(applicationCodec, app.AccountKeeper, authSimulation.RandomGenesisAccounts),
		bank.NewAppModule(applicationCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(applicationCodec, *app.CapabilityKeeper),
		gov.NewAppModule(applicationCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(applicationCodec, app.MintKeeper, app.AccountKeeper),
		staking.NewAppModule(applicationCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		distribution.NewAppModule(applicationCodec, app.DistributionKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		slashing.NewAppModule(applicationCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		params.NewAppModule(app.ParamsKeeper),
		halving.NewAppModule(applicationCodec, app.HalvingKeeper),
		authzModule.NewAppModule(applicationCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		feegrantModule.NewAppModule(applicationCodec, app.AccountKeeper, app.BankKeeper, app.FeegrantKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		transferModule,
	)

	simulationManager.RegisterStoreDecoders()
	app.simulationManager = simulationManager

	app.BaseApp.MountKVStores(keys)
	app.BaseApp.MountTransientStores(transientStoreKeys)
	app.BaseApp.MountMemoryStores(memoryKeys)

	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   app.AccountKeeper,
				BankKeeper:      app.BankKeeper,
				FeegrantKeeper:  app.FeegrantKeeper,
				SignModeHandler: encodingConfiguration.TransactionConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			IBCKeeper:         app.IBCKeeper,
			WasmConfig:        &wasmConfig,
			TXCounterStoreKey: keys[wasm.StoreKey],
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	app.BaseApp.SetAnteHandler(anteHandler)
	app.BaseApp.SetInitChainer(app.InitChainer)
	app.BaseApp.SetBeginBlocker(app.moduleManager.BeginBlock)
	app.BaseApp.SetEndBlocker(app.moduleManager.EndBlock)

	// must be before Loading version
	// requires the snapshot store to be created and registered as a BaseAppOption
	// see cmd/wasmd/root.go: 206 - 214 approx
	if manager := app.SnapshotManager(); manager != nil {
		err := manager.RegisterExtensions(
			wasmKeeper.NewWasmSnapshotter(app.CommitMultiStore(), &app.WasmKeeper),
		)
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %s", err))
		}
	}

	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx sdk.Context, _ upgradeTypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("starting the upgrade now")

			fromVM[icaTypes.ModuleName] = icaModule.ConsensusVersion()
			// create ICS27 Controller submodule params
			controllerParams := icaControllerTypes.Params{}
			// create ICS27 Host submodule params
			hostParams := icaHostTypes.Params{
				HostEnabled: true,
				AllowMessages: []string{
					authzMsgExec,
					authzMsgGrant,
					authzMsgRevoke,
					bankMsgSend,
					bankMsgMultiSend,
					distrMsgSetWithdrawAddr,
					distrMsgWithdrawValidatorCommission,
					distrMsgFundCommunityPool,
					distrMsgWithdrawDelegatorReward,
					feegrantMsgGrantAllowance,
					feegrantMsgRevokeAllowance,
					govMsgVoteWeighted,
					govMsgSubmitProposal,
					govMsgDeposit,
					govMsgVote,
					stakingMsgEditValidator,
					stakingMsgDelegate,
					stakingMsgUndelegate,
					stakingMsgBeginRedelegate,
					stakingMsgCreateValidator,
					vestingMsgCreateVestingAccount,
					transferMsgTransfer,
				},
			}
			ctx.Logger().Info("start to init interchainaccount module...")
			// initialize ICS27 module
			icaModule.InitModule(ctx, controllerParams, hostParams)

			ctx.Logger().Info("start to run module migrations...")

			// RunMigrations twice is just a way to make auth module's migrates after staking
			newVM, err := app.moduleManager.RunMigrations(ctx, app.configurator, fromVM)
			if err != nil {
				return nil, err
			}

			// Since we provide custom DefaultGenesis (privileges StoreCode) in
			// app/genesis.go rather than the wasm module, we need to set the params
			// here when migrating (is it is not customized).
			params := app.WasmKeeper.GetParams(ctx)
			params.CodeUploadAccess = wasmTypes.AllowNobody
			params.InstantiateDefaultPermission = wasmTypes.AccessTypeNobody
			app.WasmKeeper.SetParams(ctx, params)

			return newVM, nil
		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storeTypes.StoreUpgrades{
			Added: []string{icaHostTypes.StoreKey, wasm.ModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.BaseApp.SetStoreLoader(upgradeTypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}

	if loadLatest {
		if err := app.BaseApp.LoadLatestVersion(); err != nil {
			tendermintOS.Exit(err.Error())
		}
		ctx := app.BaseApp.NewUncachedContext(true, tendermintProto.Header{})

		// Initialize pinned codes in wasmvm as they are not persisted there
		if err := app.WasmKeeper.InitializePinnedCodes(ctx); err != nil {
			tendermintOS.Exit(fmt.Sprintf("failed initialize pinned codes %s", err))
		}
	}
	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	app.ScopedICAHostKeeper = scopedICAHostKeeper
	app.ScopedWasmKeeper = scopedWasmKeeper

	return app
}

func (app Application) ApplicationCodec() codec.Codec {
	return app.applicationCodec
}

func (app Application) Name() string {
	return app.BaseApp.Name()
}

func (app Application) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app Application) BeginBlocker(ctx sdk.Context, req abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	return app.moduleManager.BeginBlock(ctx, req)
}

func (app Application) EndBlocker(ctx sdk.Context, req abciTypes.RequestEndBlock) abciTypes.ResponseEndBlock {
	return app.moduleManager.EndBlock(ctx, req)
}

func (app Application) InitChainer(ctx sdk.Context, req abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
	var genesisState GenesisState
	if err := tendermintJSON.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.moduleManager.GetVersionMap())

	return app.moduleManager.InitGenesis(ctx, app.applicationCodec, genesisState)
}

func (app Application) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string) (serverTypes.ExportedApp, error) {
	context := app.BaseApp.NewContext(true, tendermintProto.Header{Height: app.BaseApp.LastBlockHeight()})

	height := app.BaseApp.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		applyWhiteList := false

		if len(jailWhiteList) > 0 {
			applyWhiteList = true
		}

		whiteListMap := make(map[string]bool)

		for _, address := range jailWhiteList {
			if _, Error := sdk.ValAddressFromBech32(address); Error != nil {
				panic(Error)
			}

			whiteListMap[address] = true
		}

		app.CrisisKeeper.AssertInvariants(context)

		app.StakingKeeper.IterateValidators(context, func(_ int64, val stakingTypes.ValidatorI) (stop bool) {
			_, _ = app.DistributionKeeper.WithdrawValidatorCommission(context, val.GetOperator())
			return false
		})

		delegations := app.StakingKeeper.GetAllDelegations(context)
		for _, delegation := range delegations {
			validatorAddress, Error := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
			if Error != nil {
				panic(Error)
			}

			delegatorAddress, Error := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
			if Error != nil {
				panic(Error)
			}

			_, _ = app.DistributionKeeper.WithdrawDelegationRewards(context, delegatorAddress, validatorAddress)
		}

		app.DistributionKeeper.DeleteAllValidatorSlashEvents(context)

		app.DistributionKeeper.DeleteAllValidatorHistoricalRewards(context)

		height := context.BlockHeight()
		context = context.WithBlockHeight(0)

		app.StakingKeeper.IterateValidators(context, func(_ int64, val stakingTypes.ValidatorI) (stop bool) {

			scraps := app.DistributionKeeper.GetValidatorOutstandingRewardsCoins(context, val.GetOperator())
			feePool := app.DistributionKeeper.GetFeePool(context)
			feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
			app.DistributionKeeper.SetFeePool(context, feePool)

			app.DistributionKeeper.Hooks().AfterValidatorCreated(context, val.GetOperator())
			return false
		})

		for _, delegation := range delegations {
			validatorAddress, Error := sdk.ValAddressFromBech32(delegation.ValidatorAddress)
			if Error != nil {
				panic(Error)
			}

			delegatorAddress, Error := sdk.AccAddressFromBech32(delegation.DelegatorAddress)
			if Error != nil {
				panic(Error)
			}

			app.DistributionKeeper.Hooks().BeforeDelegationCreated(context, delegatorAddress, validatorAddress)
			app.DistributionKeeper.Hooks().AfterDelegationModified(context, delegatorAddress, validatorAddress)
		}

		context = context.WithBlockHeight(height)

		app.StakingKeeper.IterateRedelegations(context, func(_ int64, redelegation stakingTypes.Redelegation) (stop bool) {
			for i := range redelegation.Entries {
				redelegation.Entries[i].CreationHeight = 0
			}
			app.StakingKeeper.SetRedelegation(context, redelegation)
			return false
		})

		app.StakingKeeper.IterateUnbondingDelegations(context, func(_ int64, unbondingDelegation stakingTypes.UnbondingDelegation) (stop bool) {
			for i := range unbondingDelegation.Entries {
				unbondingDelegation.Entries[i].CreationHeight = 0
			}
			app.StakingKeeper.SetUnbondingDelegation(context, unbondingDelegation)
			return false
		})

		store := context.KVStore(app.keys[stakingTypes.StoreKey])
		kvStoreReversePrefixIterator := sdk.KVStoreReversePrefixIterator(store, stakingTypes.ValidatorsKey)
		counter := int16(0)

		for ; kvStoreReversePrefixIterator.Valid(); kvStoreReversePrefixIterator.Next() {
			addr := sdk.ValAddress(kvStoreReversePrefixIterator.Key()[1:])
			validator, found := app.StakingKeeper.GetValidator(context, addr)

			if !found {
				panic("Validator not found!")
			}

			validator.UnbondingHeight = 0

			if applyWhiteList && !whiteListMap[addr.String()] {
				validator.Jailed = true
			}

			app.StakingKeeper.SetValidator(context, validator)
			counter++
		}

		_ = kvStoreReversePrefixIterator.Close()

		_, Error := app.StakingKeeper.ApplyAndReturnValidatorSetUpdates(context)
		if Error != nil {
			stdlog.Fatal(Error)
		}

		app.SlashingKeeper.IterateValidatorSigningInfos(
			context,
			func(validatorConsAddress sdk.ConsAddress, validatorSigningInfo slashingTypes.ValidatorSigningInfo) (stop bool) {
				validatorSigningInfo.StartHeight = 0
				app.SlashingKeeper.SetValidatorSigningInfo(context, validatorConsAddress, validatorSigningInfo)
				return false
			},
		)
	}

	genesisState := app.moduleManager.ExportGenesis(context, app.applicationCodec)
	applicationState, Error := codec.MarshalJSONIndent(app.legacyAmino, genesisState)

	if Error != nil {
		return serverTypes.ExportedApp{}, Error
	}

	validators, err := staking.WriteValidators(context, app.StakingKeeper)

	return serverTypes.ExportedApp{
		AppState:        applicationState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(context),
	}, err
}

func (app Application) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range ModuleAccountPermissions {
		modAccAddrs[authTypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *Application) GetSubspace(moduleName string) paramsTypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

func (app Application) SimulationManager() *module.SimulationManager {
	return app.simulationManager
}

func (app Application) ListSnapshots(snapshots abciTypes.RequestListSnapshots) abciTypes.ResponseListSnapshots {
	return app.BaseApp.ListSnapshots(snapshots)
}

func (app Application) OfferSnapshot(snapshot abciTypes.RequestOfferSnapshot) abciTypes.ResponseOfferSnapshot {
	return app.BaseApp.OfferSnapshot(snapshot)
}

func (app Application) LoadSnapshotChunk(chunk abciTypes.RequestLoadSnapshotChunk) abciTypes.ResponseLoadSnapshotChunk {
	return app.BaseApp.LoadSnapshotChunk(chunk)
}

func (app Application) ApplySnapshotChunk(chunk abciTypes.RequestApplySnapshotChunk) abciTypes.ResponseApplySnapshotChunk {
	return app.BaseApp.ApplySnapshotChunk(chunk)
}

func (app Application) RegisterGRPCServer(server grpc.Server) {
	app.BaseApp.RegisterGRPCServer(server)
}

func (app Application) RegisterAPIRoutes(apiServer *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiServer.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiServer.Router)
	// Register legacy tx routes.
	authRest.RegisterTxRoutes(clientCtx, apiServer.Router)
	// Register new tx routes from grpc-gateway.
	authTx.RegisterGRPCGatewayRoutes(clientCtx, apiServer.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiServer.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiServer.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiServer.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(apiServer.Router)
	}
}

func RegisterSwaggerAPI(rtr *mux.Router) {
	statikFS, err := fs.New()
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(statikFS)
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

func (app Application) RegisterTxService(clientContect client.Context) {
	authTx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientContect, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app Application) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

func (app Application) LoadHeight(height int64) error {
	return app.BaseApp.LoadVersion(height)
}

// initParamsKeeper init params keeper and its subspaces.
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramsKeeper.Keeper {
	paramsKeeper := paramsKeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authTypes.ModuleName)
	paramsKeeper.Subspace(bankTypes.ModuleName)
	paramsKeeper.Subspace(stakingTypes.ModuleName)
	paramsKeeper.Subspace(mintTypes.ModuleName)
	paramsKeeper.Subspace(distributionTypes.ModuleName)
	paramsKeeper.Subspace(slashingTypes.ModuleName)
	paramsKeeper.Subspace(crisisTypes.ModuleName)
	paramsKeeper.Subspace(halving.DefaultParamspace)
	paramsKeeper.Subspace(govTypes.ModuleName).WithKeyTable(govTypes.ParamKeyTable())
	paramsKeeper.Subspace(ibcTransferTypes.ModuleName)
	paramsKeeper.Subspace(ibcHost.ModuleName)
	paramsKeeper.Subspace(icaHostTypes.SubModuleName)
	paramsKeeper.Subspace(wasm.ModuleName)

	return paramsKeeper
}
