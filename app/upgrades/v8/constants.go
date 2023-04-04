package v8

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v6/modules/apps/29-fee/types"

	"github.com/persistenceOne/persistenceCore/v7/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v8"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			ibcfeetypes.ModuleName,
		},
	},
}
