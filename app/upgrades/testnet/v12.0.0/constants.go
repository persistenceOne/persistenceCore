package v11_8_0

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/persistenceOne/persistenceCore/v11/app/upgrades"
	alliancetypes "github.com/terra-money/alliance/x/alliance/types"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v12" //ONLY FOR MAINNET
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{alliancetypes.StoreKey},
	},
}
