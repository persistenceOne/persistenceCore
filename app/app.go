/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package app

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	stdlog "log"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	tendermintdb "github.com/cometbft/cometbft-db"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	tendermintjson "github.com/cometbft/cometbft/libs/json"
	tendermintlog "github.com/cometbft/cometbft/libs/log"
	tendermintos "github.com/cometbft/cometbft/libs/os"
	tendermintproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/gorilla/mux"
	pobabci "github.com/skip-mev/pob/abci"
	"github.com/skip-mev/pob/mempool"
	"github.com/spf13/cast"

	"github.com/persistenceOne/persistenceCore/v11/app/keepers"
	"github.com/persistenceOne/persistenceCore/v11/app/upgrades"
	v11_11_0 "github.com/persistenceOne/persistenceCore/v11/app/upgrades/v11.11.0"
	"github.com/persistenceOne/persistenceCore/v11/client/docs"
)

var (
	DefaultNodeHome string
	Upgrades        = []upgrades.Upgrade{v11_11_0.Upgrade}
	ModuleBasics    = module.NewBasicManager(keepers.AppModuleBasics...)
)

var (
	_ runtime.AppI            = (*Application)(nil)
	_ servertypes.Application = (*Application)(nil)
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
	*keepers.AppKeepers

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	moduleManager     *module.Manager
	configurator      module.Configurator
	simulationManager *module.SimulationManager

	// override handler for CheckTx for POB
	checkTxHandler pobabci.CheckTx
}

func NewApplication(
	logger tendermintlog.Logger,
	db tendermintdb.DB,
	traceStore io.Writer,
	loadLatest bool,
	applicationOptions servertypes.AppOptions,
	wasmOpts []wasm.Option,
	baseAppOptions ...func(*baseapp.BaseApp),
) *Application {
	encodingConfiguration := MakeEncodingConfig()

	appCodec := encodingConfiguration.Codec
	legacyAmino := encodingConfiguration.Amino
	interfaceRegistry := encodingConfiguration.InterfaceRegistry
	txConfig := encodingConfiguration.TxConfig

	baseApp := baseapp.NewBaseApp(
		AppName,
		logger,
		db,
		txConfig.TxDecoder(),
		baseAppOptions...,
	)
	baseApp.SetCommitMultiStoreTracer(traceStore)
	baseApp.SetVersion(version.Version)
	baseApp.SetInterfaceRegistry(interfaceRegistry)
	baseApp.SetTxEncoder(txConfig.TxEncoder())

	homePath := cast.ToString(applicationOptions.Get(flags.FlagHome))
	wasmDir := filepath.Join(homePath, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(applicationOptions)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}

	app := &Application{
		BaseApp:           baseApp,
		legacyAmino:       legacyAmino,
		appCodec:          appCodec,
		txConfig:          txConfig,
		interfaceRegistry: interfaceRegistry,
	}

	// Setup keepers
	app.AppKeepers = keepers.NewAppKeeper(
		appCodec,
		baseApp,
		legacyAmino,
		ModuleAccountPermissions,
		SendCoinBlockedAddrs(),
		applicationOptions,
		wasmDir,
		wasmOpts,
		Bech32MainPrefix,
	)

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(applicationOptions.Get(crisis.FlagSkipGenesisInvariants))

	app.moduleManager = module.NewManager(appModules(app, encodingConfiguration, skipGenesisInvariants)...)

	app.moduleManager.SetOrderBeginBlockers(orderBeginBlockers()...)
	app.moduleManager.SetOrderEndBlockers(orderEndBlockers()...)
	app.moduleManager.SetOrderInitGenesis(orderInitGenesis()...)
	app.moduleManager.SetOrderExportGenesis(orderInitGenesis()...)

	app.moduleManager.RegisterInvariants(app.CrisisKeeper)
	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	app.moduleManager.RegisterServices(app.configurator)

	app.simulationManager = module.NewSimulationManagerFromAppModules(
		app.moduleManager.Modules,
		overrideSimulationModules(app, encodingConfiguration, skipGenesisInvariants),
	)
	app.simulationManager.RegisterStoreDecoders()

	app.registerGRPCServices()

	app.MountKVStores(app.GetKVStoreKey())
	app.MountTransientStores(app.GetTransientStoreKey())
	app.MountMemoryStores(app.GetMemoryStoreKey())

	app.setupPOBAndAnteHandler(wasmConfig)
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// setup postHandler in this method
	// app.setupPostHandler()
	app.setupUpgradeHandlers(Upgrades)
	app.setupUpgradeStoreLoaders(Upgrades)

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

	return app
}

