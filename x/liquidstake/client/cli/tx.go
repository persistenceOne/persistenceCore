package cli

// DONTCOVER
// client is excluded from test coverage in MVP version

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"github.com/persistenceOne/persistenceCore/v17/x/liquidstake/types"
)

// GetTxCmd returns a root CLI command handler for all x/liquidstake transaction commands.
func GetTxCmd() *cobra.Command {
	liquidstakeTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Aliases:                    []string{"ls"},
		Short:                      "XPRT liquid stake transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	liquidstakeTxCmd.AddCommand(
		NewLiquidStakeCmd(),
		NewStakeToLPCmd(),
		NewLiquidUnstakeCmd(),
		NewUpdateParamsCmd(),
		NewUpdateWhitelistedValidatorsCmd(),
		NewSetModulePausedCmd(),
	)

	return liquidstakeTxCmd
}

// NewLiquidStakeCmd implements the liquid stake XPRT command handler.
func NewLiquidStakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-stake [amount]",
		Args:  cobra.ExactArgs(1),
		Short: "Liquid-stake XPRT",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Liquid-stake XPRT.

Example:
$ %s tx %s liquid-stake 1000uxprt --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			liquidStaker := clientCtx.GetFromAddress()

			stakingCoin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgLiquidStake(liquidStaker, stakingCoin)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewStakeToLPCmd implements the liquid stake XPRT command handler.
func NewStakeToLPCmd() *cobra.Command {
	bech32PrefixValAddr := sdk.GetConfig().GetBech32ValidatorAddrPrefix()

	cmd := &cobra.Command{
		Use:   "stake-to-lp [validator-addr] [staked_amount] [liquid_amount]",
		Args:  cobra.RangeArgs(2, 3),
		Short: "Convert delegation into stkXPRT and lock into LP.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Convert delegation into stkXPRT and lock into LP.
Allows to specify both staked and non-staked XPRT amounts to convert into stkXPRT and lock into LP.

Examples:
$ %s tx %s stake-to-lp %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 1000uxprt --from mykey
$ %s tx %s stake-to-lp %s1gghjut3ccd8ay0zduzj64hwre2fxs9ldmqhffj 1000uxprt 5000uxprt --from mykey
`,
				version.AppName, types.ModuleName, bech32PrefixValAddr,
				version.AppName, types.ModuleName, bech32PrefixValAddr,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			liquidStaker := clientCtx.GetFromAddress()

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			stakedCoin, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			var liquidCoin sdk.Coin
			if len(args) > 2 {
				liquidCoin, err = sdk.ParseCoinNormalized(args[2])
				if err != nil {
					return err
				}
			}

			msg := types.NewMsgStakeToLP(liquidStaker, valAddr, stakedCoin, liquidCoin)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewLiquidUnstakeCmd implements the liquid unstake XPRT command handler.
func NewLiquidUnstakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquid-unstake [amount]",
		Args:  cobra.ExactArgs(1),
		Short: "Liquid-unstake stkXPRT",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Liquid-unstake stkXPRT.

Example:
$ %s tx %s liquid-unstake 500stk/uxprt --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			liquidStaker := clientCtx.GetFromAddress()

			unstakingCoin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgLiquidUnstake(liquidStaker, unstakingCoin)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUpdateParamsCmd implements the params update command handler.
func NewUpdateParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-params [params.json]",
		Args:  cobra.ExactArgs(1),
		Short: "Update-params of liquidstake module.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`update-params param-file.

Example:
$ %s tx %s update-params ~/params.json --from mykey

Example params.json
{
  "liquid_bond_denom": "stk/uxprt",
  "whitelisted_validators": [
    {
      "validator_address": "persistencevaloper1hcqg5wj9t42zawqkqucs7la85ffyv08lmnhye9",
      "target_weight": "10"
    }
  ],
  "lsm_disabled": false,
  "unstake_fee_rate": "0.000000000000000000",
  "min_liquid_staking_amount": "10000",
  "min_liquid_staking_amount": "10000",
  "cw_locked_pool_address": "persistence14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9sjvz4fk"
}
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var params types.Params

			paramsInFile, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			err = json.Unmarshal(paramsInFile, &params)
			if err != nil {
				return err
			}
			authority := clientCtx.GetFromAddress()

			msg := types.NewMsgUpdateParams(authority, params)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewUpdateWhitelistedValidatorsCmd implements the update of whitelisted validators command handler.
func NewUpdateWhitelistedValidatorsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-whitelisted-validators [validators_list.json]",
		Args:  cobra.ExactArgs(1),
		Short: "Update whitelisted validators in params of liquidstake module.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`update-whitelisted-validators validators_list.json

Example:
$ %s tx %s update-whitelisted-validators ~/validators_list.json --from mykey

Example validators_list.json
[
  {
    "validator_address": "persistencevaloper1hcqg5wj9t42zawqkqucs7la85ffyv08lmnhye9",
    "target_weight": "10"
  }
]
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var validatorsList []types.WhitelistedValidator

			validatorsListFile, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			err = json.Unmarshal(validatorsListFile, &validatorsList)
			if err != nil {
				return err
			}
			authority := clientCtx.GetFromAddress()

			msg := types.NewMsgUpdateWhitelistedValidators(authority, validatorsList)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// NewSetModulePausedCmd implements the  command handler for updating of safety toggle that disables the module.
func NewSetModulePausedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pause-module [flag]",
		Args:  cobra.ExactArgs(1),
		Short: "Pause or unpause the liquidstake module for an emergency updates.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`pause-module [true/false]

Example:
$ %s tx %s pause-module true --from mykey
`,
				version.AppName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			isPaused := false
			if strings.ToLower(args[0]) == "true" {
				isPaused = true
			} else if strings.ToLower(args[0]) != "false" {
				err := fmt.Errorf("expected flag to be true or false – where 'true' means the module is paused")
				return err
			}

			authority := clientCtx.GetFromAddress()
			msg := types.NewMsgSetModulePaused(authority, isPaused)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
