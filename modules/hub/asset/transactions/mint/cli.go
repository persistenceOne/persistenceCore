package mint

import (
	"github.com/persistenceOne/persistenceSDK/modules/hub/asset/constants"
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
		Use:   "mint",
		Short: "Create and sign transaction to mint at asset",
		Long:  "",
		RunE: func(command *cobra.Command, args []string) error {
			transactionBuilder := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(codec))
			cliContext := context.NewCLIContext().WithCodec(codec)
			to, err := sdkTypes.AccAddressFromBech32(viper.GetString(constants.ToFlag))
			if err != nil {
				return err
			}
			message := Message{
				From:    cliContext.GetFromAddress(),
				To:      to,
				Address: viper.GetString(constants.AddressFlag),
			}

			if err := message.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdkTypes.Msg{message})
		},
	}

	command.Flags().String(constants.AddressFlag, "", "address")
	command.Flags().String(constants.ToFlag, "", "to")
	return command
}
