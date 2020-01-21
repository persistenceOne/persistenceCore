package burn

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

type Request struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`
	Asset   string       `json:"asset" yaml:"asset"`
}

func RestRequestHandler(cliContext context.CLIContext) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		var request Request
		if !rest.ReadRESTReq(responseWriter, httpRequest, cliContext.Codec, &request) {
			return
		}

		request.BaseReq = request.BaseReq.Sanitize()
		if !request.BaseReq.ValidateBasic(responseWriter) {
			return
		}

		from, error := sdkTypes.AccAddressFromBech32(request.BaseReq.From)
		if error != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, error.Error())
			return
		}

		message := Message{
			From: from,
		}
		utils.WriteGenerateStdTxResponse(responseWriter, cliContext, request.BaseReq, []sdkTypes.Msg{message})
	}
}
