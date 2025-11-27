package liquidstake

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/core/appmodule"
	abci "github.com/cometbft/cometbft/abci/types"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/client/cli"
	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/keeper"
	"github.com/persistenceOne/persistenceCore/v16/x/liquidstake/types"
)

var (
	_ appmodule.AppModule       = AppModule{}
	_ module.AppModuleBasic     = AppModuleBasic{}
	_ appmodule.HasBeginBlocker = AppModule{}
)

// AppModuleBasic defines the basic application module used by the liquidstake module.
type AppModuleBasic struct{}

// Name returns the liquidstake module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec registers the liquidstake module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the liquidstake
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the liquidstake module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ sdkclient.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	if err := cdc.UnmarshalJSON(bz, &data); err != nil {
		return fmt.Errorf("failed to unmarshal %s genesis state: %w", types.ModuleName, err)
	}

	return types.ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the liquidstake module.
func (AppModuleBasic) RegisterRESTRoutes(_ sdkclient.Context, _ *mux.Router) {}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the liquidstake module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx sdkclient.Context, mux *runtime.ServeMux) {
	if err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx)); err != nil {
		panic(err)
	}
}

// GetTxCmd returns the root tx command for the liquidstake module.
func (AppModuleBasic) GetTxCmd() *cobra.Command { return cli.GetTxCmd() }

func (b AppModuleBasic) RegisterInterfaces(reg cdctypes.InterfaceRegistry) {
	types.RegisterInterfaces(reg)
}

// AppModule implements an application module for the liquidstake module.
type AppModule struct {
	AppModuleBasic

	keeper keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(
	keeper keeper.Keeper,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		keeper:         keeper,
	}
}

func (am AppModule) IsOnePerModuleType() {}

func (am AppModule) IsAppModule() {}

// Name returns the liquidstake module's name.
func (AppModule) Name() string {
	return types.ModuleName
}

// RegisterInvariants registers the liquidstake module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// QuerierRoute returns the liquidstake module's querier route name.
func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.Querier{Keeper: am.keeper})
}

// InitGenesis performs genesis initialization for the liquidstake module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONCodec, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	am.keeper.InitGenesis(ctx, genesisState)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the liquidstake
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONCodec) json.RawMessage {
	gs := am.keeper.ExportGenesis(ctx)
	return cdc.MustMarshalJSON(gs)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }

// BeginBlock returns the begin blocker for the liquidstake module.
func (am AppModule) BeginBlock(ctx context.Context) error {
	return BeginBlock(ctx, am.keeper)
}
