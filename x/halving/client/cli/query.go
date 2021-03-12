/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/persistenceOne/persistenceCore/x/halving/types"
	"github.com/spf13/cobra"
)

// GetQueryCmd returns the cli query commands for the halving module.
func GetQueryCmd() *cobra.Command {
	halvingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the halving module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	halvingQueryCmd.AddCommand(
		GetCmdQueryParams(),
	)

	return halvingQueryCmd
}

// GetCmdQueryParams implements a command to return the current halving
// parameters.
func GetCmdQueryParams() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current halving parameters",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParameters)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := cliCtx.JSONMarshaler.UnmarshalJSON(res, &params); err != nil {
				return err
			}

			return cliCtx.PrintProto(&params)
		},
	}
}
