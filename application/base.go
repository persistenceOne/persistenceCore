/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package application

import (
	"fmt"
	icaControllerTypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/controller/types"
	icaTypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
	"github.com/strangelove-ventures/packet-forward-middleware/v2/router"
	"io"
	"log"
	stdlog "log"
	"net/http"
	"os"
	"path/filepath"

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
	"github.com/cosmos/cosmos-sdk/x/authz"
	sdkAuthzKeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	sdkAuthzModule "github.com/cosmos/cosmos-sdk/x/authz/module"
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
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	sdkFeeGrantKeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	sdkFeeGrantModule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilTypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	sdkGovKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	sdkGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
	ica "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts"
	icaHost "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host"
	icaHostKeeper "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/keeper"
	icaHostTypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	"github.com/cosmos/ibc-go/v3/modules/apps/transfer"
	ibcTransferKeeper "github.com/cosmos/ibc-go/v3/modules/apps/transfer/keeper"
	ibcTransferTypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcCore "github.com/cosmos/ibc-go/v3/modules/core"
	ibcClient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcClientTypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibcConnectionTypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
	ibcTypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibcHost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	sdkIBCKeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	"github.com/gogo/protobuf/grpc"
	"github.com/gorilla/mux"
	applicationParams "github.com/persistenceOne/persistenceCore/application/params"
	"github.com/persistenceOne/persistenceCore/x/halving"
	"github.com/rakyll/statik/fs"
	routerKeeper "github.com/strangelove-ventures/packet-forward-middleware/v2/router/keeper"
	routerTypes "github.com/strangelove-ventures/packet-forward-middleware/v2/router/types"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	tendermintJSON "github.com/tendermint/tendermint/libs/json"
	tendermintLog "github.com/tendermint/tendermint/libs/log"
	tendermintOS "github.com/tendermint/tendermint/libs/os"
	tendermintProto "github.com/tendermint/tendermint/proto/tendermint/types"
	tendermintDB "github.com/tendermint/tm-db"
	"honnef.co/go/tools/version"
)

type Application struct {
	*baseapp.BaseApp
	legacyAmino       *codec.LegacyAmino
	applicationCodec  codec.Codec
	interfaceRegistry types.InterfaceRegistry

	keys map[string]*sdkTypes.KVStoreKey

	AccountKeeper      authKeeper.AccountKeeper
	BankKeeper         sdkBankKeeper.Keeper
	CapabilityKeeper   *sdkCapabilityKeeper.Keeper
	StakingKeeper      sdkStakingKeeper.Keeper
	SlashingKeeper     slashingKeeper.Keeper
	MintKeeper         sdkMintKeeper.Keeper
	DistributionKeeper sdkDistributionKeeper.Keeper
	GovKeeper          sdkGovKeeper.Keeper
	UpgradeKeeper      sdkUpgradeKeeper.Keeper
	CrisisKeeper       sdkCrisisKeeper.Keeper
	ParamsKeeper       sdkParamsKeeper.Keeper
	IBCKeeper          *sdkIBCKeeper.Keeper
	ICAHostKeeper      icaHostKeeper.Keeper
	EvidenceKeeper     sdkEvidenceKeeper.Keeper
	TransferKeeper     ibcTransferKeeper.Keeper
	FeegrantKeeper     sdkFeeGrantKeeper.Keeper
	AuthzKeeper        sdkAuthzKeeper.Keeper
	HalvingKeeper      halving.Keeper
	RouterKeeper       routerKeeper.Keeper

	moduleManager *sdkTypesModule.Manager

	configurator      sdkTypesModule.Configurator
	simulationManager *sdkTypesModule.SimulationManager

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      sdkCapabilityKeeper.ScopedKeeper
	ScopedTransferKeeper sdkCapabilityKeeper.ScopedKeeper
	ScopedICAHostKeeper  sdkCapabilityKeeper.ScopedKeeper
}

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

func (app Application) ApplicationCodec() codec.Codec {
	return app.applicationCodec
}

func (application Application) Name() string {
	return application.BaseApp.Name()
}

