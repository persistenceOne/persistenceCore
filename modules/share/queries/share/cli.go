package share

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/persistenceOne/persistenceSDK/modules/share/constants"
	"github.com/persistenceOne/persistenceSDK/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strings"
)

func QueryCommand(codec *codec.Codec) *cobra.Command {
	command := &cobra.Command{
		Use:   constants.ShareQuery,
		Short: "Query a share.",
		Long:  "",
		RunE: func(command *cobra.Command, args []string) error {
			cliContext := context.NewCLIContext().WithCodec(codec)

			bytes := packageCodec.MustMarshalJSON(query{
				Address: viper.GetString(constants.AddressFlag),
			})

			response, _, queryWithDataError := cliContext.QueryWithData(strings.Join([]string{"", "custom", constants.QuerierRoute, constants.ShareQuery}, "/"), bytes)
			if queryWithDataError != nil {
				return queryWithDataError
			}

			var share types.Share
			unmarshalJSONError := codec.UnmarshalJSON(response, &share)
			if unmarshalJSONError != nil {
				return unmarshalJSONError
			}
			return cliContext.PrintOutput(share)
		},
	}

	command.Flags().String(constants.AddressFlag, "", "address")
	return command
}
