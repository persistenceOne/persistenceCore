package v9

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	lscosmostypes "github.com/persistenceOne/pstake-native/v2/x/lscosmos/types"

	"github.com/persistenceOne/persistenceCore/v11/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v9"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{lscosmostypes.StoreKey},
	},
}
