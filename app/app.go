/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package app

import (
	"fmt"
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
	"github.com/cosmos/cosmos-sdk/types/module"
	sdkTypesModule "github.com/cosmos/cosmos-sdk/types/module"
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
	distributionClient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	sdkDistributionKeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	sdkDistributionTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	sdkEvidenceKeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	sdkEvidenceTypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	sdkFeeGrantKeeper "github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	sdkFeeGrantModule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	sdkFeegrantModule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutilTypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	sdkGovKeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	sdkGovTypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	sdkMintKeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	mintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	sdkMintTypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsClient "github.com/cosmos/cosmos-sdk/x/params/client"
	sdkParamsKeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramsTypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramsProposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingKeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingTypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	sdkStakingKeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	sdkStakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeClient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	sdkUpgradeKeeper "github.com/cosmos/cosmos-sdk/x/upgrade/keeper"
	sdkUpgradeTypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
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
	ibcCore "github.com/cosmos/ibc-go/v3/modules/core"
	ibcCoreClient "github.com/cosmos/ibc-go/v3/modules/core/02-client"
	ibcClient "github.com/cosmos/ibc-go/v3/modules/core/02-client/client"
	ibcClientTypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	ibcConnectionTypes "github.com/cosmos/ibc-go/v3/modules/core/03-connection/types"
	ibcTypes "github.com/cosmos/ibc-go/v3/modules/core/05-port/types"
	ibcHost "github.com/cosmos/ibc-go/v3/modules/core/24-host"
	sdkIBCKeeper "github.com/cosmos/ibc-go/v3/modules/core/keeper"
	"github.com/gogo/protobuf/grpc"
	"github.com/gorilla/mux"
	"github.com/rakyll/statik/fs"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	tendermintJSON "github.com/tendermint/tendermint/libs/json"
	tendermintLog "github.com/tendermint/tendermint/libs/log"
	tendermintOS "github.com/tendermint/tendermint/libs/os"
	tendermintProto "github.com/tendermint/tendermint/proto/tendermint/types"
	tendermintDB "github.com/tendermint/tm-db"
	"honnef.co/go/tools/version"

	applicationParams "github.com/persistenceOne/persistenceCore/app/params"
	"github.com/persistenceOne/persistenceCore/x/halving"
)

var DefaultNodeHome string

