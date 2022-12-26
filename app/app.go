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
	wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authsimulation "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distributionkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantkeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts"
	icacontroller "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller"
	icacontrollerkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/keeper"
	icacontrollertypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	icahost "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host"
	icahostkeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v3/modules/core"
	ibccoreclient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcclient "github.com/cosmos/ibc-go/v3/modules/core/02-client/client"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibctypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibchost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	ibckeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	"github.com/gogo/protobuf/grpc"
	"github.com/gorilla/mux"
	"github.com/persistenceOne/persistence-sdk/x/epochs"
	epochskeeper "github.com/persistenceOne/persistence-sdk/x/epochs/keeper"
	epochstypes "github.com/persistenceOne/persistence-sdk/x/epochs/types"
	"github.com/persistenceOne/persistence-sdk/x/halving"
	"github.com/persistenceOne/persistence-sdk/x/ibchooker"
	ibchookerkeeper "github.com/persistenceOne/persistence-sdk/x/ibchooker/keeper"
	ibchookertypes "github.com/persistenceOne/persistence-sdk/x/ibchooker/types"
	"github.com/persistenceOne/persistence-sdk/x/interchainquery"
	interchainquerykeeper "github.com/persistenceOne/persistence-sdk/x/interchainquery/keeper"
	interchainquerytypes "github.com/persistenceOne/persistence-sdk/x/interchainquery/types"
	"github.com/persistenceOne/pstake-native/x/lscosmos"
	lscosmosclient "github.com/persistenceOne/pstake-native/x/lscosmos/client"
	lscosmoskeeper "github.com/persistenceOne/pstake-native/x/lscosmos/keeper"
	lscosmostypes "github.com/persistenceOne/pstake-native/x/lscosmos/types"
	"github.com/rakyll/statik/fs"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tendermintjson "github.com/tendermint/tendermint/libs/json"
	tendermintlog "github.com/tendermint/tendermint/libs/log"
	tendermintos "github.com/tendermint/tendermint/libs/os"
	tendermintproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tendermintdb "github.com/tendermint/tm-db"

	appparams "github.com/persistenceOne/persistenceCore/v6/app/params"
	upgrades "github.com/persistenceOne/persistenceCore/v6/app/upgrades/v6"
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
	authtypes.FeeCollectorName:               nil,
	distributiontypes.ModuleName:             nil,
	icatypes.ModuleName:                      nil,
	minttypes.ModuleName:                     {authtypes.Minter},
	stakingtypes.BondedPoolName:              {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName:           {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:                      {authtypes.Burner},
	ibctransfertypes.ModuleName:              {authtypes.Minter, authtypes.Burner},
	wasm.ModuleName:                          {authtypes.Burner},
	lscosmostypes.ModuleName:                 {authtypes.Minter, authtypes.Burner},
	lscosmostypes.DepositModuleAccount:       nil,
	lscosmostypes.DelegationModuleAccount:    nil,
	lscosmostypes.RewardModuleAccount:        nil,
	lscosmostypes.UndelegationModuleAccount:  nil,
	lscosmostypes.RewardBoosterModuleAccount: nil,
}

var receiveAllowedMAcc = map[string]bool{
	lscosmostypes.UndelegationModuleAccount: true,
	lscosmostypes.DelegationModuleAccount:   true,
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
			wasmclient.ProposalHandlers,
			paramsclient.ProposalHandler,
			distributionclient.ProposalHandler,
			upgradeclient.ProposalHandler,
			upgradeclient.CancelProposalHandler,
			ibcclient.UpdateClientProposalHandler,
			ibcclient.UpgradeProposalHandler,
			lscosmosclient.RegisterHostChainProposalHandler,
			lscosmosclient.MinDepositAndFeeChangeProposalHandler,
			lscosmosclient.PstakeFeeAddressChangeProposalHandler,
			lscosmosclient.AllowListValidatorSetChangeProposalHandler,
		)...,
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	wasm.AppModuleBasic{},
	halving.AppModuleBasic{},
	ica.AppModuleBasic{},
	epochs.AppModuleBasic{},
	interchainquery.AppModuleBasic{},
	ibchooker.AppModuleBasic{},
	lscosmos.AppModuleBasic{},
)

