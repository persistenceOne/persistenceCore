package asset

import (
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/queries/asset"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/burn"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/lock"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/send"
	"strings"

	"github.com/gorilla/mux"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/constants"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/mint"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func RegisterRESTRoutes(cliContext context.CLIContext, router *mux.Router) {
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, constants.MintTransaction}, "/"), burn.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, constants.MintTransaction}, "/"), lock.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, constants.MintTransaction}, "/"), mint.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, constants.SendTransaction}, "/"), send.RestRequestHandler(cliContext)).Methods("POST")

	router.HandleFunc(strings.Join([]string{"", QuerierRoute, constants.AssetQuery, "{address}"}, "/"), asset.RestQueryHandler(cliContext)).Methods("GET")
}