func (app *Application) setupPOBAndAnteHandler(wasmConfig wasmtypes.WasmConfig) {
	// Set POB's mempool into the app.
	mempool := mempool.NewAuctionMempool(app.txConfig.TxDecoder(), app.txConfig.TxEncoder(), 0, mempool.NewDefaultAuctionFactory(app.txConfig.TxDecoder()))
	app.BaseApp.SetMempool(mempool)

	anteOptions := HandlerOptions{
		HandlerOptions: ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			FeegrantKeeper:  app.FeegrantKeeper,
			SignModeHandler: app.txConfig.SignModeHandler(),
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
		IBCKeeper:         app.IBCKeeper,
		WasmConfig:        &wasmConfig,
		TXCounterStoreKey: app.GetKVStoreKey()[wasm.StoreKey],

		BuilderKeeper: app.BuilderKeeper,
		Mempool:       mempool,
		TxDecoder:     app.txConfig.TxDecoder(),
		TxEncoder:     app.txConfig.TxEncoder(),

		FeeDenomsWhitelist: app.GetFeeDenomsWhitelist(),
	}
	anteHandler, err := NewAnteHandler(anteOptions)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	// Set the proposal handlers on the BaseApp along with the custom antehandler.
	proposalHandlers := pobabci.NewProposalHandler(
		mempool,
		app.BaseApp.Logger(),
		anteHandler,
		anteOptions.TxEncoder,
		anteOptions.TxDecoder,
	)
	app.BaseApp.SetPrepareProposal(proposalHandlers.PrepareProposalHandler())
	app.BaseApp.SetProcessProposal(proposalHandlers.ProcessProposalHandler())
	app.BaseApp.SetAnteHandler(anteHandler)

	chainID := app.ChainID()
	app.BaseApp.Logger().Info("using BaseApp chainID for POB ABCI", "chainID", chainID)

	// Set the custom CheckTx handler on BaseApp.
	checkTxHandler := pobabci.NewCheckTxHandler(
		app.BaseApp,
		app.txConfig.TxDecoder(),
		mempool,
		anteHandler,
		chainID,
	)

	app.SetCheckTx(checkTxHandler.CheckTx())
}

func (app *Application) registerGRPCServices() {
	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.moduleManager.Modules))

	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)
}

// ChainID gets chainID from private fields of BaseApp
// Should be removed once SDK 0.50.x will be adopted
func (app *Application) ChainID() string {
	field := reflect.ValueOf(app.BaseApp).Elem().FieldByName("chainID")
	return field.String()
}

// GetChainBondDenom returns expected chain bond denom.
func (app *Application) GetChainBondDenom() string {
	chainID := app.ChainID()

	if strings.HasPrefix(chainID, "core-") {
		return BondDenom
	} else if strings.HasPrefix(chainID, "test-core-") {
		return BondDenom
	}

	return "stake"
}

func (app *Application) GetFeeDenomsWhitelist() []string {
	chainID := app.ChainID()

	if strings.HasPrefix(chainID, "core-") {
		return FeeDenomsWhitelistMainnet
	} else if strings.HasPrefix(chainID, "test-core-") {
		return FeeDenomsWhitelistTestnet
	}

	// Allow all denoms for random chain
	return []string{} // empty list => allow all
}

