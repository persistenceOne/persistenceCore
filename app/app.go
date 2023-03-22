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
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/gogo/protobuf/grpc"
	"github.com/gorilla/mux"
	lscosmostypes "github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
	"github.com/rakyll/statik/fs"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tendermintjson "github.com/tendermint/tendermint/libs/json"
	tendermintlog "github.com/tendermint/tendermint/libs/log"
	tendermintos "github.com/tendermint/tendermint/libs/os"
	tendermintproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tendermintdb "github.com/tendermint/tm-db"

	"github.com/persistenceOne/persistenceCore/v7/app/keepers"
	appparams "github.com/persistenceOne/persistenceCore/v7/app/params"
	"github.com/persistenceOne/persistenceCore/v7/app/upgrades"
	v7 "github.com/persistenceOne/persistenceCore/v7/app/upgrades/v7"
)

var (
	DefaultNodeHome string
	Upgrades        = []upgrades.Upgrade{v7.Upgrade}
	ModuleBasics    = module.NewBasicManager(keepers.AppModuleBasics...)
)

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

var receiveAllowedMAcc = map[string]bool{
	lscosmostypes.UndelegationModuleAccount: true,
	lscosmostypes.DelegationModuleAccount:   true,
}

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
	keepers.AppKeepers

	legacyAmino       *codec.LegacyAmino
	applicationCodec  codec.Codec
	interfaceRegistry types.InterfaceRegistry

	moduleManager     *module.Manager
	configurator      module.Configurator
	simulationManager *module.SimulationManager
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

	app := &Application{
		BaseApp:           baseApp,
		legacyAmino:       legacyAmino,
		applicationCodec:  applicationCodec,
		interfaceRegistry: interfaceRegistry,
	}

	// these blocked address will be used in distribution keeper as well
	sendCoinBlockedAddrs := make(map[string]bool)
	for acc := range moduleAccountPermissions {
		sendCoinBlockedAddrs[authtypes.NewModuleAddress(acc).String()] = !receiveAllowedMAcc[acc]
	}

	wasmDir := filepath.Join(home, "wasm")
	wasmConfig, err := wasm.ReadWasmConfig(applicationOptions)
	if err != nil {
		panic(fmt.Sprintf("error while reading wasm config: %s", err))
	}

	// Setup keepers
	appKeepers := keepers.NewAppKeeper(
		applicationCodec,
		baseApp,
		legacyAmino,
		moduleAccountPermissions,
		sendCoinBlockedAddrs,
		skipUpgradeHeights,
		home,
		invCheckPeriod,
		applicationOptions,
		wasmDir,
		enabledProposals,
		wasmOpts,
		Bech32MainPrefix,
	)
	app.AppKeepers = appKeepers

	/****  Module Options ****/
	var skipGenesisInvariants = false

	opt := applicationOptions.Get(crisis.FlagSkipGenesisInvariants)
	if opt, ok := opt.(bool); ok {
		skipGenesisInvariants = opt
	}

	app.moduleManager = module.NewManager(appModules(app, encodingConfiguration, skipGenesisInvariants)...)

	app.moduleManager.SetOrderBeginBlockers(orderBeginBlockers()...)
	app.moduleManager.SetOrderEndBlockers(orderEndBlockers()...)
	app.moduleManager.SetOrderInitGenesis(orderInitGenesis()...)

	app.moduleManager.RegisterInvariants(app.CrisisKeeper)
	app.moduleManager.RegisterRoutes(app.BaseApp.Router(), app.BaseApp.QueryRouter(), encodingConfiguration.Amino)
	app.configurator = module.NewConfigurator(app.applicationCodec, app.BaseApp.MsgServiceRouter(), app.BaseApp.GRPCQueryRouter())
	app.moduleManager.RegisterServices(app.configurator)

	simulationManager := module.NewSimulationManager(simulationModules(app, encodingConfiguration, skipGenesisInvariants)...)

	simulationManager.RegisterStoreDecoders()
	app.simulationManager = simulationManager

	app.BaseApp.MountKVStores(app.GetKVStoreKey())
	app.BaseApp.MountTransientStores(app.GetTransientStoreKey())
	app.BaseApp.MountMemoryStores(app.GetMemoryStoreKey())

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
			TXCounterStoreKey: app.GetKVStoreKey()[wasm.StoreKey],
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	app.BaseApp.SetAnteHandler(anteHandler)
	app.BaseApp.SetInitChainer(app.InitChainer)
	app.BaseApp.SetBeginBlocker(app.moduleManager.BeginBlock)
	app.BaseApp.SetEndBlocker(app.moduleManager.EndBlock)

	// setup postHandler in this method
	// app.setupPostHandler()

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

	app.setupUpgradeHandlers()
	app.setupUpgradeStoreLoaders()

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

		store := context.KVStore(app.GetKVStoreKey()[stakingtypes.StoreKey])
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

func (app *Application) setupUpgradeHandlers() {
	for _, upgrade := range Upgrades {
		app.UpgradeKeeper.SetUpgradeHandler(
			upgrade.UpgradeName,
			upgrade.CreateUpgradeHandler(
				app.moduleManager,
				app.configurator,
				&app.AppKeepers,
			),
		)
	}
}

// configure store loader that checks if version == upgradeHeight and applies store upgrades
func (app *Application) setupUpgradeStoreLoaders() {
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		return
	}

	for _, upgrade := range Upgrades {
		if upgradeInfo.Name == upgrade.UpgradeName {
			app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &upgrade.StoreUpgrades))
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
func (app *Application) setupPostHandler() {
	postDecorators := []sdk.AnteDecorator{
		// posthandler.NewTipDecorator(app.BankKeeper),
		// ... add more decorators as needed
	}
	postHandler := sdk.ChainAnteDecorators(postDecorators...)
	app.SetPostHandler(postHandler)
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
	tmservice.RegisterTendermintService(clientCtx, app.BaseApp.GRPCQueryRouter(), app.interfaceRegistry, app.Query)
}
func (app *Application) RegisterNodeService(clientCtx client.Context) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter())
}
func (app *Application) LoadHeight(height int64) error {
	return app.BaseApp.LoadVersion(height)
}
