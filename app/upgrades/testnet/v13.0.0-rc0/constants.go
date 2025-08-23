package v13_0_0_rc0

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/persistenceOne/persistenceCore/v13/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v13.0.0-rc0"
)

const capabilityStoreKey = "capability"

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Deleted: []string{capabilityStoreKey},
	},
}
