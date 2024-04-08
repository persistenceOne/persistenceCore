package v11_8_0

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	liquidstaketypes "github.com/persistenceOne/pstake-native/v2/x/liquidstake/types"
	ratesynctypes "github.com/persistenceOne/pstake-native/v2/x/ratesync/types"

	"github.com/persistenceOne/persistenceCore/v11/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v11" //ONLY FOR MAINNET
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{liquidstaketypes.StoreKey, ratesynctypes.StoreKey},
	},
}
