package lock

import (
	"github.com/asaskevich/govalidator"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

type Request struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req" valid:"required~base_req"`
	Address string       `json:"address" yaml:"address" valid:"required~address"`
	Lock    bool         `json:"lock" yaml:"lock"`
}

func RestRequestHandler(cliContext context.CLIContext) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		var request Request
		if !rest.ReadRESTReq(responseWriter, httpRequest, cliContext.Codec, &request) {
			return
		}

		request.BaseReq = request.BaseReq.Sanitize()
		if !request.BaseReq.ValidateBasic(responseWriter) {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, "")
			return
		}

		_, error := govalidator.ValidateStruct(request)
		if error != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, error.Error())
			return
		}

		from, error := sdkTypes.AccAddressFromBech32(request.BaseReq.From)
		if error != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, error.Error())
			return
		}

		message := Message{
			From:    from,
			Address: request.Address,
			Lock:    request.Lock,
		}
		utils.WriteGenerateStdTxResponse(responseWriter, cliContext, request.BaseReq, []sdkTypes.Msg{message})
	}
}
