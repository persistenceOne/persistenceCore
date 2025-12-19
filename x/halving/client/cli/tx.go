/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package cli

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/persistenceCore/v17/x/halving/types"
)

func GetTxCmd() *cobra.Command {
	halvingTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Halving transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	halvingTxCmd.AddCommand(
		GetCmdUpdateParams(),
	)

	return halvingTxCmd
}

func GetCmdUpdateParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [num-blocks]",
		Short: "Update halving parameters",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			blockHeight, ok := math.NewIntFromString(args[0])
			if !ok {
				return errors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid block height: %s", args[0])
			}

			params := types.NewParams(blockHeight.Uint64())

			msg := types.NewMsgUpdateParams(
				clientCtx.GetFromAddress().String(),
				params,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
