/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package application

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

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
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkTypesModule "github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authRest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	authKeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authSimulation "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authTx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	sdkBankKeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	sdkBankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	sdkCapabilityKeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	sdkCapabilityTypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	sdkCrisisKeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	sdkCrisisTypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	sdkDistributionKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	sdkDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	sdkEvidenceKeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	sdkEvidenceTypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilTypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	sdkGovKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	sdkGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer"
	ibcTransferKeeper "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/keeper"
	ibcTransferTypes "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/types"
	ibcCore "github.com/cosmos/cosmos-sdk/x/ibc/core"
	ibcClient "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client"
	ibcTypes "github.com/cosmos/cosmos-sdk/x/ibc/core/05-port/types"
	ibcHost "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	sdkIBCKeeper "github.com/cosmos/cosmos-sdk/x/ibc/core/keeper"
	"github.com/cosmos/cosmos-sdk/x/mint"
	sdkMintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	sdkMintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	sdkParamsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsProposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingKeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	sdkStakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	sdkStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	sdkUpgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	sdkUpgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/gogo/protobuf/grpc"
	"github.com/gorilla/mux"
	applicationParams "github.com/persistenceOne/persistenceCore/application/params"
	"github.com/persistenceOne/persistenceCore/x/halving"
	"github.com/rakyll/statik/fs"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	tendermintJSON "github.com/tendermint/tendermint/libs/json"
	tendermintLog "github.com/tendermint/tendermint/libs/log"
	tendermintOS "github.com/tendermint/tendermint/libs/os"
	tendermintProto "github.com/tendermint/tendermint/proto/tendermint/types"
	tendermintDB "github.com/tendermint/tm-db"
	"honnef.co/go/tools/version"
)

type application struct {
	baseApp           *baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	applicationCodec  codec.Marshaler
	interfaceRegistry types.InterfaceRegistry

	keys map[string]*sdkTypes.KVStoreKey

	stakingKeeper      sdkStakingKeeper.Keeper
	slashingKeeper     slashingKeeper.Keeper
	distributionKeeper sdkDistributionKeeper.Keeper
	crisisKeeper       sdkCrisisKeeper.Keeper

	moduleManager     *sdkTypesModule.Manager
	simulationManager *sdkTypesModule.SimulationManager
}

var (
	_ simapp.App              = (*application)(nil)
	_ serverTypes.Application = (*application)(nil)
)

func (application application) BaseApp() *baseapp.BaseApp {
	return application.baseApp
}

func (application application) ApplicationCodec() codec.Marshaler {
	return application.applicationCodec
}

func (application application) Name() string {
	return application.baseApp.Name()
}

func (application application) LegacyAmino() *codec.LegacyAmino {
	return application.legacyAmino
}

func (application application) BeginBlocker(ctx sdkTypes.Context, req abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	return application.moduleManager.BeginBlock(ctx, req)
}

func (application application) EndBlocker(ctx sdkTypes.Context, req abciTypes.RequestEndBlock) abciTypes.ResponseEndBlock {
	return application.moduleManager.EndBlock(ctx, req)
}

