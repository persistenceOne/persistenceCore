package zone

import (
	"encoding/json"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	tendermintDB "github.com/tendermint/tm-db"
)

const applicationName = "persistenceOne"

var DefaultClientHome = os.ExpandEnv("$HOME/.hubClient")

var DefaultNodeHome = os.ExpandEnv("$HOME/.hubNode")

func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

type PersistenceOneApplication struct {
	*baseapp.BaseApp
	cdc *codec.Codec

	invCheckPeriod uint

	keyMain          *sdk.KVStoreKey
	keyAccount       *sdk.KVStoreKey
	keyFeeCollection *sdk.KVStoreKey
	keyStaking       *sdk.KVStoreKey
	tkeyStaking      *sdk.TransientStoreKey
	keyDistribution  *sdk.KVStoreKey
	tkeyDistribution *sdk.TransientStoreKey
	keyParameter     *sdk.KVStoreKey
	tkeyParameter    *sdk.TransientStoreKey
	keySlashing      *sdk.KVStoreKey
	keySupply        *sdk.KVStoreKey
	keyMint          *sdk.KVStoreKey
	keyGovernment    *sdk.KVStoreKey

	accountKeeper      auth.AccountKeeper
	bankKeeper         bank.Keeper
	stakingKeeper      staking.Keeper
	slashingKeeper     slashing.Keeper
	distributionKeeper distribution.Keeper
	paramsKeeper       params.Keeper
	supplyKeeper       supply.Keeper
	crisisKeeper       crisis.Keeper

	moduleManager *module.Manager
}

func NewPersistenceOneApplicaiton(logger log.Logger, db tendermintDB.DB, traceStore io.Writer, loadLatest bool,
	invCheckPeriod uint, baseAppOptions ...func(*baseapp.BaseApp)) *PersistenceOneApplication {
	cdc := MakeCodec()
	baseApp := baseapp.NewBaseApp(applicationName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)
	baseApp.SetCommitMultiStoreTracer(traceStore)
	baseApp.SetAppVersion(version.Version)

	application := &PersistenceOneApplication{
		BaseApp:          baseApp,
		cdc:              cdc,
		invCheckPeriod:   invCheckPeriod,
		keyMain:          sdk.NewKVStoreKey(baseapp.MainStoreKey),
		keyAccount:       sdk.NewKVStoreKey(auth.StoreKey),
		keyStaking:       sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:      sdk.NewTransientStoreKey(staking.TStoreKey),
		keyDistribution:  sdk.NewKVStoreKey(distribution.StoreKey),
		tkeyDistribution: sdk.NewTransientStoreKey(distribution.TStoreKey),
		keyParameter:     sdk.NewKVStoreKey(params.StoreKey),
		tkeyParameter:    sdk.NewTransientStoreKey(params.TStoreKey),
		keySlashing:      sdk.NewKVStoreKey(slashing.StoreKey),
	}

	application.paramsKeeper = params.NewKeeper(application.cdc, application.keyParameter, application.tkeyParameter, params.DefaultCodespace)

	authSubspace := application.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSubspace := application.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := application.paramsKeeper.Subspace(staking.DefaultParamspace)
	distributionSubspace := application.paramsKeeper.Subspace(distribution.DefaultParamspace)
	slashingSubspace := application.paramsKeeper.Subspace(slashing.DefaultParamspace)

	application.accountKeeper = auth.NewAccountKeeper(
		application.cdc,
		application.keyAccount,
		authSubspace,
		auth.ProtoBaseAccount,
	)

	application.bankKeeper = bank.NewBaseKeeper(
		application.accountKeeper,
		bankSubspace,
		bank.DefaultCodespace,
	)

	stakingKeeper := staking.NewKeeper(
		application.cdc,
		application.keyStaking,
		application.tkeyStaking,
		application.supplyKeeper,
		stakingSubspace,
		staking.DefaultCodespace,
	)

	application.distributionKeeper = distribution.NewKeeper(
		application.cdc,
		application.keyDistribution,
		distributionSubspace,
		application.stakingKeeper,
		application.supplyKeeper,
		distribution.DefaultCodespace,
		auth.FeeCollectorName,
	)

	application.slashingKeeper = slashing.NewKeeper(
		application.cdc,
		application.keySlashing,
		&stakingKeeper,
		slashingSubspace,
		slashing.DefaultCodespace,
	)

	application.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
			application.distributionKeeper.Hooks(),
			application.slashingKeeper.Hooks()),
	)

	application.moduleManager.SetOrderInitGenesis(
		distribution.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		genutil.ModuleName,
	)

	application.moduleManager.SetOrderBeginBlockers(
		mint.ModuleName,
		distribution.ModuleName,
		slashing.ModuleName,
	)

	application.moduleManager.SetOrderEndBlockers(
		gov.ModuleName,
		staking.ModuleName,
	)

	application.moduleManager.SetOrderInitGenesis(
		genaccounts.ModuleName,
		supply.ModuleName,
		distribution.ModuleName,
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		slashing.ModuleName,
		gov.ModuleName,
		mint.ModuleName,
		crisis.ModuleName,
		genutil.ModuleName,
	)

	application.moduleManager.RegisterInvariants(&application.crisisKeeper)
	application.moduleManager.RegisterRoutes(application.Router(), application.QueryRouter())

	application.MountStores(
		application.keyMain,
		application.keyAccount,
		application.keySupply,
		application.keyStaking,
		application.keyMint,
		application.keyDistribution,
		application.keySlashing,
		application.keyGovernment,
		application.keyParameter,
		application.tkeyParameter,
		application.tkeyStaking,
		application.keyDistribution,
	)

	application.SetInitChainer(application.InitChainer)
	application.SetBeginBlocker(application.BeginBlocker)
	application.SetAnteHandler(auth.NewAnteHandler(application.accountKeeper, application.supplyKeeper, auth.DefaultSigVerificationGasConsumer))
	application.SetEndBlocker(application.EndBlocker)

	if loadLatest {
		err := application.LoadLatestVersion(application.keyMain)
		if err != nil {
			common.Exit(err.Error())
		}
	}
	return application
}

type GenesisState map[string]json.RawMessage

func (app *PersistenceOneApplication) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.moduleManager.BeginBlock(ctx, req)
}

func (app *PersistenceOneApplication) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.moduleManager.EndBlock(ctx, req)
}

func (app *PersistenceOneApplication) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	return app.moduleManager.InitGenesis(ctx, genesisState)
}

func (app *PersistenceOneApplication) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}
