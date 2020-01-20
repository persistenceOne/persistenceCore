package sign

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func TransactionCommand(cdc *codec.Codec) *cobra.Command {
	const (
		ContractFlag = "contract"
	)
	command := &cobra.Command{
		Use:   "sign",
		Short: "Create and sign transaction to sign at contract",
		Long:  "",
		RunE: func(command *cobra.Command, args []string) error {
			transactionBuilder := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliContext := context.NewCLIContext().WithCodec(cdc)

			message := Message{
				From: cliContext.GetFromAddress(),
			}

			if error := message.ValidateBasic(); error != nil {
				return error
			}

			return utils.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdkTypes.Msg{message})
		},
	}

	command.Flags().String(ContractFlag, "", "Contract")
	return command
}