func (application application) InitChainer(ctx sdkTypes.Context, req abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
	var genesisState GenesisState
	if err := tendermintJSON.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	return application.moduleManager.InitGenesis(ctx, application.applicationCodec, genesisState)
}
func (application application) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string) (serverTypes.ExportedApp, error) {
	context := application.baseApp.NewContext(true, tendermintProto.Header{Height: application.baseApp.LastBlockHeight()})

	height := application.baseApp.LastBlockHeight() + 1
	if forZeroHeight {
		height = 0
		applyWhiteList := false

		if len(jailWhiteList) > 0 {
			applyWhiteList = true
		}

		whiteListMap := make(map[string]bool)

		for _, address := range jailWhiteList {
			if _, Error := sdkTypes.ValAddressFromBech32(address); Error != nil {
				panic(Error)
			}

			whiteListMap[address] = true
		}

		application.crisisKeeper.AssertInvariants(context)

		application.stakingKeeper.IterateValidators(context, func(_ int64, val sdkStakingTypes.ValidatorI) (stop bool) {
			_, _ = application.distributionKeeper.WithdrawValidatorCommission(context, val.GetOperator())
			return false
		})

		delegations := application.stakingKeeper.GetAllDelegations(context)
		for _, delegation := range delegations {
			validatorAddress, Error := sdkTypes.ValAddressFromBech32(delegation.ValidatorAddress)
			if Error != nil {
				panic(Error)
			}

			delegatorAddress, Error := sdkTypes.AccAddressFromBech32(delegation.DelegatorAddress)
			if Error != nil {
				panic(Error)
			}

			_, _ = application.distributionKeeper.WithdrawDelegationRewards(context, delegatorAddress, validatorAddress)
		}

		application.distributionKeeper.DeleteAllValidatorSlashEvents(context)

		application.distributionKeeper.DeleteAllValidatorHistoricalRewards(context)

		height := context.BlockHeight()
		context = context.WithBlockHeight(0)

		application.stakingKeeper.IterateValidators(context, func(_ int64, val sdkStakingTypes.ValidatorI) (stop bool) {

			scraps := application.distributionKeeper.GetValidatorOutstandingRewardsCoins(context, val.GetOperator())
			feePool := application.distributionKeeper.GetFeePool(context)
			feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
			application.distributionKeeper.SetFeePool(context, feePool)

			application.distributionKeeper.Hooks().AfterValidatorCreated(context, val.GetOperator())
			return false
		})

		for _, delegation := range delegations {
			validatorAddress, Error := sdkTypes.ValAddressFromBech32(delegation.ValidatorAddress)
			if Error != nil {
				panic(Error)
			}

			delegatorAddress, Error := sdkTypes.AccAddressFromBech32(delegation.DelegatorAddress)
			if Error != nil {
				panic(Error)
			}

			application.distributionKeeper.Hooks().BeforeDelegationCreated(context, delegatorAddress, validatorAddress)
			application.distributionKeeper.Hooks().AfterDelegationModified(context, delegatorAddress, validatorAddress)
		}

		context = context.WithBlockHeight(height)

		application.stakingKeeper.IterateRedelegations(context, func(_ int64, redelegation sdkStakingTypes.Redelegation) (stop bool) {
			for i := range redelegation.Entries {
				redelegation.Entries[i].CreationHeight = 0
			}
			application.stakingKeeper.SetRedelegation(context, redelegation)
			return false
		})

		application.stakingKeeper.IterateUnbondingDelegations(context, func(_ int64, unbondingDelegation sdkStakingTypes.UnbondingDelegation) (stop bool) {
			for i := range unbondingDelegation.Entries {
				unbondingDelegation.Entries[i].CreationHeight = 0
			}
			application.stakingKeeper.SetUnbondingDelegation(context, unbondingDelegation)
			return false
		})

		store := context.KVStore(application.keys[sdkStakingTypes.StoreKey])
		kvStoreReversePrefixIterator := sdkTypes.KVStoreReversePrefixIterator(store, sdkStakingTypes.ValidatorsKey)
		counter := int16(0)

		for ; kvStoreReversePrefixIterator.Valid(); kvStoreReversePrefixIterator.Next() {
			addr := sdkTypes.ValAddress(kvStoreReversePrefixIterator.Key()[1:])
			validator, found := application.stakingKeeper.GetValidator(context, addr)

			if !found {
				panic("Validator not found!")
			}

			validator.UnbondingHeight = 0

			if applyWhiteList && !whiteListMap[addr.String()] {
				validator.Jailed = true
			}

			application.stakingKeeper.SetValidator(context, validator)
			counter++
		}

		_ = kvStoreReversePrefixIterator.Close()

		_, Error := application.stakingKeeper.ApplyAndReturnValidatorSetUpdates(context)
		if Error != nil {
			log.Fatal(Error)
		}

		application.slashingKeeper.IterateValidatorSigningInfos(
			context,
			func(validatorConsAddress sdkTypes.ConsAddress, validatorSigningInfo slashingTypes.ValidatorSigningInfo) (stop bool) {
				validatorSigningInfo.StartHeight = 0
				application.slashingKeeper.SetValidatorSigningInfo(context, validatorConsAddress, validatorSigningInfo)
				return false
			},
		)
	}

	genesisState := application.moduleManager.ExportGenesis(context, application.applicationCodec)
	applicationState, Error := codec.MarshalJSONIndent(application.legacyAmino, genesisState)

	if Error != nil {
		return serverTypes.ExportedApp{}, Error
	}

	validators, err := staking.WriteValidators(context, application.stakingKeeper)

	return serverTypes.ExportedApp{
		AppState:        applicationState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: application.baseApp.GetConsensusParams(context),
	}, err
}

