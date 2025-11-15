package v16_0_0_rc0

import (
	store "cosmossdk.io/store/types"
	protocolpooltypes "github.com/cosmos/cosmos-sdk/x/protocolpool/types"

	"github.com/persistenceOne/persistenceCore/v16/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v16.0.0-rc0"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added:   []string{protocolpooltypes.StoreKey},
		Deleted: []string{},
	},
}