var ModuleAccountPermissions = map[string][]string{
	authTypes.FeeCollectorName:     nil,
	distributionTypes.ModuleName:   nil,
	icaTypes.ModuleName:            nil,
	mintTypes.ModuleName:           {authTypes.Minter},
	stakingTypes.BondedPoolName:    {authTypes.Burner, authTypes.Staking},
	stakingTypes.NotBondedPoolName: {authTypes.Burner, authTypes.Staking},
	govTypes.ModuleName:            {authTypes.Burner},
	ibcTransferTypes.ModuleName:    {authTypes.Minter, authTypes.Burner},
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
		paramsClient.ProposalHandler,
		distributionClient.ProposalHandler,
		upgradeClient.ProposalHandler,
		upgradeClient.CancelProposalHandler,
		ibcClient.UpdateClientProposalHandler,
		ibcClient.UpgradeProposalHandler,
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	sdkFeegrantModule.AppModuleBasic{},
	sdkAuthzModule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	halving.AppModuleBasic{},
	ica.AppModuleBasic{},
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
	moduleManager      *sdkTypesModule.Manager
	configurator       sdkTypesModule.Configurator
	simulationManager  *sdkTypesModule.SimulationManager

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

func (app Application) Name() string {
	return app.BaseApp.Name()
}

func (app Application) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

func (app Application) BeginBlocker(ctx sdkTypes.Context, req abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	return app.moduleManager.BeginBlock(ctx, req)
}

func (app Application) EndBlocker(ctx sdkTypes.Context, req abciTypes.RequestEndBlock) abciTypes.ResponseEndBlock {
	return app.moduleManager.EndBlock(ctx, req)
}

func (app Application) InitChainer(ctx sdkTypes.Context, req abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
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
			if _, Error := sdkTypes.ValAddressFromBech32(address); Error != nil {
				panic(Error)
			}

			whiteListMap[address] = true
		}

		app.CrisisKeeper.AssertInvariants(context)

		app.StakingKeeper.IterateValidators(context, func(_ int64, val sdkStakingTypes.ValidatorI) (stop bool) {
			_, _ = app.DistributionKeeper.WithdrawValidatorCommission(context, val.GetOperator())
			return false
		})

		delegations := app.StakingKeeper.GetAllDelegations(context)
		for _, delegation := range delegations {
			validatorAddress, Error := sdkTypes.ValAddressFromBech32(delegation.ValidatorAddress)
			if Error != nil {
				panic(Error)
			}

			delegatorAddress, Error := sdkTypes.AccAddressFromBech32(delegation.DelegatorAddress)
			if Error != nil {
				panic(Error)
			}

			_, _ = app.DistributionKeeper.WithdrawDelegationRewards(context, delegatorAddress, validatorAddress)
		}

		app.DistributionKeeper.DeleteAllValidatorSlashEvents(context)

		app.DistributionKeeper.DeleteAllValidatorHistoricalRewards(context)

		height := context.BlockHeight()
		context = context.WithBlockHeight(0)

		app.StakingKeeper.IterateValidators(context, func(_ int64, val sdkStakingTypes.ValidatorI) (stop bool) {

			scraps := app.DistributionKeeper.GetValidatorOutstandingRewardsCoins(context, val.GetOperator())
			feePool := app.DistributionKeeper.GetFeePool(context)
			feePool.CommunityPool = feePool.CommunityPool.Add(scraps...)
			app.DistributionKeeper.SetFeePool(context, feePool)

			app.DistributionKeeper.Hooks().AfterValidatorCreated(context, val.GetOperator())
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

			app.DistributionKeeper.Hooks().BeforeDelegationCreated(context, delegatorAddress, validatorAddress)
			app.DistributionKeeper.Hooks().AfterDelegationModified(context, delegatorAddress, validatorAddress)
		}

		context = context.WithBlockHeight(height)

		app.StakingKeeper.IterateRedelegations(context, func(_ int64, redelegation sdkStakingTypes.Redelegation) (stop bool) {
			for i := range redelegation.Entries {
				redelegation.Entries[i].CreationHeight = 0
			}
			app.StakingKeeper.SetRedelegation(context, redelegation)
			return false
		})

		app.StakingKeeper.IterateUnbondingDelegations(context, func(_ int64, unbondingDelegation sdkStakingTypes.UnbondingDelegation) (stop bool) {
			for i := range unbondingDelegation.Entries {
				unbondingDelegation.Entries[i].CreationHeight = 0
			}
			app.StakingKeeper.SetUnbondingDelegation(context, unbondingDelegation)
			return false
		})

		store := context.KVStore(app.keys[sdkStakingTypes.StoreKey])
		kvStoreReversePrefixIterator := sdkTypes.KVStoreReversePrefixIterator(store, sdkStakingTypes.ValidatorsKey)
		counter := int16(0)

		for ; kvStoreReversePrefixIterator.Valid(); kvStoreReversePrefixIterator.Next() {
			addr := sdkTypes.ValAddress(kvStoreReversePrefixIterator.Key()[1:])
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
			log.Fatal(Error)
		}

		app.SlashingKeeper.IterateValidatorSigningInfos(
			context,
			func(validatorConsAddress sdkTypes.ConsAddress, validatorSigningInfo slashingTypes.ValidatorSigningInfo) (stop bool) {
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

func (app Application) SimulationManager() *sdkTypesModule.SimulationManager {
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

func NewApplication(
	applicationName string,
	encodingConfiguration applicationParams.EncodingConfig,
	moduleAccountPermissions map[string][]string,
	logger tendermintLog.Logger,
	db tendermintDB.DB,
	traceStore io.Writer,
	loadLatest bool,
	invCheckPeriod uint,
	skipUpgradeHeights map[int64]bool,
	home string,
	applicationOptions serverTypes.AppOptions,
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

	keys := sdkTypes.NewKVStoreKeys(
		authTypes.StoreKey, sdkBankTypes.StoreKey, sdkStakingTypes.StoreKey,
		sdkMintTypes.StoreKey, sdkDistributionTypes.StoreKey, slashingTypes.StoreKey,
		sdkGovTypes.StoreKey, paramsTypes.StoreKey, ibcHost.StoreKey, sdkUpgradeTypes.StoreKey,
		sdkEvidenceTypes.StoreKey, ibcTransferTypes.StoreKey, sdkCapabilityTypes.StoreKey,
		feegrant.StoreKey, sdkAuthzKeeper.StoreKey, icaHostTypes.StoreKey, halving.StoreKey,
	)

	transientStoreKeys := sdkTypes.NewTransientStoreKeys(paramsTypes.TStoreKey)
	memoryKeys := sdkTypes.NewMemoryStoreKeys(sdkCapabilityTypes.MemStoreKey)

	app := &Application{
		BaseApp:           baseApp,
		legacyAmino:       legacyAmino,
		applicationCodec:  applicationCodec,
		interfaceRegistry: interfaceRegistry,
		keys:              keys,
	}

	app.ParamsKeeper = sdkParamsKeeper.NewKeeper(
		applicationCodec,
		legacyAmino,
		keys[paramsTypes.StoreKey],
		transientStoreKeys[paramsTypes.TStoreKey],
	)
	app.BaseApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(sdkParamsKeeper.ConsensusParamsKeyTable()))

	app.CapabilityKeeper = sdkCapabilityKeeper.NewKeeper(applicationCodec, keys[sdkCapabilityTypes.StoreKey], memoryKeys[sdkCapabilityTypes.MemStoreKey])
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcHost.ModuleName)
	scopedICAHostKeeper := app.CapabilityKeeper.ScopeToModule(icaHostTypes.SubModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibcTransferTypes.ModuleName)
	app.CapabilityKeeper.Seal()

	app.AccountKeeper = authKeeper.NewAccountKeeper(
		applicationCodec,
		keys[authTypes.StoreKey],
		app.ParamsKeeper.Subspace(authTypes.ModuleName),
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

	app.BankKeeper = sdkBankKeeper.NewBaseKeeper(
		applicationCodec,
		keys[sdkBankTypes.StoreKey],
		app.AccountKeeper,
		app.ParamsKeeper.Subspace(sdkBankTypes.ModuleName),
		blacklistedAddresses,
	)

	app.AuthzKeeper = sdkAuthzKeeper.NewKeeper(
		keys[sdkAuthzKeeper.StoreKey],
		applicationCodec,
		app.BaseApp.MsgServiceRouter(),
	)

	app.FeegrantKeeper = sdkFeeGrantKeeper.NewKeeper(
		applicationCodec,
		keys[feegrant.StoreKey],
		app.AccountKeeper,
	)

	stakingKeeper := sdkStakingKeeper.NewKeeper(
		applicationCodec,
		keys[sdkStakingTypes.StoreKey],
		app.AccountKeeper,
		app.BankKeeper,
		app.ParamsKeeper.Subspace(sdkStakingTypes.ModuleName),
	)

	app.MintKeeper = sdkMintKeeper.NewKeeper(
		applicationCodec,
		keys[sdkMintTypes.StoreKey],
		app.ParamsKeeper.Subspace(sdkMintTypes.ModuleName),
		&stakingKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		authTypes.FeeCollectorName,
	)

	app.DistributionKeeper = sdkDistributionKeeper.NewKeeper(
		applicationCodec,
		keys[sdkDistributionTypes.StoreKey],
		app.ParamsKeeper.Subspace(sdkDistributionTypes.ModuleName),
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
		app.ParamsKeeper.Subspace(slashingTypes.ModuleName),
	)
	app.CrisisKeeper = sdkCrisisKeeper.NewKeeper(
		app.ParamsKeeper.Subspace(sdkCrisisTypes.ModuleName),
		invCheckPeriod,
		app.BankKeeper,
		authTypes.FeeCollectorName,
	)
	app.UpgradeKeeper = sdkUpgradeKeeper.NewKeeper(
		skipUpgradeHeights,
		keys[sdkUpgradeTypes.StoreKey],
		applicationCodec,
		home,
		app.BaseApp,
	)

	app.HalvingKeeper = halving.NewKeeper(
		keys[halving.StoreKey],
		app.ParamsKeeper.Subspace(halving.DefaultParamspace),
		app.MintKeeper,
	)

	app.StakingKeeper = *stakingKeeper.SetHooks(
		sdkStakingTypes.NewMultiStakingHooks(app.DistributionKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	app.IBCKeeper = sdkIBCKeeper.NewKeeper(
		applicationCodec,
		keys[ibcHost.StoreKey],
		app.ParamsKeeper.Subspace(ibcHost.ModuleName),
		app.StakingKeeper,
		app.UpgradeKeeper,
		scopedIBCKeeper,
	)

	govRouter := sdkGovTypes.NewRouter()
	govRouter.AddRoute(
		sdkGovTypes.RouterKey,
		sdkGovTypes.ProposalHandler,
	).AddRoute(
		paramsProposal.RouterKey,
		params.NewParamChangeProposalHandler(app.ParamsKeeper),
	).AddRoute(
		sdkDistributionTypes.RouterKey,
		distribution.NewCommunityPoolSpendProposalHandler(app.DistributionKeeper),
	).AddRoute(
		sdkUpgradeTypes.RouterKey,
		upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper),
	).AddRoute(ibcClientTypes.RouterKey, ibcCoreClient.NewClientProposalHandler(app.IBCKeeper.ClientKeeper))

	app.GovKeeper = sdkGovKeeper.NewKeeper(
		applicationCodec,
		keys[sdkGovTypes.StoreKey],
		app.ParamsKeeper.Subspace(sdkGovTypes.ModuleName).WithKeyTable(sdkGovTypes.ParamKeyTable()),
		app.AccountKeeper,
		app.BankKeeper,
		&stakingKeeper,
		govRouter,
	)

	app.TransferKeeper = ibcTransferKeeper.NewKeeper(
		applicationCodec,
		keys[ibcTransferTypes.StoreKey],
		app.ParamsKeeper.Subspace(ibcTransferTypes.ModuleName),
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
		app.ParamsKeeper.Subspace(icaHostTypes.SubModuleName),
		app.IBCKeeper.ChannelKeeper,
		&app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(),
	)

	icaModule := ica.NewAppModule(nil, &app.ICAHostKeeper)
	icaHostIBCModule := icaHost.NewIBCModule(app.ICAHostKeeper)

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcTypes.NewRouter()
	ibcRouter.AddRoute(icaHostTypes.SubModuleName, icaHostIBCModule).
		AddRoute(ibcTransferTypes.ModuleName, transferIBCModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	evidenceKeeper := sdkEvidenceKeeper.NewKeeper(
		applicationCodec,
		keys[sdkEvidenceTypes.StoreKey],
		&app.StakingKeeper,
		app.SlashingKeeper,
	)
	app.EvidenceKeeper = *evidenceKeeper

	/****  Module Options ****/
	var skipGenesisInvariants = false

	opt := applicationOptions.Get(crisis.FlagSkipGenesisInvariants)
	if opt, ok := opt.(bool); ok {
		skipGenesisInvariants = opt
	}

	app.moduleManager = sdkTypesModule.NewManager(
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
		sdkFeeGrantModule.NewAppModule(applicationCodec, app.AccountKeeper, app.BankKeeper, app.FeegrantKeeper, app.interfaceRegistry),
		sdkAuthzModule.NewAppModule(applicationCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibcCore.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		halving.NewAppModule(applicationCodec, app.HalvingKeeper),
		transferModule,
		icaModule,
	)

	app.moduleManager.SetOrderBeginBlockers(
		sdkUpgradeTypes.ModuleName,
		sdkCapabilityTypes.ModuleName,
		sdkCrisisTypes.ModuleName,
		sdkGovTypes.ModuleName,
		sdkStakingTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcHost.ModuleName,
		icaTypes.ModuleName,
		authTypes.ModuleName,
		sdkBankTypes.ModuleName,
		sdkDistributionTypes.ModuleName,
		slashingTypes.ModuleName,
		sdkMintTypes.ModuleName,
		genutilTypes.ModuleName,
		sdkEvidenceTypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramsTypes.ModuleName,
		vestingTypes.ModuleName,
		halving.ModuleName,
	)
	app.moduleManager.SetOrderEndBlockers(
		sdkCrisisTypes.ModuleName,
		sdkGovTypes.ModuleName,
		sdkStakingTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcHost.ModuleName,
		icaTypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		sdkCapabilityTypes.ModuleName,
		authTypes.ModuleName,
		sdkBankTypes.ModuleName,
		sdkDistributionTypes.ModuleName,
		slashingTypes.ModuleName,
		sdkMintTypes.ModuleName,
		genutilTypes.ModuleName,
		sdkEvidenceTypes.ModuleName,
		paramsTypes.ModuleName,
		sdkUpgradeTypes.ModuleName,
		vestingTypes.ModuleName,
		halving.ModuleName,
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.moduleManager.SetOrderInitGenesis(
		sdkCapabilityTypes.ModuleName,
		sdkBankTypes.ModuleName,
		sdkDistributionTypes.ModuleName,
		sdkStakingTypes.ModuleName,
		slashingTypes.ModuleName,
		sdkGovTypes.ModuleName,
		sdkMintTypes.ModuleName,
		sdkCrisisTypes.ModuleName,
		ibcTransferTypes.ModuleName,
		ibcHost.ModuleName,
		icaTypes.ModuleName,
		sdkEvidenceTypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		authTypes.ModuleName,
		genutilTypes.ModuleName,
		paramsTypes.ModuleName,
		sdkUpgradeTypes.ModuleName,
		vestingTypes.ModuleName,
		halving.ModuleName,
	)

	app.moduleManager.RegisterInvariants(&app.CrisisKeeper)
	app.moduleManager.RegisterRoutes(app.BaseApp.Router(), app.BaseApp.QueryRouter(), encodingConfiguration.Amino)
	app.configurator = sdkTypesModule.NewConfigurator(app.applicationCodec, app.BaseApp.MsgServiceRouter(), app.BaseApp.GRPCQueryRouter())
	app.moduleManager.RegisterServices(app.configurator)

	simulationManager := sdkTypesModule.NewSimulationManager(
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
		sdkAuthzModule.NewAppModule(applicationCodec, app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		sdkFeeGrantModule.NewAppModule(applicationCodec, app.AccountKeeper, app.BankKeeper, app.FeegrantKeeper, app.interfaceRegistry),
		ibcCore.NewAppModule(app.IBCKeeper),
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
			IBCKeeper: app.IBCKeeper,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	app.BaseApp.SetAnteHandler(anteHandler)
	app.BaseApp.SetInitChainer(app.InitChainer)
	app.BaseApp.SetBeginBlocker(app.moduleManager.BeginBlock)
	app.BaseApp.SetEndBlocker(app.moduleManager.EndBlock)

	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeName,
		func(ctx sdkTypes.Context, _ sdkUpgradeTypes.Plan, fromVM sdkTypesModule.VersionMap) (sdkTypesModule.VersionMap, error) {
			app.IBCKeeper.ConnectionKeeper.SetParams(ctx, ibcConnectionTypes.DefaultParams())
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
			return app.moduleManager.RunMigrations(ctx, app.configurator, fromVM)

		},
	)

	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(fmt.Sprintf("failed to read upgrade info from disk %s", err))
	}

	if upgradeInfo.Name == UpgradeName && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		storeUpgrades := storeTypes.StoreUpgrades{
			Added: []string{authz.ModuleName, feegrant.ModuleName},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.BaseApp.SetStoreLoader(sdkUpgradeTypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}

	if loadLatest {
		if err := app.BaseApp.LoadLatestVersion(); err != nil {
			tendermintOS.Exit(err.Error())
		}
	}
	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	return app
}