func (application Application) LegacyAmino() *codec.LegacyAmino {
	return application.legacyAmino
}

func (application Application) BeginBlocker(ctx sdkTypes.Context, req abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	return application.moduleManager.BeginBlock(ctx, req)
}

func (application Application) EndBlocker(ctx sdkTypes.Context, req abciTypes.RequestEndBlock) abciTypes.ResponseEndBlock {
	return application.moduleManager.EndBlock(ctx, req)
}

func (application Application) InitChainer(ctx sdkTypes.Context, req abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
	var genesisState GenesisState
	if err := tendermintJSON.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}

	application.UpgradeKeeper.SetModuleVersionMap(ctx, application.moduleManager.GetVersionMap())

	return application.moduleManager.InitGenesis(ctx, application.applicationCodec, genesisState)
}
func (application Application) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string) (serverTypes.ExportedApp, error) {
	context := application.BaseApp.NewContext(true, tendermintProto.Header{Height: application.BaseApp.LastBlockHeight()})

	height := application.BaseApp.LastBlockHeight() + 1
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

		application.CrisisKeeper.AssertInvariants(context)

		application.StakingKeeper.IterateValidators(context, func(_ int64, val sdkStakingTypes.ValidatorI) (stop bool) {
			_, _ = application.DistributionKeeper.WithdrawValidatorCommission(context, val.GetOperator())
			return false
		})

		delegations := application.StakingKeeper.GetAllDelegations(context)
		for _, delegation := range delegations {
			validatorAddress, Error := sdkTypes.ValAddressFromBech32(delegation.ValidatorAddress)
			if Error != nil {
				panic(Error)
			}

			delegatorAddress, Error := sdkTypes.AccAddressFromBech32(delegation.DelegatorAddress)
			if Error != nil {
				panic(Error)
			}

			_, _ = application.DistributionKeeper.WithdrawDelegationRewards(context, delegatorAddress, validatorAddress)
		}

		application.DistributionKeeper.DeleteAllValidatorSlashEvents(context)

		application.DistributionKeeper.DeleteAllValidatorHistoricalRewards(context)

		height := context.BlockHeight()
		context = context.WithBlockHeight(0)

		application.StakingKeeper.IterateValidators(context, func(_ int64, val sdkStakingTypes.ValidatorI) (stop bool) {

			scraps := application.DistributionKeeper.GetValidatorOutstandingRewardsCoins(context, val.GetOperator())
			feePool := application.DistributionKeeper.GetFeePool(context)
			feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
			application.DistributionKeeper.SetFeePool(context, feePool)

			application.DistributionKeeper.Hooks().AfterValidatorCreated(context, val.GetOperator())
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

			application.DistributionKeeper.Hooks().BeforeDelegationCreated(context, delegatorAddress, validatorAddress)
			application.DistributionKeeper.Hooks().AfterDelegationModified(context, delegatorAddress, validatorAddress)
		}

		context = context.WithBlockHeight(height)

		application.StakingKeeper.IterateRedelegations(context, func(_ int64, redelegation sdkStakingTypes.Redelegation) (stop bool) {
			for i := range redelegation.Entries {
				redelegation.Entries[i].CreationHeight = 0
			}
			application.StakingKeeper.SetRedelegation(context, redelegation)
			return false
		})

		application.StakingKeeper.IterateUnbondingDelegations(context, func(_ int64, unbondingDelegation sdkStakingTypes.UnbondingDelegation) (stop bool) {
			for i := range unbondingDelegation.Entries {
				unbondingDelegation.Entries[i].CreationHeight = 0
			}
			application.StakingKeeper.SetUnbondingDelegation(context, unbondingDelegation)
			return false
		})

		store := context.KVStore(application.keys[sdkStakingTypes.StoreKey])
		kvStoreReversePrefixIterator := sdkTypes.KVStoreReversePrefixIterator(store, sdkStakingTypes.ValidatorsKey)
		counter := int16(0)

		for ; kvStoreReversePrefixIterator.Valid(); kvStoreReversePrefixIterator.Next() {
			addr := sdkTypes.ValAddress(kvStoreReversePrefixIterator.Key()[1:])
			validator, found := application.StakingKeeper.GetValidator(context, addr)

			if !found {
				panic("Validator not found!")
			}

			validator.UnbondingHeight = 0

			if applyWhiteList && !whiteListMap[addr.String()] {
				validator.Jailed = true
			}

			application.StakingKeeper.SetValidator(context, validator)
			counter++
		}

		_ = kvStoreReversePrefixIterator.Close()

		_, Error := application.StakingKeeper.ApplyAndReturnValidatorSetUpdates(context)
		if Error != nil {
			log.Fatal(Error)
		}

		application.SlashingKeeper.IterateValidatorSigningInfos(
			context,
			func(validatorConsAddress sdkTypes.ConsAddress, validatorSigningInfo slashingTypes.ValidatorSigningInfo) (stop bool) {
				validatorSigningInfo.StartHeight = 0
				application.SlashingKeeper.SetValidatorSigningInfo(context, validatorConsAddress, validatorSigningInfo)
				return false
			},
		)
	}

	genesisState := application.moduleManager.ExportGenesis(context, application.applicationCodec)
	applicationState, Error := codec.MarshalJSONIndent(application.legacyAmino, genesisState)

	if Error != nil {
		return serverTypes.ExportedApp{}, Error
	}

	validators, err := staking.WriteValidators(context, application.StakingKeeper)

	return serverTypes.ExportedApp{
		AppState:        applicationState,
		Validators:      validators,
		Height:          height,
		ConsensusParams: application.BaseApp.GetConsensusParams(context),
	}, err
}

