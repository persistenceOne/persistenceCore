package v16_0_0

import (
	store "cosmossdk.io/store/types"
	"github.com/persistenceOne/persistenceCore/v16/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v16.0.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{}, //TODO add protocolpool
		Deleted: []string{},
	},
}
