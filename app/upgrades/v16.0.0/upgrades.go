package v16_0_0

import (
	"context"
	"fmt"

	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"

	"github.com/persistenceOne/persistenceCore/v16/app/keepers"
	"github.com/persistenceOne/persistenceCore/v16/app/upgrades"
)

func CreateUpgradeHandler(args upgrades.UpgradeHandlerArgs) upgradetypes.UpgradeHandler {
	return func(ctx context.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		sdkCtx := sdk.UnwrapSDKContext(ctx)
		sdkCtx.Logger().Info("Running upgrade handler")

		sdkCtx.Logger().Info("running module migrations...")
		vm, err := args.ModuleManager.RunMigrations(sdkCtx, args.Configurator, vm)
		if err != nil {
			return vm, err
		}

		err = ResetIBCTransferVersions(sdkCtx, args.Keepers)
		if err != nil {
			return vm, err
		}

		sdkCtx.Logger().Info("Upgrade complete")
		return vm, nil
	}
}

func ResetIBCTransferVersions(sdkCtx sdk.Context, keepers *keepers.AppKeepers) error {
	channels := keepers.IBCKeeper.ChannelKeeper.GetAllChannels(sdkCtx)
	for _, channel := range channels {
		if channel.PortId == ibctransfertypes.PortID && channel.Version != ibctransfertypes.V1 {
			channelDb, ok := keepers.IBCKeeper.ChannelKeeper.GetChannel(sdkCtx, channel.PortId, channel.ChannelId)
			if !ok {
				return fmt.Errorf("channel %s not found", channel.ChannelId)
			}
			channelDb.Version = ibctransfertypes.V1
			keepers.IBCKeeper.ChannelKeeper.SetChannel(sdkCtx, channel.PortId, channel.ChannelId, channelDb)
		}
	}

	return nil
}
