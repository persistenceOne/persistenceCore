package asset

import (
	"encoding/json"

	abciTypes "github.com/tendermint/tendermint/abci/types"

	sdkTypes "github.com/cosmos/cosmos-sdk/types"

	"github.com/persistenceOne/persistenceSDK/modules/asset/constants"

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
func (AppModuleBasic) RegisterCodec(codec *codec.Codec) {
	RegisterCodec(codec)
}
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return packageCodec.MustMarshalJSON(DefaultGenesisState())
}
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var genesisState GenesisState
	error := packageCodec.UnmarshalJSON(bz, &genesisState)
	if error != nil {
		return error
	}
	return ValidateGenesis(genesisState)
}
func (AppModuleBasic) RegisterRESTRoutes(cliContext context.CLIContext, router *mux.Router) {
	RegisterRESTRoutes(cliContext, router)
}
func (AppModuleBasic) GetTxCmd(codec *codec.Codec) *cobra.Command {
	return GetCLIRootTransactionCommand(codec)
}
func (AppModuleBasic) GetQueryCmd(codec *codec.Codec) *cobra.Command {
	return GetCLIRootQueryCommand(codec)
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
	packageCodec.MustUnmarshalJSON(data, &genesisState)
	InitializeGenesisState(context, appModule.keeper, genesisState)
	return []abciTypes.ValidatorUpdate{}
}
func (appModule AppModule) ExportGenesis(context sdkTypes.Context) json.RawMessage {
	gs := ExportGenesis(context, appModule.keeper)
	return packageCodec.MustMarshalJSON(gs)
}
func (AppModule) BeginBlock(_ sdkTypes.Context, _ abciTypes.RequestBeginBlock) {}

func (AppModule) EndBlock(_ sdkTypes.Context, _ abciTypes.RequestEndBlock) []abciTypes.ValidatorUpdate {
	return []abciTypes.ValidatorUpdate{}
}
