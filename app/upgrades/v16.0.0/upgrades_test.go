package v16_0_0_test

import (
	"testing"

	"cosmossdk.io/log"
	"github.com/CosmWasm/wasmd/x/wasm"
	dbm "github.com/cosmos/cosmos-db"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	ibctransfertypes "github.com/cosmos/ibc-go/v10/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v10/modules/core/04-channel/types"
	"github.com/persistenceOne/persistenceCore/v16/app/constants"
	v1600 "github.com/persistenceOne/persistenceCore/v16/app/upgrades/v16.0.0"
	"github.com/stretchr/testify/require"

	"github.com/persistenceOne/persistenceCore/v16/app"
)

func TestResetIBCTransferVersions(t *testing.T) {
	constants.SetConfig()
	testApp := app.NewApplication(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.NewAppOptionsWithFlagHome(t.TempDir()), []wasm.Option{})
	ctx := testApp.NewContext(true)

	type args struct {
		portID    string
		channelID string
		version   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "correct case",
			args: args{ibctransfertypes.PortID, "channel-0", ibctransfertypes.V1},
		},
		{
			name: "fee version should be converted to v1",
			args: args{ibctransfertypes.PortID, "channel-1", "{\"fee_version\":\"ics29-1\",\"app_version\":\"ics20-1\"}"},
		},
		{
			name: "random port ids should not be affected",
			args: args{"ica-controller-xxx", "channel-2", ibctransfertypes.V1},
		},
		{
			name: "random port ids should not be affected 2",
			args: args{"ica-controller-xxxyyy", "channel-3", "ibcfee-0,controller-1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testApp.IBCKeeper.ChannelKeeper.SetChannel(ctx, tt.args.portID, tt.args.channelID, channeltypes.Channel{
				Version: tt.args.version,
			})
			err := v1600.ResetIBCTransferVersions(ctx, testApp.AppKeepers)
			require.NoError(t, err)
			channel, ok := testApp.IBCKeeper.ChannelKeeper.GetChannel(ctx, tt.args.portID, tt.args.channelID)
			require.True(t, ok)
			if tt.args.portID == ibctransfertypes.PortID {
				require.Equal(t, ibctransfertypes.V1, channel.Version)
			} else {
				require.Equal(t, tt.args.version, channel.Version)

			}

		})
	}
}