// CheckTx will check the transaction with the provided checkTxHandler. We override the default
// handler so that we can verify bid transactions before they are inserted into the mempool.
// With the POB CheckTx, we can verify the bid transaction and all of the bundled transactions
// before inserting the bid transaction into the mempool.
func (app *Application) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	return app.checkTxHandler(req)
}

// SetCheckTx sets the checkTxHandler for the app.
func (app *Application) SetCheckTx(handler pobabci.CheckTx) {
	app.checkTxHandler = handler
}

func (app *Application) Name() string {
	return app.BaseApp.Name()
}

func (app *Application) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns Application's InterfaceRegistry
func (app *Application) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// TxConfig returns Application's TxConfig
func (app *Application) TxConfig() client.TxConfig {
	return app.txConfig
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

	return app.moduleManager.InitGenesis(ctx, app.appCodec, genesisState)
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

func (app *Application) ModuleManager() *module.Manager {
	return app.moduleManager
}

func (app *Application) SimulationManager() *module.SimulationManager {
	return app.simulationManager
}

func (app *Application) RegisterAPIRoutes(apiServer *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiServer.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiServer.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiServer.GRPCGatewayRouter)
	// Register node gRPC service for grpc-gateway.
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiServer.GRPCGatewayRouter)
	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiServer.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if apiConfig.Swagger {
		RegisterSwaggerAPI(apiServer.Router)
	}
}

func (app *Application) setupUpgradeHandlers(upgradeversions []upgrades.Upgrade) {
	for _, upgrade := range upgradeversions {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(upgrades.UpgradeHandlerArgs{
				ModuleManager: app.moduleManager,
				Configurator:  app.configurator,
				Keepers:       app.AppKeepers,
				Codec:         app.appCodec,
			}),
		)
	}
}

// configure store loader that checks if version == upgradeHeight and applies store upgrades
func (app *Application) setupUpgradeStoreLoaders(upgradeversions []upgrades.Upgrade) {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, upgrade := range upgradeversions {
		if upgradeInfo.Name == upgrade.UpgradeName {
			storeUpgrades := upgrade.StoreUpgrades
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
		}
	}
}

// PostHandlers are like AnteHandlers (they have the same signature), but they are run after runMsgs.
// One use case for PostHandlers is transaction tips,
// but other use cases like unused gas refund can also be enabled by PostHandlers.
//
// In baseapp, postHandlers are run in the same store branch as `runMsgs`,
// meaning that both `runMsgs` and `postHandler` state will be committed if
// both are successful, and both will be reverted if any of the two fails.
// nolint:unused // post handle is not used for now (as there is no requirement of it)
func (app *Application) setupPostHandler() {
	postDecorators := []sdk.PostDecorator{
		// posthandler.NewTipDecorator(app.BankKeeper),
		// ... add more decorators as needed
	}
	postHandler := sdk.ChainPostDecorators(postDecorators...)
	app.SetPostHandler(postHandler)
}

func RegisterSwaggerAPI(rtr *mux.Router) {
	swaggerDir, err := fs.Sub(docs.SwaggerUI, "swagger-ui")
	if err != nil {
		panic(err)
	}

	staticServer := http.FileServer(http.FS(swaggerDir))
	rtr.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
}

func (app *Application) RegisterTxService(clientContect client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientContect, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *Application) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.Query)
}
func (app *Application) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}
func (app *Application) LoadHeight(height int64) error {
	return app.BaseApp.LoadVersion(height)
}

// DefaultGenesis returns a default genesis from the registered AppModuleBasic's.
func (app *Application) DefaultGenesis() map[string]json.RawMessage {
	return ModuleBasics.DefaultGenesis(app.appCodec)
}

func SendCoinBlockedAddrs() map[string]bool {
	sendCoinBlockedAddrs := make(map[string]bool)
	for acc := range ModuleAccountPermissions {
		sendCoinBlockedAddrs[authtypes.NewModuleAddress(acc).String()] = !receiveAllowedMAcc[acc]
	}
	return sendCoinBlockedAddrs
}
