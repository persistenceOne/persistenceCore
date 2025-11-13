/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package app

import (
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"reflect"
	"strings"

	"cosmossdk.io/client/v2/autocli"
	clienthelpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/core/appmodule"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	"cosmossdk.io/log"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	abci "github.com/cometbft/cometbft/abci/types"
	tendermintjson "github.com/cometbft/cometbft/libs/json"
	tendermintos "github.com/cometbft/cometbft/libs/os"
	tendermintproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	tmservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	sigtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	txmodule "github.com/cosmos/cosmos-sdk/x/auth/tx/config"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/spf13/cast"

	"github.com/persistenceOne/persistenceCore/v16/app/constants"
	"github.com/persistenceOne/persistenceCore/v16/app/keepers"
	"github.com/persistenceOne/persistenceCore/v16/app/upgrades"
	v1600rc0 "github.com/persistenceOne/persistenceCore/v16/app/upgrades/testnet/v16.0.0-rc0"
)

var (
	DefaultNodeHome string
	Upgrades        = []upgrades.Upgrade{v1600rc0.Upgrade}
)

var (
	_ runtime.AppI            = (*Application)(nil)
	_ servertypes.Application = (*Application)(nil)
)

func init() {
	var err error
	DefaultNodeHome, err = clienthelpers.GetNodeHomeDirectory(".persistenceCore")
	if err != nil {
		panic(err)
	}
}

type Application struct {
	*baseapp.BaseApp
	*keepers.AppKeepers

	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry types.InterfaceRegistry

	ModuleBasicManager module.BasicManager
	moduleManager      *module.Manager
	configurator       module.Configurator
	simulationManager  *module.SimulationManager
}

func NewApplication(
	logger log.Logger,
	db dbm.DB,
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

	// may be in future - vote extensions

	baseApp := baseapp.NewBaseApp(
		constants.AppName,
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
	nodeConfig, err := wasm.ReadNodeConfig(applicationOptions)
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
		logger,
	)

	enabledSignModes := append(tx.DefaultSignModes, sigtypes.SignMode_SIGN_MODE_TEXTUAL)
	txConfigOpts := tx.ConfigOptions{
		EnabledSignModes:           enabledSignModes,
		TextualCoinMetadataQueryFn: txmodule.NewBankKeeperCoinMetadataQueryFn(app.BankKeeper),
	}
	txConfig, err = tx.NewTxConfigWithOptions(
		appCodec,
		txConfigOpts,
	)
	if err != nil {
		panic(err)
	}
	app.txConfig = txConfig

	app.moduleManager = module.NewManager(appModules(app, app.appCodec, app.txConfig)...)

	app.ModuleBasicManager = module.NewBasicManagerFromManager(app.moduleManager,
		map[string]module.AppModuleBasic{
			genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			govtypes.ModuleName: gov.NewAppModuleBasic(
				[]govclient.ProposalHandler{},
			),
		})
	// register interfaces, amino codecs
	AppendModuleCodecs(app.legacyAmino, app.interfaceRegistry, app.ModuleBasicManager)

	app.moduleManager.SetOrderPreBlockers(orderPreBlockers()...)
	app.moduleManager.SetOrderBeginBlockers(orderBeginBlockers()...)
	app.moduleManager.SetOrderEndBlockers(orderEndBlockers()...)
	app.moduleManager.SetOrderInitGenesis(orderInitGenesis()...)
	app.moduleManager.SetOrderExportGenesis(orderExportGenesis()...)
	//app.moduleManager.SetOrderMigrations()

	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	err = app.moduleManager.RegisterServices(app.configurator)
	if err != nil {
		panic(err)
	}

	app.simulationManager = module.NewSimulationManagerFromAppModules(
		app.moduleManager.Modules,
		overrideSimulationModules(app, app.appCodec),
	)
	app.simulationManager.RegisterStoreDecoders()

	app.MountKVStores(app.GetKVStoreKey())
	app.MountTransientStores(app.GetTransientStoreKey())
	app.MountMemoryStores(app.GetMemoryStoreKey())

	app.setupAnteHandler(nodeConfig)
	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.setPostHandler()

	// setup postHandler in this method
	// app.setupPostHandler()
	app.setupUpgradeHandlers(Upgrades)
	app.setupUpgradeStoreLoaders(Upgrades)

	app.registerGRPCServices()

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

func (app *Application) setupAnteHandler(nodeConfig wasmtypes.NodeConfig) {

	anteOptions := HandlerOptions{
		HandlerOptions: ante.HandlerOptions{
			AccountKeeper:   app.AccountKeeper,
			BankKeeper:      app.BankKeeper,
			FeegrantKeeper:  app.FeegrantKeeper,
			SignModeHandler: app.TxConfig().SignModeHandler(),
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			SigVerifyOptions: []ante.SigVerificationDecoratorOption{
				ante.WithUnorderedTxGasCost(ante.DefaultUnorderedTxGasCost),
				ante.WithMaxUnorderedTxTimeoutDuration(ante.DefaultMaxTimeoutDuration),
			},
		},
		IBCKeeper:             app.IBCKeeper,
		WasmKeeper:            app.WasmKeeper,
		NodeConfig:            &nodeConfig,
		TXCounterStoreService: runtime.NewKVStoreService(app.GetKVStoreKey()[wasmtypes.StoreKey]),

		TxDecoder: app.txConfig.TxDecoder(),
		TxEncoder: app.txConfig.TxEncoder(),

		FeeDenomsWhitelist: app.GetFeeDenomsWhitelist(),
	}
	anteHandler, err := NewAnteHandler(anteOptions)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	app.BaseApp.SetAnteHandler(anteHandler)
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
		return constants.BondDenom
	} else if strings.HasPrefix(chainID, "test-core-") {
		return constants.BondDenom
	}

	return "stake"
}