func (application application) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range ModuleAccountPermissions {
		modAccAddrs[authTypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (application application) SimulationManager() *sdkTypesModule.SimulationManager {
	return application.simulationManager
}

func (application application) ListSnapshots(snapshots abciTypes.RequestListSnapshots) abciTypes.ResponseListSnapshots {
	return application.baseApp.ListSnapshots(snapshots)
}

func (application application) OfferSnapshot(snapshot abciTypes.RequestOfferSnapshot) abciTypes.ResponseOfferSnapshot {
	return application.baseApp.OfferSnapshot(snapshot)
}

func (application application) LoadSnapshotChunk(chunk abciTypes.RequestLoadSnapshotChunk) abciTypes.ResponseLoadSnapshotChunk {
	return application.baseApp.LoadSnapshotChunk(chunk)
}

func (application application) ApplySnapshotChunk(chunk abciTypes.RequestApplySnapshotChunk) abciTypes.ResponseApplySnapshotChunk {
	return application.baseApp.ApplySnapshotChunk(chunk)
}

func (application application) RegisterGRPCServer(context client.Context, server grpc.Server) {
	application.baseApp.RegisterGRPCServer(context, server)
}

func (application application) RegisterAPIRoutes(apiServer *api.Server, apiConfig config.APIConfig) {
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

func (application application) RegisterTxService(clientContect client.Context) {
	authTx.RegisterTxService(application.baseApp.GRPCQueryRouter(), clientContect, application.baseApp.Simulate, application.interfaceRegistry)
}

func (application application) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(application.baseApp.GRPCQueryRouter(), clientCtx, application.interfaceRegistry)
}

func (application application) Info(requestInfo abciTypes.RequestInfo) abciTypes.ResponseInfo {
	return application.baseApp.Info(requestInfo)
}

func (application application) SetOption(requestSetOption abciTypes.RequestSetOption) abciTypes.ResponseSetOption {
	return application.baseApp.SetOption(requestSetOption)
}

func (application application) Query(requestQuery abciTypes.RequestQuery) abciTypes.ResponseQuery {
	return application.baseApp.Query(requestQuery)
}

func (application application) CheckTx(requestCheckTx abciTypes.RequestCheckTx) abciTypes.ResponseCheckTx {
	return application.baseApp.CheckTx(requestCheckTx)
}

func (application application) InitChain(requestInitChain abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
	return application.baseApp.InitChain(requestInitChain)
}

func (application application) BeginBlock(requestBeginBlock abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	return application.baseApp.BeginBlock(requestBeginBlock)
}

func (application application) DeliverTx(requestDeliverTx abciTypes.RequestDeliverTx) abciTypes.ResponseDeliverTx {
	return application.baseApp.DeliverTx(requestDeliverTx)
}

func (application application) EndBlock(requestEndBlock abciTypes.RequestEndBlock) abciTypes.ResponseEndBlock {
	return application.baseApp.EndBlock(requestEndBlock)
}

func (application application) Commit() abciTypes.ResponseCommit {
	return application.baseApp.Commit()
}

func (application application) LoadHeight(height int64) error {
	return application.baseApp.LoadVersion(height)
}
func (application application) Initialize(applicationName string, encodingConfiguration applicationParams.EncodingConfiguration, moduleAccountPermissions map[string][]string, logger tendermintLog.Logger, db tendermintDB.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string, applicationOptions serverTypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp)) application {
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
	baseApp.SetAppVersion(version.Version)
	baseApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdkTypes.NewKVStoreKeys(
		authTypes.StoreKey, sdkBankTypes.StoreKey, sdkStakingTypes.StoreKey,
		sdkMintTypes.StoreKey, sdkDistributionTypes.StoreKey, slashingTypes.StoreKey,
		sdkGovTypes.StoreKey, paramsTypes.StoreKey, ibcHost.StoreKey, sdkUpgradeTypes.StoreKey,
		sdkEvidenceTypes.StoreKey, ibcTransferTypes.StoreKey, sdkCapabilityTypes.StoreKey,
		halving.StoreKey,
	)

	transientStoreKeys := sdkTypes.NewTransientStoreKeys(paramsTypes.TStoreKey)
	memoryKeys := sdkTypes.NewMemoryStoreKeys(sdkCapabilityTypes.MemStoreKey)

	application.baseApp = baseApp
	application.legacyAmino = legacyAmino
	application.applicationCodec = applicationCodec
	application.interfaceRegistry = interfaceRegistry
	application.keys = keys

	paramsKeeper := sdkParamsKeeper.NewKeeper(
		applicationCodec,
		legacyAmino,
		keys[paramsTypes.StoreKey],
		transientStoreKeys[paramsTypes.TStoreKey],
	)
	application.baseApp.SetParamStore(paramsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(sdkParamsKeeper.ConsensusParamsKeyTable()))

	capabilityKeeper := sdkCapabilityKeeper.NewKeeper(applicationCodec, keys[sdkCapabilityTypes.StoreKey], memoryKeys[sdkCapabilityTypes.MemStoreKey])
	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibcHost.ModuleName)
	scopedTransferKeeper := capabilityKeeper.ScopeToModule(ibcTransferTypes.ModuleName)

	accountKeeper := authKeeper.NewAccountKeeper(
		applicationCodec,
		keys[authTypes.StoreKey],
		paramsKeeper.Subspace(authTypes.ModuleName),
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

	bankKeeper := sdkBankKeeper.NewBaseKeeper(
		applicationCodec,
		keys[sdkBankTypes.StoreKey],
		accountKeeper,
		paramsKeeper.Subspace(sdkBankTypes.ModuleName),
		blacklistedAddresses,
	)

	stakingKeeper := sdkStakingKeeper.NewKeeper(
		applicationCodec,
		keys[sdkStakingTypes.StoreKey],
		accountKeeper,
		bankKeeper,
		paramsKeeper.Subspace(sdkStakingTypes.ModuleName),
	)

	mintKeeper := sdkMintKeeper.NewKeeper(
		applicationCodec,
		keys[sdkMintTypes.StoreKey],
		paramsKeeper.Subspace(sdkMintTypes.ModuleName),
		&stakingKeeper,
		accountKeeper,
		bankKeeper,
		authTypes.FeeCollectorName,
	)

	application.distributionKeeper = sdkDistributionKeeper.NewKeeper(
		applicationCodec,
		keys[sdkDistributionTypes.StoreKey],
		paramsKeeper.Subspace(sdkDistributionTypes.ModuleName),
		accountKeeper,
		bankKeeper,
		&stakingKeeper,
		authTypes.FeeCollectorName,
		blackListedModuleAddresses,
	)
	application.slashingKeeper = slashingKeeper.NewKeeper(
		applicationCodec,
		keys[slashingTypes.StoreKey],
		&stakingKeeper,
		paramsKeeper.Subspace(slashingTypes.ModuleName),
	)
	application.crisisKeeper = sdkCrisisKeeper.NewKeeper(
		paramsKeeper.Subspace(sdkCrisisTypes.ModuleName),
		invCheckPeriod,
		bankKeeper,
		authTypes.FeeCollectorName,
	)
	upgradeKeeper := sdkUpgradeKeeper.NewKeeper(
		skipUpgradeHeights,
		keys[sdkUpgradeTypes.StoreKey],
		applicationCodec,
		home,
	)
	halvingKeeper := halving.NewKeeper(keys[halving.StoreKey], paramsKeeper.Subspace(halving.DefaultParamspace), mintKeeper)

	application.stakingKeeper = *stakingKeeper.SetHooks(
		sdkStakingTypes.NewMultiStakingHooks(application.distributionKeeper.Hooks(), application.slashingKeeper.Hooks()),
	)

	ibcKeeper := sdkIBCKeeper.NewKeeper(
		applicationCodec, keys[ibcHost.StoreKey], paramsKeeper.Subspace(ibcHost.ModuleName), application.stakingKeeper, scopedIBCKeeper,
	)

	govRouter := sdkGovTypes.NewRouter()
	govRouter.AddRoute(
		sdkGovTypes.RouterKey,
		sdkGovTypes.ProposalHandler,
	).AddRoute(
		paramsProposal.RouterKey,
		params.NewParamChangeProposalHandler(paramsKeeper),
	).AddRoute(
		sdkDistributionTypes.RouterKey,
		distribution.NewCommunityPoolSpendProposalHandler(application.distributionKeeper),
	).AddRoute(
		sdkUpgradeTypes.RouterKey,
		upgrade.NewSoftwareUpgradeProposalHandler(upgradeKeeper),
	).AddRoute(ibcHost.RouterKey, ibcClient.NewClientUpdateProposalHandler(ibcKeeper.ClientKeeper))

	transferKeeper := ibcTransferKeeper.NewKeeper(
		applicationCodec, keys[ibcTransferTypes.StoreKey], paramsKeeper.Subspace(ibcTransferTypes.ModuleName),
		ibcKeeper.ChannelKeeper, &ibcKeeper.PortKeeper,
		accountKeeper, bankKeeper, scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(transferKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcTypes.NewRouter()
	ibcRouter.AddRoute(ibcTransferTypes.ModuleName, transferModule)

	evidenceKeeper := sdkEvidenceKeeper.NewKeeper(
		applicationCodec,
		keys[sdkEvidenceTypes.StoreKey],
		&stakingKeeper,
		application.slashingKeeper,
	)

	ibcKeeper.SetRouter(ibcRouter)

	govKeeper := sdkGovKeeper.NewKeeper(
		applicationCodec,
		keys[sdkGovTypes.StoreKey],
		paramsKeeper.Subspace(sdkGovTypes.ModuleName).WithKeyTable(sdkGovTypes.ParamKeyTable()),
		accountKeeper,
		bankKeeper,
		&stakingKeeper,
		govRouter,
	)
	/****  Module Options ****/
	var skipGenesisInvariants = false

	opt := applicationOptions.Get(crisis.FlagSkipGenesisInvariants)
	if opt, ok := opt.(bool); ok {
		skipGenesisInvariants = opt
	}

	application.moduleManager = sdkTypesModule.NewManager(
		genutil.NewAppModule(
			accountKeeper, stakingKeeper, application.baseApp.DeliverTx,
			encodingConfiguration.TransactionConfig,
		),
		auth.NewAppModule(applicationCodec, accountKeeper, nil),
		vesting.NewAppModule(accountKeeper, bankKeeper),
		bank.NewAppModule(applicationCodec, bankKeeper, accountKeeper),
		capability.NewAppModule(applicationCodec, *capabilityKeeper),
		crisis.NewAppModule(&application.crisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(applicationCodec, govKeeper, accountKeeper, bankKeeper),
		mint.NewAppModule(applicationCodec, mintKeeper, accountKeeper),
		slashing.NewAppModule(applicationCodec, application.slashingKeeper, accountKeeper, bankKeeper, application.stakingKeeper),
		distribution.NewAppModule(applicationCodec, application.distributionKeeper, accountKeeper, bankKeeper, application.stakingKeeper),
		staking.NewAppModule(applicationCodec, application.stakingKeeper, accountKeeper, bankKeeper),
		upgrade.NewAppModule(upgradeKeeper),
		evidence.NewAppModule(*evidenceKeeper),
		ibcCore.NewAppModule(ibcKeeper),
		params.NewAppModule(paramsKeeper),
		transferModule,
		halving.NewAppModule(applicationCodec, halvingKeeper),
	)

	application.moduleManager.SetOrderBeginBlockers(
		sdkUpgradeTypes.ModuleName, sdkMintTypes.ModuleName, sdkDistributionTypes.ModuleName, slashingTypes.ModuleName,
		sdkEvidenceTypes.ModuleName, sdkStakingTypes.ModuleName, ibcHost.ModuleName,
	)
	application.moduleManager.SetOrderEndBlockers(sdkCrisisTypes.ModuleName, sdkGovTypes.ModuleName, sdkStakingTypes.ModuleName, halving.ModuleName)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	application.moduleManager.SetOrderInitGenesis(
		sdkCapabilityTypes.ModuleName, authTypes.ModuleName, sdkBankTypes.ModuleName, sdkDistributionTypes.ModuleName, sdkStakingTypes.ModuleName,
		slashingTypes.ModuleName, sdkGovTypes.ModuleName, sdkMintTypes.ModuleName, sdkCrisisTypes.ModuleName,
		ibcHost.ModuleName, genutilTypes.ModuleName, sdkEvidenceTypes.ModuleName, ibcTransferTypes.ModuleName, halving.ModuleName,
	)

	application.moduleManager.RegisterInvariants(&application.crisisKeeper)
	application.moduleManager.RegisterRoutes(application.baseApp.Router(), application.baseApp.QueryRouter(), encodingConfiguration.Amino)
	application.moduleManager.RegisterServices(sdkTypesModule.NewConfigurator(application.baseApp.MsgServiceRouter(), application.baseApp.GRPCQueryRouter()))

	simulationManager := sdkTypesModule.NewSimulationManager(
		auth.NewAppModule(applicationCodec, accountKeeper, authSimulation.RandomGenesisAccounts),
		bank.NewAppModule(applicationCodec, bankKeeper, accountKeeper),
		gov.NewAppModule(applicationCodec, govKeeper, accountKeeper, bankKeeper),
		mint.NewAppModule(applicationCodec, mintKeeper, accountKeeper),
		staking.NewAppModule(applicationCodec, application.stakingKeeper, accountKeeper, bankKeeper),
		distribution.NewAppModule(applicationCodec, application.distributionKeeper, accountKeeper, bankKeeper, application.stakingKeeper),
		slashing.NewAppModule(applicationCodec, application.slashingKeeper, accountKeeper, bankKeeper, application.stakingKeeper),
		params.NewAppModule(paramsKeeper),
		halving.NewAppModule(applicationCodec, halvingKeeper),
	)

	simulationManager.RegisterStoreDecoders()
	application.simulationManager = simulationManager

	application.baseApp.MountKVStores(keys)
	application.baseApp.MountTransientStores(transientStoreKeys)
	application.baseApp.MountMemoryStores(memoryKeys)

	application.baseApp.SetBeginBlocker(application.moduleManager.BeginBlock)
	application.baseApp.SetEndBlocker(application.moduleManager.EndBlock)
	application.baseApp.SetInitChainer(func(context sdkTypes.Context, requestInitChain abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
		var genesisState map[string]json.RawMessage
		legacyAmino.MustUnmarshalJSON(requestInitChain.AppStateBytes, &genesisState)
		return application.moduleManager.InitGenesis(context, applicationCodec, genesisState)
	})
	application.baseApp.SetAnteHandler(ante.NewAnteHandler(accountKeeper, bankKeeper, ante.DefaultSigVerificationGasConsumer, encodingConfiguration.TransactionConfig.SignModeHandler()))

	if loadLatest {
		if err := application.baseApp.LoadLatestVersion(); err != nil {
			tendermintOS.Exit(err.Error())
		}

		ctx := application.baseApp.NewUncachedContext(true, tendermintProto.Header{})
		capabilityKeeper.InitializeAndSeal(ctx)
	}

	return application
}

func NewApplication() *application {
	return &application{}
}
