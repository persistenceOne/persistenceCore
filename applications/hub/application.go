package hub

import (
	"io"
	"os"

	abciTypes "github.com/tendermint/tendermint/abci/types"
	tendermintCommon "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tendermintDB "github.com/tendermint/tm-db"
	"honnef.co/go/tools/version"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsClient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

const applicationName = "CommitHubApplication"

var DefaultClientHome = os.ExpandEnv("$HOME/.hubClient")
var DefaultNodeHome = os.ExpandEnv("$HOME/.hubNode")
var moduleAccountPermissions = map[string][]string{
	auth.FeeCollectorName:     nil,
	distribution.ModuleName:   nil,
	mint.ModuleName:           {supply.Minter},
	staking.BondedPoolName:    {supply.Burner, supply.Staking},
	staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	gov.ModuleName:            {supply.Burner},
}
var ModuleBasics = module.NewBasicManager(
	genaccounts.AppModuleBasic{},
	genutil.AppModuleBasic{},
	auth.AppModuleBasic{},
	bank.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.NewAppModuleBasic(paramsClient.ProposalHandler, distribution.ProposalHandler),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	supply.AppModuleBasic{},
)

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	bank.RegisterCodec(cdc)
	staking.RegisterCodec(cdc)
	distribution.RegisterCodec(cdc)
	slashing.RegisterCodec(cdc)
	gov.RegisterCodec(cdc)
	auth.RegisterCodec(cdc)
	crisis.RegisterCodec(cdc)
	sdkTypes.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

type CommitHubApplication struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	keys          map[string]*sdkTypes.KVStoreKey
	transientKeys map[string]*sdkTypes.TransientStoreKey

	accountKeeper      auth.AccountKeeper
	bankKeeper         bank.Keeper
	supplyKeeper       supply.Keeper
	stakingKeeper      staking.Keeper
	slashingKeeper     slashing.Keeper
	mintKeeper         mint.Keeper
	distributionKeeper distribution.Keeper
	govKeeper          gov.Keeper
	crisisKeeper       crisis.Keeper
	parameterKeeper    params.Keeper

	moduleManager *module.Manager
}