var (
	_ simapp.App                          = (*Application)(nil)
	_ servertypes.Application             = (*Application)(nil)
	_ servertypes.ApplicationQueryService = (*Application)(nil)
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
	ICAHostKeeper         *icahostkeeper.Keeper
	EvidenceKeeper        *evidencekeeper.Keeper
	TransferKeeper        *ibctransferkeeper.Keeper
	FeegrantKeeper        *feegrantkeeper.Keeper
	AuthzKeeper           *authzkeeper.Keeper
	HalvingKeeper         *halving.Keeper
	WasmKeeper            *wasm.Keeper
	EpochsKeeper          *epochskeeper.Keeper
	ICAControllerKeeper   *icacontrollerkeeper.Keeper
	LSCosmosKeeper        *lscosmoskeeper.Keeper
	InterchainQueryKeeper *interchainquerykeeper.Keeper
	TransferHooksKeeper   *ibchookerkeeper.Keeper

	moduleManager     *module.Manager
	configurator      module.Configurator
	simulationManager *module.SimulationManager

	// make scoped keepers public for test purposes
	ScopedIBCKeeper           capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper      capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper       capabilitykeeper.ScopedKeeper
	ScopedICAControllerKeeper capabilitykeeper.ScopedKeeper
	ScopedLSCosmosKeeper      capabilitykeeper.ScopedKeeper
	ScopedWasmKeeper          capabilitykeeper.ScopedKeeper
}

