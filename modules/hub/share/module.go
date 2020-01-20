package share

import (
	"encoding/json"

	abciTypes "github.com/tendermint/tendermint/abci/types"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/persistenceSDK/modules/hub/share/constants"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct {
}

func (AppModuleBasic) Name() string {
	return constants.ModuleName
}
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return cdc.MustMarshalJSON(DefaultGenesisState())
}
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var genesisState GenesisState
	error := cdc.UnmarshalJSON(bz, &genesisState)
	if error != nil {
		return error
	}
	return ValidateGenesis(genesisState)
}
func (AppModuleBasic) RegisterRESTRoutes(cliContext context.CLIContext, router *mux.Router) {
	RegisterRESTRoutes(cliContext, router)
}
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return GetCLIRootTransactionCommand(cdc)
}
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return GetCLIRootQueryCommand(cdc)
}

type AppModule struct {
	AppModuleBasic
	keeper Keeper
}

func NewAppModule(keeper Keeper) AppModule {
	return AppModule{keeper: keeper}
}
func (AppModule) Name() string {
	return ModuleName
}
func (appModule AppModule) RegisterInvariants(_ sdkTypes.InvariantRegistry) {}
func (AppModule) Route() string {
	return TransactionRoute
}
func (appModule AppModule) NewHandler() sdkTypes.Handler {
	return NewHandler(appModule.keeper)
}
func (AppModule) QuerierRoute() string {
	return QuerierRoute
}
func (appModule AppModule) NewQuerierHandler() sdkTypes.Querier {
	return NewQuerier(appModule.keeper)
}
func (appModule AppModule) InitGenesis(context sdkTypes.Context, data json.RawMessage) []abciTypes.ValidatorUpdate {
	var genesisState GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)
	InitializeGenesisState(context, appModule.keeper, genesisState)
	return []abciTypes.ValidatorUpdate{}
}
func (appModule AppModule) ExportGenesis(context sdkTypes.Context) json.RawMessage {
	gs := ExportGenesis(context, appModule.keeper)
	return cdc.MustMarshalJSON(gs)
}
func (AppModule) BeginBlock(_ sdkTypes.Context, _ abciTypes.RequestBeginBlock) {}

func (AppModule) EndBlock(_ sdkTypes.Context, _ abciTypes.RequestEndBlock) []abciTypes.ValidatorUpdate {
	return []abciTypes.ValidatorUpdate{}
}