func (application Application) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range ModuleAccountPermissions {
		modAccAddrs[authTypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (application Application) SimulationManager() *sdkTypesModule.SimulationManager {
	return application.simulationManager
}

func (application Application) ListSnapshots(snapshots abciTypes.RequestListSnapshots) abciTypes.ResponseListSnapshots {
	return application.BaseApp.ListSnapshots(snapshots)
}

func (application Application) OfferSnapshot(snapshot abciTypes.RequestOfferSnapshot) abciTypes.ResponseOfferSnapshot {
	return application.BaseApp.OfferSnapshot(snapshot)
}

func (application Application) LoadSnapshotChunk(chunk abciTypes.RequestLoadSnapshotChunk) abciTypes.ResponseLoadSnapshotChunk {
	return application.BaseApp.LoadSnapshotChunk(chunk)
}

func (application Application) ApplySnapshotChunk(chunk abciTypes.RequestApplySnapshotChunk) abciTypes.ResponseApplySnapshotChunk {
	return application.BaseApp.ApplySnapshotChunk(chunk)
}

func (application Application) RegisterGRPCServer(server grpc.Server) {
	application.BaseApp.RegisterGRPCServer(server)
}

func (application Application) RegisterAPIRoutes(apiServer *api.Server, apiConfig config.APIConfig) {
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

func (application Application) RegisterTxService(clientContect client.Context) {
	authTx.RegisterTxService(application.BaseApp.GRPCQueryRouter(), clientContect, application.BaseApp.Simulate, application.interfaceRegistry)
}

func (application Application) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(application.BaseApp.GRPCQueryRouter(), clientCtx, application.interfaceRegistry)
}

func (application Application) LoadHeight(height int64) error {
	return application.BaseApp.LoadVersion(height)
}
func (application Application) Initialize(applicationName string, encodingConfiguration applicationParams.EncodingConfiguration, moduleAccountPermissions map[string][]string, logger tendermintLog.Logger, db tendermintDB.DB, traceStore io.Writer, loadLatest bool, invCheckPeriod uint, skipUpgradeHeights map[int64]bool, home string, applicationOptions serverTypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp)) Application {
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

	keys := sdkTypes.NewKVStoreKeys(
		authTypes.StoreKey, sdkBankTypes.StoreKey, sdkStakingTypes.StoreKey,
		sdkMintTypes.StoreKey, sdkDistributionTypes.StoreKey, slashingTypes.StoreKey,
		sdkGovTypes.StoreKey, paramsTypes.StoreKey, ibcHost.StoreKey, sdkUpgradeTypes.StoreKey,
		sdkEvidenceTypes.StoreKey, ibcTransferTypes.StoreKey, sdkCapabilityTypes.StoreKey,
		halving.StoreKey, sdkAuthzKeeper.StoreKey, feegrant.StoreKey,
	)

	transientStoreKeys := sdkTypes.NewTransientStoreKeys(paramsTypes.TStoreKey)
	memoryKeys := sdkTypes.NewMemoryStoreKeys(sdkCapabilityTypes.MemStoreKey)

	application.BaseApp = baseApp
	application.legacyAmino = legacyAmino
	application.applicationCodec = applicationCodec
	application.interfaceRegistry = interfaceRegistry
	application.keys = keys

	application.ParamsKeeper = sdkParamsKeeper.NewKeeper(
		applicationCodec,
		legacyAmino,
		keys[paramsTypes.StoreKey],
		transientStoreKeys[paramsTypes.TStoreKey],
	)
	application.BaseApp.SetParamStore(application.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(sdkParamsKeeper.ConsensusParamsKeyTable()))

	application.CapabilityKeeper = sdkCapabilityKeeper.NewKeeper(applicationCodec, keys[sdkCapabilityTypes.StoreKey], memoryKeys[sdkCapabilityTypes.MemStoreKey])
	scopedIBCKeeper := application.CapabilityKeeper.ScopeToModule(ibcHost.ModuleName)
	scopedICAHostKeeper := application.CapabilityKeeper.ScopeToModule(icaHostTypes.SubModuleName)
	scopedTransferKeeper := application.CapabilityKeeper.ScopeToModule(ibcTransferTypes.ModuleName)
	application.CapabilityKeeper.Seal()

	application.AccountKeeper = authKeeper.NewAccountKeeper(
		applicationCodec,
		keys[authTypes.StoreKey],
		application.ParamsKeeper.Subspace(authTypes.ModuleName),
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

	application.BankKeeper = sdkBankKeeper.NewBaseKeeper(
		applicationCodec,
		keys[sdkBankTypes.StoreKey],
		application.AccountKeeper,
		application.ParamsKeeper.Subspace(sdkBankTypes.ModuleName),
		blacklistedAddresses,
	)

	application.AuthzKeeper = sdkAuthzKeeper.NewKeeper(
		keys[sdkAuthzKeeper.StoreKey],
		applicationCodec,
		application.BaseApp.MsgServiceRouter(),
	)

	application.FeegrantKeeper = sdkFeeGrantKeeper.NewKeeper(
		applicationCodec,
		keys[feegrant.StoreKey],
		application.AccountKeeper,
	)

	stakingKeeper := sdkStakingKeeper.NewKeeper(
		applicationCodec,
		keys[sdkStakingTypes.StoreKey],
		application.AccountKeeper,
		application.BankKeeper,
		application.ParamsKeeper.Subspace(sdkStakingTypes.ModuleName),
	)

	application.MintKeeper = sdkMintKeeper.NewKeeper(
		applicationCodec,
		keys[sdkMintTypes.StoreKey],
		application.ParamsKeeper.Subspace(sdkMintTypes.ModuleName),
		&stakingKeeper,
		application.AccountKeeper,
		application.BankKeeper,
		authTypes.FeeCollectorName,
	)

	application.DistributionKeeper = sdkDistributionKeeper.NewKeeper(
		applicationCodec,
		keys[sdkDistributionTypes.StoreKey],
		application.ParamsKeeper.Subspace(sdkDistributionTypes.ModuleName),
		application.AccountKeeper,
		application.BankKeeper,
		&stakingKeeper,
		authTypes.FeeCollectorName,
		blackListedModuleAddresses,
	)
	application.SlashingKeeper = slashingKeeper.NewKeeper(
		applicationCodec,
		keys[slashingTypes.StoreKey],
		&stakingKeeper,
		application.ParamsKeeper.Subspace(slashingTypes.ModuleName),
	)
	application.CrisisKeeper = sdkCrisisKeeper.NewKeeper(
		application.ParamsKeeper.Subspace(sdkCrisisTypes.ModuleName),
		invCheckPeriod,
		application.BankKeeper,
		authTypes.FeeCollectorName,
	)
	application.UpgradeKeeper = sdkUpgradeKeeper.NewKeeper(
		skipUpgradeHeights,
		keys[sdkUpgradeTypes.StoreKey],
		applicationCodec,
		home,
		application.BaseApp,
	)
	application.HalvingKeeper = halving.NewKeeper(
		keys[halving.StoreKey],
		application.ParamsKeeper.Subspace(halving.DefaultParamspace),
		application.MintKeeper,
	)

	application.StakingKeeper = *stakingKeeper.SetHooks(
		sdkStakingTypes.NewMultiStakingHooks(application.DistributionKeeper.Hooks(), application.SlashingKeeper.Hooks()),
	)

	application.IBCKeeper = sdkIBCKeeper.NewKeeper(
		applicationCodec,
		keys[ibcHost.StoreKey],
		application.ParamsKeeper.Subspace(ibcHost.ModuleName),
		application.StakingKeeper,
		application.UpgradeKeeper,
		scopedIBCKeeper,
	)

	govRouter := sdkGovTypes.NewRouter()
	govRouter.AddRoute(
		sdkGovTypes.RouterKey,
		sdkGovTypes.ProposalHandler,
	).AddRoute(
		paramsProposal.RouterKey,
		params.NewParamChangeProposalHandler(application.ParamsKeeper),
	).AddRoute(
		sdkDistributionTypes.RouterKey,
		distribution.NewCommunityPoolSpendProposalHandler(application.DistributionKeeper),
	).AddRoute(
		sdkUpgradeTypes.RouterKey,
		upgrade.NewSoftwareUpgradeProposalHandler(application.UpgradeKeeper),
	).AddRoute(ibcClientTypes.RouterKey, ibcClient.NewClientProposalHandler(application.IBCKeeper.ClientKeeper))

	application.GovKeeper = sdkGovKeeper.NewKeeper(
		applicationCodec,
		keys[sdkGovTypes.StoreKey],
		application.ParamsKeeper.Subspace(sdkGovTypes.ModuleName).WithKeyTable(sdkGovTypes.ParamKeyTable()),
		application.AccountKeeper,
		application.BankKeeper,
		&stakingKeeper,
		govRouter,
	)

	application.TransferKeeper = ibcTransferKeeper.NewKeeper(
		applicationCodec,
		keys[ibcTransferTypes.StoreKey],
		application.ParamsKeeper.Subspace(ibcTransferTypes.ModuleName),
		application.IBCKeeper.ChannelKeeper,
		application.IBCKeeper.ChannelKeeper,
		&application.IBCKeeper.PortKeeper,
		application.AccountKeeper,
		application.BankKeeper,
		scopedTransferKeeper,
	)
	transferModule := transfer.NewAppModule(application.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(application.TransferKeeper)

	application.ICAHostKeeper = icaHostKeeper.NewKeeper(
		applicationCodec,
		keys[icaHostTypes.StoreKey],
		application.ParamsKeeper.Subspace(icaHostTypes.SubModuleName),
		application.IBCKeeper.ChannelKeeper,
		&application.IBCKeeper.PortKeeper,
		application.AccountKeeper,
		scopedICAHostKeeper,
		application.MsgServiceRouter(),
	)

	icaModule := ica.NewAppModule(nil, &application.ICAHostKeeper)
	icaHostIBCModule := icaHost.NewIBCModule(application.ICAHostKeeper)

	application.RouterKeeper = routerKeeper.NewKeeper(
		applicationCodec,
		keys[routerTypes.StoreKey],
		application.ParamsKeeper.Subspace(routerTypes.ModuleName),
		application.TransferKeeper,
		application.DistributionKeeper)

	routerModule := router.NewAppModule(application.RouterKeeper, transferIBCModule)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcTypes.NewRouter()
	ibcRouter.AddRoute(icaHostTypes.SubModuleName, icaHostIBCModule).
		AddRoute(ibcTransferTypes.ModuleName, transferIBCModule)
	application.IBCKeeper.SetRouter(ibcRouter)

	evidenceKeeper := sdkEvidenceKeeper.NewKeeper(
		applicationCodec,
		keys[sdkEvidenceTypes.StoreKey],
		&application.StakingKeeper,
		application.SlashingKeeper,
	)
	application.EvidenceKeeper = *evidenceKeeper

	/****  Module Options ****/
	var skipGenesisInvariants = false

	opt := applicationOptions.Get(crisis.FlagSkipGenesisInvariants)
	if opt, ok := opt.(bool); ok {
		skipGenesisInvariants = opt
	}

	application.moduleManager = sdkTypesModule.NewManager(
		genutil.NewAppModule(
			application.AccountKeeper, application.StakingKeeper, application.BaseApp.DeliverTx,
			encodingConfiguration.TransactionConfig,
		),
		auth.NewAppModule(applicationCodec, application.AccountKeeper, nil),
		vesting.NewAppModule(application.AccountKeeper, application.BankKeeper),
		bank.NewAppModule(applicationCodec, application.BankKeeper, application.AccountKeeper),
		capability.NewAppModule(applicationCodec, *application.CapabilityKeeper),
		crisis.NewAppModule(&application.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(applicationCodec, application.GovKeeper, application.AccountKeeper, application.BankKeeper),
		mint.NewAppModule(applicationCodec, application.MintKeeper, application.AccountKeeper),
		slashing.NewAppModule(applicationCodec, application.SlashingKeeper, application.AccountKeeper, application.BankKeeper, application.StakingKeeper),
		distribution.NewAppModule(applicationCodec, application.DistributionKeeper, application.AccountKeeper, application.BankKeeper, application.StakingKeeper),
		staking.NewAppModule(applicationCodec, application.StakingKeeper, application.AccountKeeper, application.BankKeeper),
		upgrade.NewAppModule(application.UpgradeKeeper),
		evidence.NewAppModule(application.EvidenceKeeper),
		sdkFeeGrantModule.NewAppModule(applicationCodec, application.AccountKeeper, application.BankKeeper, application.FeegrantKeeper, application.interfaceRegistry),
		sdkAuthzModule.NewAppModule(applicationCodec, application.AuthzKeeper, application.AccountKeeper, application.BankKeeper, application.interfaceRegistry),
		ibcCore.NewAppModule(application.IBCKeeper),
		params.NewAppModule(application.ParamsKeeper),
		halving.NewAppModule(applicationCodec, application.HalvingKeeper),
		transferModule,
		icaModule,
		routerModule,
	)

	application.moduleManager.SetOrderBeginBlockers(
		sdkUpgradeTypes.ModuleName,
		sdkCapabilityTypes.ModuleName,
		sdkMintTypes.ModuleName,
		sdkDistributionTypes.ModuleName,
		slashingTypes.ModuleName,
		sdkEvidenceTypes.ModuleName,
		sdkStakingTypes.ModuleName,
		ibcHost.ModuleName,
	)
	application.moduleManager.SetOrderEndBlockers(
		sdkCrisisTypes.ModuleName,
		sdkGovTypes.ModuleName,
		sdkStakingTypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		halving.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	application.moduleManager.SetOrderInitGenesis(
		sdkCapabilityTypes.ModuleName,
		authTypes.ModuleName,
		sdkBankTypes.ModuleName,
		sdkDistributionTypes.ModuleName,
		sdkStakingTypes.ModuleName,
		slashingTypes.ModuleName,
		sdkGovTypes.ModuleName,
		sdkMintTypes.ModuleName,
		sdkCrisisTypes.ModuleName,
		ibcHost.ModuleName,
		genutilTypes.ModuleName,
		sdkEvidenceTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		halving.ModuleName,
	)

	application.moduleManager.RegisterInvariants(&application.CrisisKeeper)
	application.moduleManager.RegisterRoutes(application.BaseApp.Router(), application.BaseApp.QueryRouter(), encodingConfiguration.Amino)
	application.configurator = sdkTypesModule.NewConfigurator(application.applicationCodec, application.BaseApp.MsgServiceRouter(), application.BaseApp.GRPCQueryRouter())
	application.moduleManager.RegisterServices(application.configurator)

	simulationManager := sdkTypesModule.NewSimulationManager(
		auth.NewAppModule(applicationCodec, application.AccountKeeper, authSimulation.RandomGenesisAccounts),
		bank.NewAppModule(applicationCodec, application.BankKeeper, application.AccountKeeper),
		capability.NewAppModule(applicationCodec, *application.CapabilityKeeper),
		gov.NewAppModule(applicationCodec, application.GovKeeper, application.AccountKeeper, application.BankKeeper),
		mint.NewAppModule(applicationCodec, application.MintKeeper, application.AccountKeeper),
		staking.NewAppModule(applicationCodec, application.StakingKeeper, application.AccountKeeper, application.BankKeeper),
		distribution.NewAppModule(applicationCodec, application.DistributionKeeper, application.AccountKeeper, application.BankKeeper, application.StakingKeeper),
		slashing.NewAppModule(applicationCodec, application.SlashingKeeper, application.AccountKeeper, application.BankKeeper, application.StakingKeeper),
		params.NewAppModule(application.ParamsKeeper),
		halving.NewAppModule(applicationCodec, application.HalvingKeeper),
		sdkAuthzModule.NewAppModule(applicationCodec, application.AuthzKeeper, application.AccountKeeper, application.BankKeeper, application.interfaceRegistry),
		sdkFeeGrantModule.NewAppModule(applicationCodec, application.AccountKeeper, application.BankKeeper, application.FeegrantKeeper, application.interfaceRegistry),
		ibcCore.NewAppModule(application.IBCKeeper),
		transferModule,
	)

	simulationManager.RegisterStoreDecoders()
	application.simulationManager = simulationManager

	application.BaseApp.MountKVStores(keys)
	application.BaseApp.MountTransientStores(transientStoreKeys)
	application.BaseApp.MountMemoryStores(memoryKeys)

	anteHandler, err := NewAnteHandler(
		HandlerOptions{
			HandlerOptions: ante.HandlerOptions{
				AccountKeeper:   application.AccountKeeper,
				BankKeeper:      application.BankKeeper,
				FeegrantKeeper:  application.FeegrantKeeper,
				SignModeHandler: encodingConfiguration.TransactionConfig.SignModeHandler(),
				SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
			},
			IBCKeeper: application.IBCKeeper,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	application.BaseApp.SetAnteHandler(anteHandler)
	application.BaseApp.SetInitChainer(application.InitChainer)
	application.BaseApp.SetBeginBlocker(application.moduleManager.BeginBlock)
	application.BaseApp.SetEndBlocker(application.moduleManager.EndBlock)

	application.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx sdkTypes.Context, _ sdkUpgradeTypes.Plan, fromVM sdkTypesModule.VersionMap) (sdkTypesModule.VersionMap, error) {
			application.IBCKeeper.ConnectionKeeper.SetParams(ctx, ibcConnectionTypes.DefaultParams())
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
					liquidityMsgCreatePool,
					liquidityMsgSwapWithinBatch,
					liquidityMsgDepositWithinBatch,
					liquidityMsgWithdrawWithinBatch,
				},
			}
			ctx.Logger().Info("start to init interchainaccount module...")
			// initialize ICS27 module
			icaModule.InitModule(ctx, controllerParams, hostParams)

			ctx.Logger().Info("start to run module migrations...")

			// RunMigrations twice is just a way to make auth module's migrates after staking
			return application.moduleManager.RunMigrations(ctx, application.configurator, fromVM)

		},
	)

	upgradeInfo, err := application.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == UpgradeName && !application.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storeTypes.StoreUpgrades{
			Added: []string{authz.ModuleName, feegrant.ModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		application.BaseApp.SetStoreLoader(sdkUpgradeTypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}

	if loadLatest {
		if err := application.BaseApp.LoadLatestVersion(); err != nil {
			tendermintOS.Exit(err.Error())
		}
	}
	application.ScopedIBCKeeper = scopedIBCKeeper
	application.ScopedTransferKeeper = scopedTransferKeeper

	return application
}

func NewApplication() *Application {
	return &Application{}
}
