package asset

import (
	"strings"

	"github.com/commitHub/commitBlockchain/modules/hub/asset/constants"
	"github.com/commitHub/commitBlockchain/modules/hub/asset/transactions/mint"
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func RegisterRESTRoutes(cliContext context.CLIContext, router *mux.Router) {
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.MintTransaction}, "/"), mint.RestRequestHandler(cliContext)).Methods("POST")
	router.HandleFunc(strings.Join([]string{"", constants.ModuleName, constants.MintTransaction}, "/"), mint.QueryRequestHandler(cliContext)).Methods("GET")
}
