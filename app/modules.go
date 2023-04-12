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
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts"
	icatypes "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts/types"
	ibcfee "github.com/cosmos/ibc-go/v6/modules/apps/29-fee"
	ibcfeetypes "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v6/modules/core"
	ibchost "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/persistenceOne/persistence-sdk/v2/x/epochs"
	epochstypes "github.com/persistenceOne/persistence-sdk/v2/x/epochs/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/halving"
	ibchookertypes "github.com/persistenceOne/persistence-sdk/v2/x/ibchooker/types"
	interchainquerytypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/lsnative/distribution"
	distributiontypes "github.com/persistenceOne/persistence-sdk/v2/x/lsnative/distribution/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/lsnative/genutil"
	genutiltypes "github.com/persistenceOne/persistence-sdk/v2/x/lsnative/genutil/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/lsnative/slashing"
	slashingtypes "github.com/persistenceOne/persistence-sdk/v2/x/lsnative/slashing/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/lsnative/staking"
	stakingtypes "github.com/persistenceOne/persistence-sdk/v2/x/lsnative/staking/types"
	lscosmostypes "github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"

	appparams "github.com/persistenceOne/persistenceCore/v7/app/params"
)

var ModuleAccountPermissions = map[string][]string{
	authtypes.FeeCollectorName:               nil,
	distributiontypes.ModuleName:             nil,
	icatypes.ModuleName:                      nil,
	minttypes.ModuleName:                     {authtypes.Minter},
	stakingtypes.BondedPoolName:              {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName:           {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:                      {authtypes.Burner},
	ibctransfertypes.ModuleName:              {authtypes.Minter, authtypes.Burner},
	ibcfeetypes.ModuleName:                   nil,
	wasm.ModuleName:                          {authtypes.Burner},
	lscosmostypes.ModuleName:                 {authtypes.Minter, authtypes.Burner},
	lscosmostypes.DepositModuleAccount:       nil,
	lscosmostypes.DelegationModuleAccount:    nil,
	lscosmostypes.RewardModuleAccount:        nil,
	lscosmostypes.UndelegationModuleAccount:  nil,
	lscosmostypes.RewardBoosterModuleAccount: nil,
}

func appModules(
	app *Application,
	encodingConfig appparams.EncodingConfig,
	skipGenesisInvariants bool,
) []module.AppModule {
	appCodec := encodingConfig.Marshaler

	return []module.AppModule{
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TransactionConfig,
		),
		auth.NewAppModule(appCodec, *app.AccountKeeper, nil),
		vesting.NewAppModule(*app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, *app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		gov.NewAppModule(appCodec, *app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, *app.MintKeeper, app.AccountKeeper, nil), // nil -> SDK's default inflation function.
		slashing.NewAppModule(appCodec, *app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper),
		distribution.NewAppModule(appCodec, *app.DistributionKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper),
		staking.NewAppModule(appCodec, *app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(*app.UpgradeKeeper),
		evidence.NewAppModule(*app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, *app.FeegrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, *app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		ibcfee.NewAppModule(*app.IBCFeeKeeper),
		params.NewAppModule(*app.ParamsKeeper),
		halving.NewAppModule(appCodec, *app.HalvingKeeper),
		app.TransferModule,
		app.IBCTransferHooksMiddleware,
		ica.NewAppModule(app.ICAControllerKeeper, app.ICAHostKeeper),
		wasm.NewAppModule(appCodec, app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		epochs.NewAppModule(*app.EpochsKeeper),
		app.InterchainQueryModule,
		app.LSCosmosModule,
		crisis.NewAppModule(app.CrisisKeeper, skipGenesisInvariants), // always be last to make sure that it checks for all invariants and not only part of them
	}
}

func simulationModules(
	app *Application,
	encodingConfig appparams.EncodingConfig,
	_ bool,
) []module.AppModuleSimulation {
	appCodec := encodingConfig.Marshaler

	return []module.AppModuleSimulation{
		auth.NewAppModule(appCodec, *app.AccountKeeper, authsimulation.RandomGenesisAccounts),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		gov.NewAppModule(appCodec, *app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, *app.MintKeeper, app.AccountKeeper, nil), // nil -> SDK's default inflation function.
		staking.NewAppModule(appCodec, *app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		distribution.NewAppModule(appCodec, *app.DistributionKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper),
		slashing.NewAppModule(appCodec, *app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper),
		params.NewAppModule(*app.ParamsKeeper),
		halving.NewAppModule(appCodec, *app.HalvingKeeper),
		authzmodule.NewAppModule(appCodec, *app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, *app.FeegrantKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		ibcfee.NewAppModule(*app.IBCFeeKeeper),
		app.TransferModule,
		app.InterchainQueryModule,
		app.LSCosmosModule,
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
		ibchost.ModuleName,
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
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		halving.ModuleName,
		wasm.ModuleName,
		ibchookertypes.ModuleName,
		interchainquerytypes.ModuleName,
		lscosmostypes.ModuleName,
	}
}

func orderEndBlockers() []string {
	return []string{
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibchost.ModuleName,
		icatypes.ModuleName,
		ibcfeetypes.ModuleName,
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
		ibchost.ModuleName,
		icatypes.ModuleName,
		ibcfeetypes.ModuleName,
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
	}
}
