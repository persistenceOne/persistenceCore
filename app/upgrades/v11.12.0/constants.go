package v11_12_0

import (
	store "github.com/cosmos/cosmos-sdk/store/types"

	"github.com/persistenceOne/persistenceCore/v11/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v11.12.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades:        store.StoreUpgrades{},
}
