package mint

import (
	"bufio"
	"github.com/persistenceOne/persistenceSDK/modules/share/constants"
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
		Use:   constants.MintTransaction,
		Short: "Create and sign transaction to mint an share",
		Long:  "",
		RunE: func(command *cobra.Command, args []string) error {
			bufioReader := bufio.NewReader(command.InOrStdin())
			transactionBuilder := auth.NewTxBuilderFromCLI(bufioReader).WithTxEncoder(auth.DefaultTxEncoder(codec))
			cliContext := context.NewCLIContextWithInput(bufioReader).WithCodec(codec)
			to, err := sdkTypes.AccAddressFromBech32(viper.GetString(constants.ToFlag))
			if err != nil {
				return err
			}
			message := Message{
				From:    cliContext.GetFromAddress(),
				To:      to,
				Address: viper.GetString(constants.AddressFlag),
				Lock:    viper.GetBool(constants.LockFlag),
			}

			if err := message.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdkTypes.Msg{message})
		},
	}

	command.Flags().String(constants.AddressFlag, "", "address")
	command.Flags().String(constants.ToFlag, "", "to")
	command.Flags().Bool(constants.LockFlag, false, "lock")
	return command
}