func (app *Application) GetFeeDenomsWhitelist() []string {
	chainID := app.ChainID()

	if strings.HasPrefix(chainID, "core-") {
		return constants.FeeDenomsWhitelistMainnet
	} else if strings.HasPrefix(chainID, "test-core-") {
		return constants.FeeDenomsWhitelistTestnet
	}

	// Allow all denoms for random chain
	return []string{} // empty list => allow all
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

func (app *Application) setPostHandler() {
	postHandler, err := posthandler.NewPostHandler(
		posthandler.HandlerOptions{},
	)
	if err != nil {
		panic(err)
	}

	app.SetPostHandler(postHandler)
}

func (app *Application) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	return app.moduleManager.PreBlock(ctx)
}

func (app *Application) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.moduleManager.BeginBlock(ctx)
}

func (app *Application) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.moduleManager.EndBlock(ctx)
}

func (app *Application) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	if err := tendermintjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	err := app.UpgradeKeeper.SetModuleVersionMap(ctx, app.moduleManager.GetVersionMap())
	if err != nil {
		panic(err)
	}

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
	app.ModuleBasicManager.RegisterGRPCGatewayRoutes(clientCtx, apiServer.GRPCGatewayRouter)

	// register swagger API from root so that other applications can override easily
	if err := server.RegisterSwaggerAPI(apiServer.ClientCtx, apiServer.Router, apiConfig.Swagger); err != nil {
		panic(err)
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

func (app *Application) RegisterTxService(clientContect client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientContect, app.BaseApp.Simulate, app.interfaceRegistry)
}

func (app *Application) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.Query)
}
func (app *Application) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// AutoCliOpts returns the autocli options for the app.
func (app *Application) AutoCliOpts() autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule, 0)
	for _, m := range app.moduleManager.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservices.ExtractAutoCLIOptions(app.moduleManager.Modules),
		AddressCodec:          address.NewBech32Codec(constants.Bech32PrefixAccAddr),
		ValidatorAddressCodec: address.NewBech32Codec(constants.Bech32PrefixValAddr),
		ConsensusAddressCodec: address.NewBech32Codec(constants.Bech32PrefixConsAddr),
	}
}

func (app *Application) LoadHeight(height int64) error {
	return app.BaseApp.LoadVersion(height)
}

// DefaultGenesis returns a default genesis from the registered AppModuleBasic's.
func (app *Application) DefaultGenesis() map[string]json.RawMessage {
	return app.ModuleBasicManager.DefaultGenesis(app.appCodec)
}

func SendCoinBlockedAddrs() map[string]bool {
	sendCoinBlockedAddrs := make(map[string]bool)
	for acc := range ModuleAccountPermissions {
		sendCoinBlockedAddrs[authtypes.NewModuleAddress(acc).String()] = !receiveAllowedMAcc[acc]
	}
	return sendCoinBlockedAddrs
}