func NewApplication(
	applicationName string,
	encodingConfiguration appparams.EncodingConfig,
	moduleAccountPermissions map[string][]string,
	logger tendermintlog.Logger,
	db tendermintdb.DB,
	traceStore io.Writer,
	loadLatest bool,
	invCheckPeriod uint,
	skipUpgradeHeights map[int64]bool,
	home string,
	enabledProposals []wasm.ProposalType,
	applicationOptions servertypes.AppOptions,
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
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distributiontypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		feegrant.StoreKey, authzkeeper.StoreKey, icahosttypes.StoreKey, halving.StoreKey, wasm.StoreKey,
		icacontrollertypes.StoreKey, epochstypes.StoreKey, lscosmostypes.StoreKey, interchainquerytypes.StoreKey,
	)

	transientStoreKeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memoryKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &Application{
		BaseApp:           baseApp,
		legacyAmino:       legacyAmino,
		applicationCodec:  applicationCodec,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
	}

	paramsKeeper := initParamsKeeper(
		applicationCodec,
		legacyAmino,
		keys[paramstypes.StoreKey],
		transientStoreKeys[paramstypes.TStoreKey],
	)
	app.ParamsKeeper = &paramsKeeper
	app.BaseApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	app.CapabilityKeeper = capabilitykeeper.NewKeeper(applicationCodec, keys[capabilitytypes.StoreKey], memoryKeys[capabilitytypes.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedWasmKeeper := app.CapabilityKeeper.ScopeToModule(wasm.ModuleName)
	scopedICAControllerKeeper := app.CapabilityKeeper.ScopeToModule(icacontrollertypes.SubModuleName)
	scopedLSCosmosKeeper := app.CapabilityKeeper.ScopeToModule(lscosmostypes.ModuleName)
	app.CapabilityKeeper.Seal()

	accountKeeper := authkeeper.NewAccountKeeper(
		applicationCodec,
		keys[authtypes.StoreKey],
		app.GetSubspace(authtypes.ModuleName),
		authtypes.ProtoBaseAccount,
		moduleAccountPermissions,
	)
	app.AccountKeeper = &accountKeeper

	blockedModuleAddrs := make(map[string]bool)
	for moduleAccount := range moduleAccountPermissions {
		blockedModuleAddrs[authtypes.NewModuleAddress(moduleAccount).String()] = true
	}
	sendCoinBlockedAddrs := make(map[string]bool)
	for acc := range moduleAccountPermissions {
		sendCoinBlockedAddrs[authtypes.NewModuleAddress(acc).String()] = !receiveAllowedMAcc[acc]
	}

	bankKeeper := bankkeeper.NewBaseKeeper(
		applicationCodec,
		keys[banktypes.StoreKey],
		app.AccountKeeper,
		app.GetSubspace(banktypes.ModuleName),
		sendCoinBlockedAddrs,
	)
	app.BankKeeper = &bankKeeper

	authzKeeper := authzkeeper.NewKeeper(
		keys[authzkeeper.StoreKey],
		applicationCodec,
		app.BaseApp.MsgServiceRouter(),
	)
	app.AuthzKeeper = &authzKeeper

	feegrantKeeper := feegrantkeeper.NewKeeper(
		applicationCodec,
		keys[feegrant.StoreKey],
		app.AccountKeeper,
	)
	app.FeegrantKeeper = &feegrantKeeper

	stakingKeeper := stakingkeeper.NewKeeper(
		applicationCodec,
		keys[stakingtypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(stakingtypes.ModuleName),
	)

	mintKeeper := mintkeeper.NewKeeper(
		applicationCodec,
		keys[minttypes.StoreKey],
		app.GetSubspace(minttypes.ModuleName),
		&stakingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)
	app.MintKeeper = &mintKeeper

	distributionKeeper := distributionkeeper.NewKeeper(
		applicationCodec,
		keys[distributiontypes.StoreKey],
		app.GetSubspace(distributiontypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		authtypes.FeeCollectorName,
		blockedModuleAddrs,
	)
	app.DistributionKeeper = &distributionKeeper

	slashingKeeper := slashingkeeper.NewKeeper(
		applicationCodec,
		keys[slashingtypes.StoreKey],
		&stakingKeeper,
		app.GetSubspace(slashingtypes.ModuleName),
	)
	app.SlashingKeeper = &slashingKeeper

	crisisKeeper := crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName),
		invCheckPeriod,
		app.BankKeeper,
		authtypes.FeeCollectorName,
	)
	app.CrisisKeeper = &crisisKeeper

	upgradeKeeper := upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		keys[upgradetypes.StoreKey],
		applicationCodec,
		home,
		app.BaseApp,
	)
	app.UpgradeKeeper = &upgradeKeeper

	halvingKeeper := halving.NewKeeper(
		keys[halving.StoreKey],
		app.GetSubspace(halving.DefaultParamspace),
		app.MintKeeper,
	)
	app.HalvingKeeper = &halvingKeeper

	app.StakingKeeper = stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistributionKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	epochsKeeper := *epochskeeper.NewKeeper(
		applicationCodec,
		keys[epochstypes.StoreKey],
	)

	app.IBCKeeper = ibckeeper.NewKeeper(
		applicationCodec,
		keys[ibchost.StoreKey],
		app.GetSubspace(ibchost.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	transferKeeper := ibctransferkeeper.NewKeeper(
		applicationCodec,
		keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
	)
	app.TransferKeeper = &transferKeeper

	transferModule := transfer.NewAppModule(*app.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(*app.TransferKeeper)

	icaHostKeeper := icahostkeeper.NewKeeper(
		applicationCodec,
		keys[icahosttypes.StoreKey],
		app.GetSubspace(icahosttypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
	)
	app.ICAHostKeeper = &icaHostKeeper

	icaControllerKeeper := icacontrollerkeeper.NewKeeper(
		applicationCodec, keys[icacontrollertypes.StoreKey],
		app.GetSubspace(icacontrollertypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper, // may be replaced with middleware such as ics29 fee
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		scopedICAControllerKeeper, app.MsgServiceRouter(),
	)
	app.ICAControllerKeeper = &icaControllerKeeper

	icaModule := ica.NewAppModule(app.ICAControllerKeeper, app.ICAHostKeeper)
	icaHostIBCModule := icahost.NewIBCModule(*app.ICAHostKeeper)

	interchainQueryKeeper := interchainquerykeeper.NewKeeper(
		applicationCodec,
		keys[interchainquerytypes.StoreKey],
		app.IBCKeeper,
	)
	app.InterchainQueryKeeper = &interchainQueryKeeper

	interchainQueryModule := interchainquery.NewAppModule(applicationCodec, *app.InterchainQueryKeeper)

	lsCosmosKeeper := lscosmoskeeper.NewKeeper(
		applicationCodec,
		keys[lscosmostypes.StoreKey],
		memoryKeys[lscosmostypes.MemStoreKey],
		app.GetSubspace(lscosmostypes.ModuleName),
		app.BankKeeper,
		app.AccountKeeper,
		epochsKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.TransferKeeper,
		app.ICAControllerKeeper,
		app.InterchainQueryKeeper,
		scopedLSCosmosKeeper,
		app.MsgServiceRouter(),
	)
	app.LSCosmosKeeper = &lsCosmosKeeper

	err := app.InterchainQueryKeeper.SetCallbackHandler(lscosmostypes.ModuleName, app.LSCosmosKeeper.CallbackHandler())
	if err != nil {
		panic(err)
	}

	app.EpochsKeeper = epochsKeeper.SetHooks(
		epochstypes.NewMultiEpochHooks(app.LSCosmosKeeper.NewEpochHooks()),
	)
	// Information will flow: ibc-port -> icaController -> lscosmos.
	lscosmosModule := lscosmos.NewAppModule(applicationCodec, *app.LSCosmosKeeper, app.AccountKeeper, app.BankKeeper)
	icaControllerIBCModule := icacontroller.NewIBCModule(*app.ICAControllerKeeper, lscosmosModule)

	ibcTransferHooksKeeper := ibchookerkeeper.NewKeeper()
	app.TransferHooksKeeper = ibcTransferHooksKeeper.SetHooks(ibchookertypes.NewMultiStakingHooks(app.LSCosmosKeeper.NewIBCTransferHooks()))
	ibcTransferHooksMiddleware := ibchooker.NewAppModule(*app.TransferHooksKeeper, transferIBCModule)

	evidenceKeeper := evidencekeeper.NewKeeper(
		applicationCodec,
		keys[evidencetypes.StoreKey],
		app.StakingKeeper,
		app.SlashingKeeper,
	)
	app.EvidenceKeeper = evidenceKeeper

	wasmDir := filepath.Join(home, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(applicationOptions)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}

	// The last arguments can contain custom message handlers, and custom query handlers,
	// if we want to allow any custom callbacks
	supportedFeatures := "iterator,staking,stargate"
	wasmKeeper := wasm.NewKeeper(
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
	app.WasmKeeper = &wasmKeeper

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibctypes.NewRouter()
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostIBCModule).
		AddRoute(ibctransfertypes.ModuleName, ibcTransferHooksMiddleware).
		AddRoute(icacontrollertypes.SubModuleName, icaControllerIBCModule).
		AddRoute(lscosmostypes.ModuleName, icaControllerIBCModule).
		AddRoute(wasm.ModuleName, wasm.NewIBCHandler(app.WasmKeeper, app.IBCKeeper.ChannelKeeper))
	app.IBCKeeper.SetRouter(ibcRouter)

	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(
		govtypes.RouterKey,
		govtypes.ProposalHandler,
	).AddRoute(
		paramsproposal.RouterKey,
		params.NewParamChangeProposalHandler(*app.ParamsKeeper),
	).AddRoute(
		distributiontypes.RouterKey,
		distribution.NewCommunityPoolSpendProposalHandler(*app.DistributionKeeper),
	).AddRoute(
		upgradetypes.RouterKey,
		upgrade.NewSoftwareUpgradeProposalHandler(*app.UpgradeKeeper),
	).AddRoute(
		ibcclienttypes.RouterKey, ibccoreclient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper),
	).AddRoute(lscosmostypes.RouterKey, lscosmos.NewLSCosmosProposalHandler(*app.LSCosmosKeeper))

	if len(enabledProposals) != 0 {
		govRouter.AddRoute(wasm.RouterKey, wasm.NewWasmProposalHandler(app.WasmKeeper, enabledProposals))
	}
	govKeeper := govkeeper.NewKeeper(
		applicationCodec,
		keys[govtypes.StoreKey],
		app.GetSubspace(govtypes.ModuleName),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		govRouter,
	)
	app.GovKeeper = &govKeeper

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
		auth.NewAppModule(applicationCodec, *app.AccountKeeper, nil),
		vesting.NewAppModule(*app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(applicationCodec, *app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(applicationCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(applicationCodec, *app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(applicationCodec, *app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(applicationCodec, *app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper),
		distribution.NewAppModule(applicationCodec, *app.DistributionKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper),
		staking.NewAppModule(applicationCodec, *app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(*app.UpgradeKeeper),
		evidence.NewAppModule(*app.EvidenceKeeper),
		feegrantmodule.NewAppModule(applicationCodec, app.AccountKeeper, app.BankKeeper, *app.FeegrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(applicationCodec, *app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(*app.ParamsKeeper),
		halving.NewAppModule(applicationCodec, *app.HalvingKeeper),
		transferModule,
		ibcTransferHooksMiddleware,
		icaModule,
		wasm.NewAppModule(applicationCodec, app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		epochs.NewAppModule(*app.EpochsKeeper),
		interchainQueryModule,
		lscosmosModule,
	)

	app.moduleManager.SetOrderBeginBlockers(
		upgradetypes.ModuleName,
		epochstypes.ModuleName,
		capabilitytypes.ModuleName,
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distributiontypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		halving.ModuleName,
		wasm.ModuleName,
		ibchookertypes.ModuleName,
		interchainquerytypes.ModuleName,
		lscosmostypes.ModuleName,
	)
	app.moduleManager.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distributiontypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		halving.ModuleName,
		wasm.ModuleName,
		epochstypes.ModuleName,
		ibchookertypes.ModuleName,
		interchainquerytypes.ModuleName,
		lscosmostypes.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.moduleManager.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		banktypes.ModuleName,
		distributiontypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		authtypes.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		halving.ModuleName,
		wasm.ModuleName,
		epochstypes.ModuleName,
		ibchookertypes.ModuleName,
		interchainquerytypes.ModuleName,
		lscosmostypes.ModuleName,
	)

	app.moduleManager.RegisterInvariants(app.CrisisKeeper)
	app.moduleManager.RegisterRoutes(app.BaseApp.Router(), app.BaseApp.QueryRouter(), encodingConfiguration.Amino)
	app.configurator = module.NewConfigurator(app.applicationCodec, app.BaseApp.MsgServiceRouter(), app.BaseApp.GRPCQueryRouter())
	app.moduleManager.RegisterServices(app.configurator)

	simulationManager := module.NewSimulationManager(
		auth.NewAppModule(applicationCodec, *app.AccountKeeper, authsimulation.RandomGenesisAccounts),
		bank.NewAppModule(applicationCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(applicationCodec, *app.CapabilityKeeper),
		gov.NewAppModule(applicationCodec, *app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(applicationCodec, *app.MintKeeper, app.AccountKeeper),
		staking.NewAppModule(applicationCodec, *app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		distribution.NewAppModule(applicationCodec, *app.DistributionKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		slashing.NewAppModule(applicationCodec, *app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		params.NewAppModule(*app.ParamsKeeper),
		halving.NewAppModule(applicationCodec, *app.HalvingKeeper),
		authzmodule.NewAppModule(applicationCodec, *app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		feegrantmodule.NewAppModule(applicationCodec, app.AccountKeeper, app.BankKeeper, *app.FeegrantKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		transferModule,
		interchainQueryModule,
		lscosmosModule,
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
			wasmkeeper.NewWasmSnapshotter(app.CommitMultiStore(), app.WasmKeeper),
		)
		if err != nil {
			panic(fmt.Errorf("failed to register snapshot extension: %s", err))
		}
	}

	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx sdk.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			ctx.Logger().Info("start to run upgrade migration...")

			//add more upgrade instructions
			ctx.Logger().Info("running revert of tombstoning")
			err := upgrades.RevertCosTombstoning(ctx, app.SlashingKeeper, app.MintKeeper, app.BankKeeper, app.StakingKeeper)
			if err != nil {
				panic(fmt.Sprintf("failed to revert tombstoning: %s", err))
			}

			//err = upgrades.MintPstakeTokens(ctx, app.LSCosmosKeeper)
			//if err != nil {
			//	panic(fmt.Sprintf("failed to mint pstake tokens: %s", err))
			//}

			ctx.Logger().Info("start to run module migrations...")
			newVM, err := app.moduleManager.RunMigrations(ctx, app.configurator, fromVM)
			if err != nil {
				return nil, err
			}
			return newVM, nil
		},
	)

	// commenting out for this upgrade, since not needed
	//upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	//if err != nil {
	//	panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	//}

	//if upgradeInfo.Name == UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
	//	storeUpgrades := storetypes.StoreUpgrades{
	//		Added: []string{},
	//	}
	//
	//	// configure store loader that checks if version == upgradeHeight and applies store upgrades
	//	app.BaseApp.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	//}

	if loadLatest {
		if err := app.BaseApp.LoadLatestVersion(); err != nil {
			tendermintos.Exit(err.Error())
		}
		ctx := app.BaseApp.NewUncachedContext(true, tendermintproto.Header{})

		// Initialize pinned codes in wasmvm as they are not persisted there
		if err := app.WasmKeeper.InitializePinnedCodes(ctx); err != nil {
			tendermintos.Exit(fmt.Sprintf("failed initialize pinned codes %s", err))
		}
	}
	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	app.ScopedICAHostKeeper = scopedICAHostKeeper
	app.ScopedWasmKeeper = scopedWasmKeeper
	app.ScopedICAControllerKeeper = scopedICAControllerKeeper
	app.ScopedLSCosmosKeeper = scopedLSCosmosKeeper

	return app
}

func (app *Application) CreateUpgradeHandler(mm *module.Manager, configurator module.Configurator) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, _ upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		return mm.RunMigrations(ctx, configurator, vm)
	}
}

func (app *Application) ApplicationCodec() codec.Codec {
	return app.applicationCodec
}

func (app *Application) Name() string {
	return app.BaseApp.Name()
}

func (app *Application) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app *Application) BeginBlocker(ctx sdk.Context, req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	return app.moduleManager.BeginBlock(ctx, req)
}

func (app *Application) EndBlocker(ctx sdk.Context, req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return app.moduleManager.EndBlock(ctx, req)
}

func (app *Application) InitChainer(ctx sdk.Context, req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	var genesisState GenesisState
	if err := tendermintjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	app.UpgradeKeeper.SetModuleVersionMap(ctx, app.moduleManager.GetVersionMap())

	return app.moduleManager.InitGenesis(ctx, app.applicationCodec, genesisState)
}

func (app *Application) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string) (servertypes.ExportedApp, error) {
	context := app.BaseApp.NewContext(true, tendermintproto.Header{Height: app.BaseApp.LastBlockHeight()})

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

		app.StakingKeeper.IterateValidators(context, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {
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

		app.StakingKeeper.IterateValidators(context, func(_ int64, val stakingtypes.ValidatorI) (stop bool) {

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

		app.StakingKeeper.IterateRedelegations(context, func(_ int64, redelegation stakingtypes.Redelegation) (stop bool) {
			for i := range redelegation.Entries {
				redelegation.Entries[i].CreationHeight = 0
			}
			app.StakingKeeper.SetRedelegation(context, redelegation)
			return false
		})

		app.StakingKeeper.IterateUnbondingDelegations(context, func(_ int64, unbondingDelegation stakingtypes.UnbondingDelegation) (stop bool) {
			for i := range unbondingDelegation.Entries {
				unbondingDelegation.Entries[i].CreationHeight = 0
			}
			app.StakingKeeper.SetUnbondingDelegation(context, unbondingDelegation)
			return false
		})

		store := context.KVStore(app.keys[stakingtypes.StoreKey])
		kvStoreReversePrefixIterator := sdk.KVStoreReversePrefixIterator(store, stakingtypes.ValidatorsKey)
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
			func(validatorConsAddress sdk.ConsAddress, validatorSigningInfo slashingtypes.ValidatorSigningInfo) (stop bool) {
				validatorSigningInfo.StartHeight = 0
				app.SlashingKeeper.SetValidatorSigningInfo(context, validatorConsAddress, validatorSigningInfo)
				return false
			},
		)
	}

	genesisState := app.moduleManager.ExportGenesis(context, app.applicationCodec)
	applicationState, Error := codec.MarshalJSONIndent(app.legacyAmino, genesisState)

	if Error != nil {
		return servertypes.ExportedApp{}, Error
	}

	validators, err := staking.WriteValidators(context, *app.StakingKeeper)

	return servertypes.ExportedApp{
		AppState:        applicationState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: app.BaseApp.GetConsensusParams(context),
	}, err
}

func (app *Application) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range ModuleAccountPermissions {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *Application) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

func (app *Application) SimulationManager() *module.SimulationManager {
	return app.simulationManager
}

func (app *Application) ListSnapshots(snapshots abcitypes.RequestListSnapshots) abcitypes.ResponseListSnapshots {
	return app.BaseApp.ListSnapshots(snapshots)
}

func (app *Application) OfferSnapshot(snapshot abcitypes.RequestOfferSnapshot) abcitypes.ResponseOfferSnapshot {
	return app.BaseApp.OfferSnapshot(snapshot)
}

func (app *Application) LoadSnapshotChunk(chunk abcitypes.RequestLoadSnapshotChunk) abcitypes.ResponseLoadSnapshotChunk {
	return app.BaseApp.LoadSnapshotChunk(chunk)
}

func (app *Application) ApplySnapshotChunk(chunk abcitypes.RequestApplySnapshotChunk) abcitypes.ResponseApplySnapshotChunk {
	return app.BaseApp.ApplySnapshotChunk(chunk)
}

func (app *Application) RegisterGRPCServer(server grpc.Server) {
	app.BaseApp.RegisterGRPCServer(server)
}

func (app *Application) RegisterAPIRoutes(apiServer *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiServer.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiServer.Router)
	// Register legacy tx routes.
	authrest.RegisterTxRoutes(clientCtx, apiServer.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiServer.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiServer.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiServer.GRPCGatewayRouter)
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

func (app *Application) RegisterTxService(clientContect client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientContect, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *Application) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}
func (app *Application) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}
func (app *Application) LoadHeight(height int64) error {
	return app.BaseApp.LoadVersion(height)
}

// initParamsKeeper init params keeper and its subspaces.
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distributiontypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(halving.DefaultParamspace)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(wasm.ModuleName)
	paramsKeeper.Subspace(icacontrollertypes.SubModuleName)
	paramsKeeper.Subspace(lscosmostypes.ModuleName)
	paramsKeeper.Subspace(interchainquerytypes.ModuleName)

	return paramsKeeper
}
