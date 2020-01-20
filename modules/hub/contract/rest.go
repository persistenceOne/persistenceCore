package asset

import (
	"strings"

	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/constants"
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/transactions/mint"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func RegisterRESTRoutes(cliContext context.CLIContext, router *mux.Router) {
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.MintTransaction}, "/"), mint.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.MintTransaction}, "/"), QueryRequestHandler(cliContext)).Methods("GET")
}

func QueryRequestHandler(cliContext context.CLIContext) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")

		rest.PostProcessResponse(responseWriter, cliContext, "dfsdfsdfsdfs")

	}
}
