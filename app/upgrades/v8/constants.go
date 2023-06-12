package v8

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	consensusparamstypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	oracletypes "github.com/persistenceOne/persistence-sdk/v2/x/oracle/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v2/x/liquidstakeibc/types"
	lspersistencetypes "github.com/persistenceOne/pstake-native/v2/x/lspersistence/types"

	"github.com/persistenceOne/persistenceCore/v8/app/upgrades"
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
			oracletypes.ModuleName,
			crisistypes.ModuleName,
			liquidstakeibctypes.ModuleName,
			lspersistencetypes.ModuleName,
			consensusparamstypes.ModuleName,
		},
	},
}
