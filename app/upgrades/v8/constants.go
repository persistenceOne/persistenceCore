package v8

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	consensusparamstypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	group "github.com/cosmos/cosmos-sdk/x/group"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	ibchookstypes "github.com/persistenceOne/persistence-sdk/v2/x/ibc-hooks/types"
	oracletypes "github.com/persistenceOne/persistence-sdk/v2/x/oracle/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	routertypes "github.com/strangelove-ventures/packet-forward-middleware/v7/router/types"

	"github.com/persistenceOne/persistenceCore/v8/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v8"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			ibcfeetypes.ModuleName,
			oracletypes.ModuleName,
			crisistypes.ModuleName,
			group.ModuleName,
			liquidstakeibctypes.ModuleName,
			consensusparamstypes.ModuleName,
			ibchookstypes.StoreKey,
			routertypes.ModuleName,
		},
	},
}
