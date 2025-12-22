package v17_0_0_rc0

import (
	store "cosmossdk.io/store/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/persistenceOne/persistenceCore/v17/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v17.0.0-rc0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{},
		Deleted: []string{paramstypes.StoreKey, crisistypes.StoreKey},
	},
}
