package asset

import (
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/queries/asset"
	"strings"

	"github.com/gorilla/mux"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/constants"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/mint"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func RegisterRESTRoutes(cliContext context.CLIContext, router *mux.Router) {
	router.HandleFunc(strings.Join([]string{"", constants.TransactionRoute, constants.MintTransaction}, "/"), mint.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", constants.QuerierRoute, constants.AssetQuery, "{address}"}, "/"), asset.RestQueryHandler(cliContext)).Methods("GET")
}
