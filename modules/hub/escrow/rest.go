package escrow

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/persistenceOne/persistenceSDK/modules/hub/escrow/constants"
	"github.com/persistenceOne/persistenceSDK/modules/hub/escrow/transactions/execute"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterRESTRoutes(cliContext context.CLIContext, router *mux.Router) {
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.ExecuteTransaction}, "/"), execute.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.ExecuteTransaction}, "/"), QueryRequestHandler(cliContext)).Methods("GET")
}

func QueryRequestHandler(cliContext context.CLIContext) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")

		rest.PostProcessResponse(responseWriter, cliContext, "dfsdfsdfsdfs")

	}
}
