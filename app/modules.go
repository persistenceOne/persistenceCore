package app

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authsimulation "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	group "github.com/cosmos/cosmos-sdk/x/group"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	packetforward "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/router"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/router/types"
	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	icatypes "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/persistenceOne/persistence-sdk/v2/x/epochs"
	epochstypes "github.com/persistenceOne/persistence-sdk/v2/x/epochs/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/halving"
	ibchookertypes "github.com/persistenceOne/persistence-sdk/v2/x/ibchooker/types"
	interchainquerytypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/oracle"
	oracletypes "github.com/persistenceOne/persistence-sdk/v2/x/oracle/types"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	lscosmostypes "github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"
	"github.com/skip-mev/pob/x/builder"
	buildertypes "github.com/skip-mev/pob/x/builder/types"

	appparams "github.com/persistenceOne/persistenceCore/v8/app/params"
)

var ModuleAccountPermissions = map[string][]string{
	authtypes.FeeCollectorName:                    nil,
	distributiontypes.ModuleName:                  nil,
	icatypes.ModuleName:                           nil,
	minttypes.ModuleName:                          {authtypes.Minter},
	stakingtypes.BondedPoolName:                   {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName:                {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:                           {authtypes.Burner},
	ibctransfertypes.ModuleName:                   {authtypes.Minter, authtypes.Burner},
	ibcfeetypes.ModuleName:                        nil,
	wasm.ModuleName:                               {authtypes.Burner},
	lscosmostypes.ModuleName:                      {authtypes.Minter, authtypes.Burner},
	lscosmostypes.DepositModuleAccount:            nil,
	lscosmostypes.DelegationModuleAccount:         nil,
	lscosmostypes.RewardModuleAccount:             nil,
	lscosmostypes.UndelegationModuleAccount:       nil,
	lscosmostypes.RewardBoosterModuleAccount:      nil,
	oracletypes.ModuleName:                        nil,
	liquidstakeibctypes.ModuleName:                {authtypes.Minter, authtypes.Burner},
	liquidstakeibctypes.DepositModuleAccount:      nil,
	liquidstakeibctypes.UndelegationModuleAccount: {authtypes.Burner},
	buildertypes.ModuleName:                       nil,
}

var receiveAllowedMAcc = map[string]bool{
	lscosmostypes.UndelegationModuleAccount:       true,
	lscosmostypes.DelegationModuleAccount:         true,
	liquidstakeibctypes.DepositModuleAccount:      true,
	liquidstakeibctypes.UndelegationModuleAccount: true,
}

func appModules(
	app *Application,
	encodingConfig appparams.EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	appCodec := encodingConfig.Codec

	return []module.AppModule{
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, *app.AccountKeeper, authsimulation.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(*app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, *app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, false),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, *app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)), // nil -> SDK's default inflation function.
		slashing.NewAppModule(appCodec, *app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName)),
		distribution.NewAppModule(appCodec, *app.DistributionKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper, app.GetSubspace(distributiontypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(*app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, *app.FeegrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, *app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		groupmodule.NewAppModule(appCodec, *app.GroupKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		ibcfee.NewAppModule(*app.IBCFeeKeeper),
		packetforward.NewAppModule(app.PacketForwardKeeper),
		params.NewAppModule(*app.ParamsKeeper),
		halving.NewAppModule(appCodec, *app.HalvingKeeper),
		app.TransferModule,
		app.IBCTransferHooksMiddleware,
		ibchooks.NewAppModule(*app.AccountKeeper),
		ica.NewAppModule(app.ICAControllerKeeper, app.ICAHostKeeper),
		wasm.NewAppModule(appCodec, app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasm.ModuleName)),
		epochs.NewAppModule(*app.EpochsKeeper),
		app.InterchainQueryModule,
		app.LSCosmosModule,
		liquidstakeibc.NewAppModule(*app.LiquidStakeIBCKeeper),
		oracle.NewAppModule(appCodec, *app.OracleKeeper, app.AccountKeeper, app.BankKeeper),
		builder.NewAppModule(appCodec, *app.BuilderKeeper),
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)), // always be last to make sure that it checks for all invariants and not only part of them
	}
}

func overrideSimulationModules(
	app *Application,
	encodingConfig appparams.EncodingConfig,
	_ bool,
) map[string]module.AppModuleSimulation {
	appCodec := encodingConfig.Codec

	return map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(appCodec, *app.AccountKeeper, authsimulation.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
	}
}

func orderBeginBlockers() []string {
	return []string{
		upgradetypes.ModuleName,
		epochstypes.ModuleName,
		capabilitytypes.ModuleName,
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		ibcfeetypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distributiontypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		buildertypes.ModuleName,
		consensusparamtypes.ModuleName,
		halving.ModuleName,
		ibchookstypes.ModuleName,
		packetforwardtypes.ModuleName,
		wasm.ModuleName,
		ibchookertypes.ModuleName,
		interchainquerytypes.ModuleName,
		liquidstakeibctypes.ModuleName,
		oracletypes.ModuleName,
		lscosmostypes.ModuleName,
	}
}

func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		ibcfeetypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		group.ModuleName,
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
		consensusparamtypes.ModuleName,
		halving.ModuleName,
		ibchookstypes.ModuleName,
		packetforwardtypes.ModuleName,
		wasm.ModuleName,
		epochstypes.ModuleName,
		ibchookertypes.ModuleName,
		interchainquerytypes.ModuleName,
		liquidstakeibctypes.ModuleName,
		oracletypes.ModuleName,
		lscosmostypes.ModuleName,
		buildertypes.ModuleName,
	}
}

// orderInitGenesis returns the order in which genesis is initialzed for modules
// NOTE: The genutils module must occur after staking so that pools are
// properly initialized with tokens from genesis accounts.
// NOTE: Capability module must occur first so that it can initialize any capabilities
// so that other modules that want to create or claim capabilities afterwards in InitChain
// can do so safely.
func orderInitGenesis() []string {
	return []string{
		capabilitytypes.ModuleName,
		banktypes.ModuleName,
		distributiontypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		ibcfeetypes.ModuleName,
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		group.ModuleName,
		authtypes.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		halving.ModuleName,
		ibchookstypes.ModuleName,
		packetforwardtypes.ModuleName,
		wasm.ModuleName,
		epochstypes.ModuleName,
		ibchookertypes.ModuleName,
		interchainquerytypes.ModuleName,
		liquidstakeibctypes.ModuleName,
		oracletypes.ModuleName,
		lscosmostypes.ModuleName,
		buildertypes.ModuleName,
	}
}
