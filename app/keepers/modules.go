package keepers

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/consensus"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	ica "github.com/cosmos/ibc-go/v7/modules/apps/27-interchain-accounts"
	ibcfee "github.com/cosmos/ibc-go/v7/modules/apps/29-fee"
	"github.com/cosmos/ibc-go/v7/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v7/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v7/modules/core/02-client/client"
	ibctm "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"
	"github.com/persistenceOne/persistence-sdk/v2/x/epochs"
	"github.com/persistenceOne/persistence-sdk/v2/x/halving"
	"github.com/persistenceOne/persistence-sdk/v2/x/ibchooker"
	"github.com/persistenceOne/persistence-sdk/v2/x/interchainquery"
	"github.com/persistenceOne/persistence-sdk/v2/x/oracle"
	"github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos"
	lscosmosclient "github.com/persistenceOne/pstake-native/v2/x/lscosmos/client"
	"github.com/persistenceOne/pstake-native/v2/x/lspersistence"
)

// AppModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var AppModuleBasics = []module.AppModuleBasic{
	auth.AppModuleBasic{},
	genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	gov.NewAppModuleBasic(
		[]govclient.ProposalHandler{
			paramsclient.ProposalHandler,
			upgradeclient.LegacyProposalHandler,
			upgradeclient.LegacyCancelProposalHandler,
			ibcclient.UpdateClientProposalHandler,
			ibcclient.UpgradeProposalHandler,
			lscosmosclient.MinDepositAndFeeChangeProposalHandler,
			lscosmosclient.PstakeFeeAddressChangeProposalHandler,
			lscosmosclient.AllowListValidatorSetChangeProposalHandler,
		},
	),
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	feegrantmodule.AppModuleBasic{},
	authzmodule.AppModuleBasic{},
	ibc.AppModuleBasic{},
	ibctm.AppModuleBasic{},
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
	ibcfee.AppModuleBasic{},
	oracle.AppModuleBasic{},
	liquidstakeibc.AppModuleBasic{},
	lspersistence.AppModuleBasic{},
	consensus.AppModuleBasic{},
	groupmodule.AppModuleBasic{},
}
