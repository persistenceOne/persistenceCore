package v8

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	consensusparamstypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	group "github.com/cosmos/cosmos-sdk/x/group"
	packetforwardtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
	ibchookstypes "github.com/cosmos/ibc-apps/modules/ibc-hooks/v7/types"
	ibcfeetypes "github.com/cosmos/ibc-go/v7/modules/apps/29-fee/types"
	oracletypes "github.com/persistenceOne/persistence-sdk/v3/x/oracle/types"
	liquidstakeibctypes "github.com/persistenceOne/pstake-native/v3/x/liquidstakeibc/types"
	buildertypes "github.com/skip-mev/pob/x/builder/types"

	"github.com/persistenceOne/persistenceCore/v12/app/upgrades"
)

const (
	// UpgradeName defines the on-chain upgrade name.
	UpgradeName = "v8"

	// BondDenom defines current active bond denom for testnet/mainnet.
	BondDenom = "uxprt"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			ibcfeetypes.ModuleName,
			oracletypes.ModuleName,
			crisistypes.ModuleName,
			group.ModuleName,
			liquidstakeibctypes.ModuleName,
			consensusparamstypes.ModuleName,
			ibchookstypes.StoreKey,
			packetforwardtypes.ModuleName,
			buildertypes.ModuleName,
		},
	},
}
