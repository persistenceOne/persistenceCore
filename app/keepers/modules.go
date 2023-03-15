package keepers

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmclient "github.com/CosmWasm/wasmd/x/wasm/client"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	feegrantmodule "github.com/cosmos/cosmos-sdk/x/feegrant/module"
	"github.com/cosmos/cosmos-sdk/x/gov"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/cosmos/cosmos-sdk/x/upgrade/client"
	ica "github.com/cosmos/ibc-go/v6/modules/apps/27-interchain-accounts"
	ibcfee "github.com/cosmos/ibc-go/v6/modules/apps/29-fee"
	"github.com/cosmos/ibc-go/v6/modules/apps/transfer"
	ibc "github.com/cosmos/ibc-go/v6/modules/core"
	ibcclient "github.com/cosmos/ibc-go/v6/modules/core/02-client/client"
	"github.com/persistenceOne/persistence-sdk/v2/x/epochs"
	"github.com/persistenceOne/persistence-sdk/v2/x/halving"
	"github.com/persistenceOne/persistence-sdk/v2/x/ibchooker"
	"github.com/persistenceOne/persistence-sdk/v2/x/interchainquery"
	"github.com/persistenceOne/persistence-sdk/v2/x/lsnative/distribution"
	distributionclient "github.com/persistenceOne/persistence-sdk/v2/x/lsnative/distribution/client"
	"github.com/persistenceOne/persistence-sdk/v2/x/lsnative/genutil"
	"github.com/persistenceOne/persistence-sdk/v2/x/lsnative/slashing"
	"github.com/persistenceOne/persistence-sdk/v2/x/lsnative/staking"
	"github.com/persistenceOne/persistence-sdk/v2/x/oracle"
	"github.com/persistenceOne/pstake-native/v2/x/lscosmos"
	lscosmosclient "github.com/persistenceOne/pstake-native/v2/x/lscosmos/client"
)

// AppModuleBasics defines the module BasicManager is in charge of setting up basic,
// non-dependant module elements, such as codec registration
// and genesis verification.
var AppModuleBasics = []module.AppModuleBasic{
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
			upgradeclient.LegacyProposalHandler,
			upgradeclient.LegacyCancelProposalHandler,
			ibcclient.UpdateClientProposalHandler,
			ibcclient.UpgradeProposalHandler,
			lscosmosclient.MinDepositAndFeeChangeProposalHandler,
			lscosmosclient.PstakeFeeAddressChangeProposalHandler,
			lscosmosclient.AllowListValidatorSetChangeProposalHandler,
		),
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
	ibcfee.AppModuleBasic{},
	oracle.AppModuleBasic{},
}
