package share

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/constants"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/burn"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/lock"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/modules/hub/share/transactions/send"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterRESTRoutes(cliContext context.CLIContext, router *mux.Router) {
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.BurnTransaction}, "/"), burn.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.LockTransaction}, "/"), lock.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.MintTransaction}, "/"), mint.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.SendTransaction}, "/"), send.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.ShareQuery}, "/"), QueryRequestHandler(cliContext)).Methods("GET")
}

func QueryRequestHandler(cliContext context.CLIContext) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		rest.PostProcessResponse(responseWriter, cliContext, constants.ShareQuery)

	}
}
