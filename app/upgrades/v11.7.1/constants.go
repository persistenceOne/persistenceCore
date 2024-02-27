package v11_7_1

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	liquidstaketypes "github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
	ratesynctypes "github.com/persistenceOne/pstake-native/v2/x/ratesync/types"

	"github.com/persistenceOne/persistenceCore/v11/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v11.7.1"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{liquidstaketypes.StoreKey, ratesynctypes.StoreKey},
	},
}
