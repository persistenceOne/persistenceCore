package app

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"
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
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributionclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts"
	icatypes "github.com/cosmos/ibc-go/v4/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v4/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v4/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v4/modules/core/02-client/client"
	ibchost "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	"github.com/persistenceOne/persistence-sdk/v2/x/epochs"
	epochstypes "github.com/persistenceOne/persistence-sdk/v2/x/epochs/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/halving"
	"github.com/persistenceOne/persistence-sdk/v2/x/ibchooker"
	ibchookertypes "github.com/persistenceOne/persistence-sdk/v2/x/ibchooker/types"
	"github.com/persistenceOne/persistence-sdk/v2/x/interchainquery"
	interchainquerytypes "github.com/persistenceOne/persistence-sdk/v2/x/interchainquery/types"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos"
	lscosmosclient "github.com/persistenceOne/pstake-native/v2/x/lscosmos/client"
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
	wasm.ModuleName:                          {authtypes.Burner},
	lscosmostypes.ModuleName:                 {authtypes.Minter, authtypes.Burner},
	lscosmostypes.DepositModuleAccount:       nil,
	lscosmostypes.DelegationModuleAccount:    nil,
	lscosmostypes.RewardModuleAccount:        nil,
	lscosmostypes.UndelegationModuleAccount:  nil,
	lscosmostypes.RewardBoosterModuleAccount: nil,
}

// ModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.NewAppModuleBasic(
		append(
			wasmclient.ProposalHandlers,
			paramsclient.ProposalHandler,
			distributionclient.ProposalHandler,
			upgradeclient.ProposalHandler,
			upgradeclient.CancelProposalHandler,
			ibcclient.UpdateClientProposalHandler,
			ibcclient.UpgradeProposalHandler,
			lscosmosclient.RegisterHostChainProposalHandler,
			lscosmosclient.MinDepositAndFeeChangeProposalHandler,
			lscosmosclient.PstakeFeeAddressChangeProposalHandler,
			lscosmosclient.AllowListValidatorSetChangeProposalHandler,
		)...,
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	transfer.AppModuleBasic{},
	vesting.AppModuleBasic{},
	wasm.AppModuleBasic{},
	halving.AppModuleBasic{},
	ica.AppModuleBasic{},
	epochs.AppModuleBasic{},
	interchainquery.AppModuleBasic{},
	ibchooker.AppModuleBasic{},
	lscosmos.AppModuleBasic{},
)

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
		mint.NewAppModule(appCodec, *app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, *app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper),
		distribution.NewAppModule(appCodec, *app.DistributionKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper),
		staking.NewAppModule(appCodec, *app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(*app.UpgradeKeeper),
		evidence.NewAppModule(*app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, *app.FeegrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, *app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
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
		mint.NewAppModule(appCodec, *app.MintKeeper, app.AccountKeeper),
		staking.NewAppModule(appCodec, *app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		distribution.NewAppModule(appCodec, *app.DistributionKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		slashing.NewAppModule(appCodec, *app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		params.NewAppModule(*app.ParamsKeeper),
		halving.NewAppModule(appCodec, *app.HalvingKeeper),
		authzmodule.NewAppModule(appCodec, *app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, *app.FeegrantKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
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
