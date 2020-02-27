package share

import (
	"github.com/persistenceOne/persistenceSDK/modules/share/queries/share"
	"strings"

	"github.com/gorilla/mux"
	"github.com/persistenceOne/persistenceSDK/modules/share/constants"
	"github.com/persistenceOne/persistenceSDK/modules/share/transactions/burn"
	"github.com/persistenceOne/persistenceSDK/modules/share/transactions/lock"
	"github.com/persistenceOne/persistenceSDK/modules/share/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/modules/share/transactions/send"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func RegisterRESTRoutes(cliContext context.CLIContext, router *mux.Router) {
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, constants.MintTransaction}, "/"), burn.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, constants.MintTransaction}, "/"), lock.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, constants.MintTransaction}, "/"), mint.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", TransactionRoute, constants.SendTransaction}, "/"), send.RestRequestHandler(cliContext)).Methods("POST")

	router.HandleFunc(strings.Join([]string{"", QuerierRoute, constants.ShareQuery, "{address}"}, "/"), share.RestQueryHandler(cliContext)).Methods("GET")
}
