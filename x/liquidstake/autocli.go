package liquidstake

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	liquidstakev1beta1 "github.com/persistenceOne/persistenceCore/v17/x/liquidstake/types"
)

// AutoCLIOptions returns the AutoCLI configuration for the liquidstake module.
// This enables `pstaked` to auto-generate CLI commands for this module based
// on the Msg and Query protobuf service definitions.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: liquidstakev1beta1.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the liquidstake module parameters",
				},
				{
					RpcMethod: "LiquidValidators",
					Use:       "liquid-validators",
					Short:     "Query all liquid validators",
				},
				{
					RpcMethod: "States",
					Use:       "states",
					Short:     "Query liquidstake module states (net amount, mint rate)",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: liquidstakev1beta1.Msg_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "LiquidStake",
					Use:       "liquid-stake [amount]",
					Short:     "Liquid-stake XPRT",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "amount"},
					},
				},
				{
					RpcMethod: "StakeToLP",
					Use:       "stake-to-lp [validator-addr] [staked_amount] [liquid_amount]",
					Short:     "Convert delegation into stkXPRT and lock into LP",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "validator_address"},
						{ProtoField: "staked_amount"},
						{ProtoField: "liquid_amount", Optional: true},
					},
				},
				{
					RpcMethod: "LiquidUnstake",
					Use:       "liquid-unstake [amount]",
					Short:     "Liquid-unstake stkXPRT",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "amount"},
					},
				},
				// Authority gated admin txs are kept but discoverable via autocli; typically used with --from authority
				{
					RpcMethod: "UpdateParams",
					Use:       "update-params",
					Short:     "Update liquidstake module params (authority only)",
				},
				{
					RpcMethod: "UpdateWhitelistedValidators",
					Use:       "update-whitelisted-validators",
					Short:     "Update whitelisted validators (authority only)",
				},
				{
					RpcMethod: "SetModulePaused",
					Use:       "set-module-paused",
					Short:     "Set module paused flag (authority only)",
				},
			},
		},
	}
}
