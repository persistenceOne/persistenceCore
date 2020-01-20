package contract

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/persistenceOne/persistenceSDK/modules/hub/contract/constants"
	"github.com/persistenceOne/persistenceSDK/modules/hub/contract/transactions/sign"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterRESTRoutes(cliContext context.CLIContext, router *mux.Router) {
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.SignTransaction}, "/"), sign.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.SignTransaction}, "/"), QueryRequestHandler(cliContext)).Methods("GET")
}

func QueryRequestHandler(cliContext context.CLIContext) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")

		rest.PostProcessResponse(responseWriter, cliContext, "dfsdfsdfsdfs")

	}
}
