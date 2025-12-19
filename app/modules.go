package app

import (
	"cosmossdk.io/x/evidence"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
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
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/epochs"
	epochstypes "github.com/cosmos/cosmos-sdk/x/epochs/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/cosmos/cosmos-sdk/x/mint"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/protocolpool"
	protocolpooltypes "github.com/cosmos/cosmos-sdk/x/protocolpool/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/cosmos/gaia/v25/x/liquid"
	liquidtypes "github.com/cosmos/gaia/v25/x/liquid/types"
	packetforward "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v10/packetforward/types"
	ibchooks "github.com/cosmos/ibc-apps/modules/ibc-hooks/v10"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v10/types"
	ica "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts"
	icatypes "github.com/cosmos/ibc-go/v10/modules/apps/27-interchain-accounts/types"
	ibctransfer "github.com/cosmos/ibc-go/v10/modules/apps/transfer"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v10/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v10/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v10/modules/light-clients/07-tendermint"
	"github.com/persistenceOne/persistenceCore/v17/x/halving"
	"github.com/persistenceOne/persistenceCore/v17/x/liquidstake"
	liquidstaketypes "github.com/persistenceOne/persistenceCore/v17/x/liquidstake/types"
)

var ModuleAccountPermissions = map[string][]string{
	authtypes.FeeCollectorName:                  nil,
	distributiontypes.ModuleName:                nil,
	icatypes.ModuleName:                         nil,
	minttypes.ModuleName:                        {authtypes.Minter},
	stakingtypes.BondedPoolName:                 {authtypes.Burner, authtypes.Staking},
	stakingtypes.NotBondedPoolName:              {authtypes.Burner, authtypes.Staking},
	govtypes.ModuleName:                         {authtypes.Burner},
	ibctransfertypes.ModuleName:                 {authtypes.Minter, authtypes.Burner},
	wasmtypes.ModuleName:                        {authtypes.Burner},
	liquidstaketypes.ModuleName:                 {authtypes.Minter, authtypes.Burner},
	protocolpooltypes.ModuleName:                nil,
	protocolpooltypes.ProtocolPoolEscrowAccount: nil,
}

var receiveAllowedMAcc = map[string]bool{
	liquidstaketypes.ModuleName: true,
}

func appModules(
	app *Application,
	appCodec codec.Codec,
	txCfg client.TxConfig,
) []module.AppModule {

	return []module.AppModule{
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app, txCfg,
		),
		auth.NewAppModule(appCodec, *app.AccountKeeper, authsimulation.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		vesting.NewAppModule(*app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, *app.BankKeeper, app.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, *app.MintKeeper, app.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)), // nil -> SDK's default inflation function.
		slashing.NewAppModule(appCodec, *app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		distribution.NewAppModule(appCodec, *app.DistributionKeeper, app.AccountKeeper, app.BankKeeper, *app.StakingKeeper, app.GetSubspace(distributiontypes.ModuleName)),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.UpgradeKeeper, app.AccountKeeper.AddressCodec()),
		evidence.NewAppModule(*app.EvidenceKeeper),
		feegrantmodule.NewAppModule(appCodec, app.AccountKeeper, app.BankKeeper, *app.FeegrantKeeper, app.interfaceRegistry),
		authzmodule.NewAppModule(appCodec, *app.AuthzKeeper, app.AccountKeeper, app.BankKeeper, app.interfaceRegistry),
		ibc.NewAppModule(app.IBCKeeper),
		packetforward.NewAppModule(app.PacketForwardKeeper, app.GetSubspace(packetforwardtypes.ModuleName)),
		params.NewAppModule(*app.ParamsKeeper),
		halving.NewAppModule(appCodec, *app.HalvingKeeper),
		ibctransfer.NewAppModule(*app.TransferKeeper),
		ibchooks.NewAppModule(*app.AccountKeeper),
		ica.NewAppModule(app.ICAControllerKeeper, app.ICAHostKeeper),
		wasm.NewAppModule(appCodec, app.WasmKeeper, app.StakingKeeper, app.AccountKeeper, app.BankKeeper, app.MsgServiceRouter(), app.GetSubspace(wasm.ModuleName)),
		epochs.NewAppModule(*app.EpochsKeeper),
		liquid.NewAppModule(appCodec, app.LiquidKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		liquidstake.NewAppModule(*app.LiquidStakeKeeper),
		crisis.NewAppModule(app.CrisisKeeper, false, app.GetSubspace(crisistypes.ModuleName)), // skipGenesisInvariants: false, always be last to make sure that it checks for all invariants and not only part of them
		ibctm.NewAppModule(app.TMLightClientModule),
		protocolpool.NewAppModule(*app.ProtocolPoolKeeper, app.AccountKeeper, app.BankKeeper),
	}
}

func overrideSimulationModules(
	app *Application,
	appCodec codec.Codec,
) map[string]module.AppModuleSimulation {
	return map[string]module.AppModuleSimulation{
		authtypes.ModuleName: auth.NewAppModule(appCodec, *app.AccountKeeper, authsimulation.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
	}
}

func orderPreBlockers() []string {
	return []string{
		upgradetypes.ModuleName,
		authtypes.ModuleName,
	}
}

func orderBeginBlockers() []string {
	return []string{
		upgradetypes.ModuleName,
		epochstypes.ModuleName,
		minttypes.ModuleName,
		distributiontypes.ModuleName,
		protocolpooltypes.ModuleName,
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		slashingtypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		consensusparamtypes.ModuleName,
		halving.ModuleName,
		ibchookstypes.ModuleName,
		packetforwardtypes.ModuleName,
		wasmtypes.ModuleName,
		liquidtypes.ModuleName,
		liquidstaketypes.ModuleName,
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
		feegrant.ModuleName,
		authz.ModuleName,
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
		wasmtypes.ModuleName,
		epochstypes.ModuleName,
		liquidtypes.ModuleName,
		liquidstaketypes.ModuleName,
		protocolpooltypes.ModuleName,
	}
}

// orderInitGenesis returns the order in which genesis is initialized for modules
// NOTE: The genutils module must occur after staking so that pools are
// properly initialized with tokens from genesis accounts.
// NOTE: Capability module must occur first so that it can initialize any capabilities
// so that other modules that want to create or claim capabilities afterwards in InitChain
// can do so safely.
func orderInitGenesis() []string {
	return []string{
		authtypes.ModuleName,
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
		evidencetypes.ModuleName,
		feegrant.ModuleName,
		authz.ModuleName,
		genutiltypes.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		halving.ModuleName,
		ibchookstypes.ModuleName,
		packetforwardtypes.ModuleName,
		wasmtypes.ModuleName,
		epochstypes.ModuleName,
		liquidtypes.ModuleName,
		liquidstaketypes.ModuleName,
		protocolpooltypes.ModuleName,
	}
}

func orderExportGenesis() []string {
	return []string{
		consensusparamtypes.ModuleName,
		authtypes.ModuleName,
		protocolpooltypes.ModuleName, // Must be exported before bank
		banktypes.ModuleName,
		distributiontypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		halving.ModuleName,
		ibchookstypes.ModuleName,
		packetforwardtypes.ModuleName,
		epochstypes.ModuleName,
		liquidtypes.ModuleName,
		liquidstaketypes.ModuleName,
		wasmtypes.ModuleName,
	}
}
