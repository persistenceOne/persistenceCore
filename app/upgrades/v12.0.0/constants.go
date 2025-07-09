package v12_0_0

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/cosmos/cosmos-sdk/x/group"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	interchainquerytypes "github.com/persistenceOne/persistence-sdk/v3/x/interchainquery/types"
	oracletypes "github.com/persistenceOne/persistence-sdk/v3/x/oracle/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v3/x/liquidstakeibc/types"
	ratesynctypes "github.com/persistenceOne/pstake-native/v3/x/ratesync/types"
	buildertypes "github.com/skip-mev/pob/x/builder/types"

	"github.com/persistenceOne/persistenceCore/v12/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v12.0.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Deleted: []string{ibcfeetypes.StoreKey, group.StoreKey,
			interchainquerytypes.StoreKey, liquidstakeibctypes.StoreKey, ratesynctypes.StoreKey,
			oracletypes.StoreKey, buildertypes.StoreKey,
		},
	},
}
