package v14_0_0

import (
	store "cosmossdk.io/store/types"
	liquidtypes "github.com/cosmos/gaia/v24/x/liquid/types"
	"github.com/persistenceOne/persistenceCore/v15/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v14.0.0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{liquidtypes.ModuleName},
		Deleted: []string{},
	},
}
