package burn

import (
	"github.com/persistenceOne/persistenceSDK/modules/asset/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func TransactionCommand(codec *codec.Codec) *cobra.Command {

	command := &cobra.Command{
		Use:   constants.BurnTransaction,
		Short: "Create and sign transaction to burn an asset",
		Long:  "",
		RunE: func(command *cobra.Command, args []string) error {
			transactionBuilder := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(codec))
			cliContext := context.NewCLIContext().WithCodec(codec)

			message := Message{
				From:    cliContext.GetFromAddress(),
				Address: viper.GetString(constants.AddressFlag),
			}

			if err := message.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdkTypes.Msg{message})
		},
	}

	command.Flags().String(constants.AddressFlag, "", "address")
	return command
}
