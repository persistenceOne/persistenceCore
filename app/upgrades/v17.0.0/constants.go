package v17_0_0

import (
	store "cosmossdk.io/store/types"
	"github.com/persistenceOne/persistenceCore/v17/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v17.0.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{},
	},
}
