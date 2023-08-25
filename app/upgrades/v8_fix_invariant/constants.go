package v8_fix_invariant

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/persistenceOne/persistenceCore/v8/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v8-fix-invariant"
	ValidatorAddress = "persistencevaloper1qxrztmfe2pevv8nam2znfe5p3eep3r4q9wn86p"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{},
	},
}
