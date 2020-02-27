package reputation

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/persistenceOne/persistenceSDK/modules/reputation/constants"
	"github.com/persistenceOne/persistenceSDK/modules/reputation/transactions/feedback"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func RegisterRESTRoutes(cliContext context.CLIContext, router *mux.Router) {
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.FeedbackTransaction}, "/"), feedback.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.FeedbackTransaction}, "/"), QueryRequestHandler(cliContext)).Methods("GET")
}

func QueryRequestHandler(cliContext context.CLIContext) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")

		rest.PostProcessResponse(responseWriter, cliContext, "dfsdfsdfsdfs")

	}
}