func NewCommitHubApplication(
	logger log.Logger,
	db tendermintDB.DB,
	traceStore io.Writer,
	loadLatest bool,
	invCheckPeriod uint,
	baseAppOptions ...func(*baseapp.BaseApp),
) *CommitHubApplication {

	cdc := MakeCodec()

	baseApp := baseapp.NewBaseApp(
		applicationName,
		logger,
		db,
		auth.DefaultTxDecoder(cdc),
		baseAppOptions...,
	)
	baseApp.SetCommitMultiStoreTracer(traceStore)
	baseApp.SetAppVersion(version.Version)

	keys := sdkTypes.NewKVStoreKeys(
		baseapp.MainStoreKey,
		auth.StoreKey,
		staking.StoreKey,
		supply.StoreKey,
		mint.StoreKey,
		distribution.StoreKey,
		slashing.StoreKey,
		gov.StoreKey,
		params.StoreKey,
	)
	transientKeys := sdkTypes.NewTransientStoreKeys(
		staking.TStoreKey,
		params.TStoreKey,
	)

	var application = &CommitHubApplication{
		BaseApp:        baseApp,
		cdc:            cdc,
		invCheckPeriod: invCheckPeriod,
		keys:           keys,
		transientKeys:  transientKeys,
	}

	application.parameterKeeper = params.NewKeeper(
		application.cdc,
		keys[params.StoreKey],
		transientKeys[params.TStoreKey],
		params.DefaultCodespace,
	)
	authSubspace := application.parameterKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := application.parameterKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := application.parameterKeeper.Subspace(staking.DefaultParamspace)
	mintSubspace := application.parameterKeeper.Subspace(mint.DefaultParamspace)
	distributionSubspace := application.parameterKeeper.Subspace(distribution.DefaultParamspace)
	slashingSubspace := application.parameterKeeper.Subspace(slashing.DefaultParamspace)
	govSubspace := application.parameterKeeper.Subspace(gov.DefaultParamspace)
	crisisSubspace := application.parameterKeeper.Subspace(crisis.DefaultParamspace)

	application.accountKeeper = auth.NewAccountKeeper(
		application.cdc,
		keys[auth.StoreKey],
		authSubspace,
		auth.ProtoBaseAccount,
	)

	application.bankKeeper = bank.NewBaseKeeper(
		application.accountKeeper,
		bankSubspace,
		bank.DefaultCodespace,
		application.ModuleAccountAddress(),
	)

	application.supplyKeeper = supply.NewKeeper(
		application.cdc,
		keys[supply.StoreKey],
		application.accountKeeper,
		application.bankKeeper,
		moduleAccountPermissions,
	)

	stakingKeeper := staking.NewKeeper(
		application.cdc,
		keys[staking.StoreKey],
		transientKeys[staking.TStoreKey],
		application.supplyKeeper,
		stakingSubspace,
		staking.DefaultCodespace,
	)
	application.mintKeeper = mint.NewKeeper(
		application.cdc,
		keys[mint.StoreKey],
		mintSubspace,
		&stakingKeeper,
		application.supplyKeeper,
		auth.FeeCollectorName,
	)
	application.distributionKeeper = distribution.NewKeeper(
		application.cdc,
		keys[distribution.StoreKey],
		distributionSubspace,
		&stakingKeeper,
		application.supplyKeeper,
		distribution.DefaultCodespace,
		auth.FeeCollectorName,
		application.ModuleAccountAddress(),
	)
	application.slashingKeeper = slashing.NewKeeper(
		application.cdc,
		keys[slashing.StoreKey],
		&stakingKeeper,
		slashingSubspace,
		slashing.DefaultCodespace,
	)
	application.crisisKeeper = crisis.NewKeeper(
		crisisSubspace,
		invCheckPeriod,
		application.supplyKeeper,
		auth.FeeCollectorName,
	)
	govRouter := gov.NewRouter()
	govRouter.AddRoute(
		gov.RouterKey,
		gov.ProposalHandler,
	).AddRoute(
		params.RouterKey,
		params.NewParamChangeProposalHandler(application.parameterKeeper),
	).AddRoute(
		distribution.RouterKey,
		distribution.NewCommunityPoolSpendProposalHandler(application.distributionKeeper),
	)
	application.govKeeper = gov.NewKeeper(
		application.cdc,
		keys[gov.StoreKey],
		application.parameterKeeper,
		govSubspace,
		application.supplyKeeper,
		&stakingKeeper,
		gov.DefaultCodespace,
		govRouter,
	)

	application.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
			application.distributionKeeper.Hooks(),
			application.slashingKeeper.Hooks(),
		),
	)

	application.moduleManager = module.NewManager(
		genaccounts.NewAppModule(application.accountKeeper),
		genutil.NewAppModule(application.accountKeeper, application.stakingKeeper, application.BaseApp.DeliverTx),
		auth.NewAppModule(application.accountKeeper),
		bank.NewAppModule(application.bankKeeper, application.accountKeeper),
		crisis.NewAppModule(&application.crisisKeeper),
		supply.NewAppModule(application.supplyKeeper, application.accountKeeper),
		distribution.NewAppModule(application.distributionKeeper, application.supplyKeeper),
		gov.NewAppModule(application.govKeeper, application.supplyKeeper),
		mint.NewAppModule(application.mintKeeper),
		slashing.NewAppModule(application.slashingKeeper, application.stakingKeeper),
		staking.NewAppModule(application.stakingKeeper, application.distributionKeeper, application.accountKeeper, application.supplyKeeper),
	)

	application.moduleManager.SetOrderBeginBlockers(mint.ModuleName, distribution.ModuleName, slashing.ModuleName)

	application.moduleManager.SetOrderEndBlockers(crisis.ModuleName, gov.ModuleName, staking.ModuleName)

	application.moduleManager.SetOrderInitGenesis(
		genaccounts.ModuleName, distribution.ModuleName, staking.ModuleName,
		auth.ModuleName, bank.ModuleName, slashing.ModuleName, gov.ModuleName,
		mint.ModuleName, supply.ModuleName, crisis.ModuleName, genutil.ModuleName,
	)

	application.moduleManager.RegisterInvariants(&application.crisisKeeper)
	application.moduleManager.RegisterRoutes(application.Router(), application.QueryRouter())

	application.MountKVStores(keys)
	application.MountTransientStores(transientKeys)

	application.SetInitChainer(application.InitChainer)
	application.SetBeginBlocker(application.BeginBlocker)
	application.SetAnteHandler(auth.NewAnteHandler(application.accountKeeper, application.supplyKeeper, auth.DefaultSigVerificationGasConsumer))
	application.SetEndBlocker(application.EndBlocker)

	if loadLatest {
		err := application.LoadLatestVersion(application.keys[baseapp.MainStoreKey])
		if err != nil {
			tendermintCommon.Exit(err.Error())
		}
	}

	return application
}
func (application *CommitHubApplication) BeginBlocker(ctx sdkTypes.Context, req abciTypes.RequestBeginBlock) abciTypes.ResponseBeginBlock {
	return application.moduleManager.BeginBlock(ctx, req)
}
func (application *CommitHubApplication) EndBlocker(ctx sdkTypes.Context, req abciTypes.RequestEndBlock) abciTypes.ResponseEndBlock {
	return application.moduleManager.EndBlock(ctx, req)
}
func (application *CommitHubApplication) InitChainer(ctx sdkTypes.Context, req abciTypes.RequestInitChain) abciTypes.ResponseInitChain {
	var genesisState GenesisState
	application.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return application.moduleManager.InitGenesis(ctx, genesisState)
}
func (application *CommitHubApplication) LoadHeight(height int64) error {
	return application.LoadVersion(height, application.keys[baseapp.MainStoreKey])
}
func (application *CommitHubApplication) ModuleAccountAddress() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range moduleAccountPermissions {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}
