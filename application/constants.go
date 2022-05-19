/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package application

import (
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
)

const (
	Name             = "PersistenceCore"
	Bech32MainPrefix = "persistence"
	UpgradeName      = "v3"
	CoinType         = 750

	FullFundraiserPath = "44'/750'/0'/0/0"

	Bech32PrefixAccAddr  = Bech32MainPrefix
	Bech32PrefixAccPub   = Bech32MainPrefix + sdkTypes.PrefixPublic
	Bech32PrefixValAddr  = Bech32MainPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator
	Bech32PrefixValPub   = Bech32MainPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixOperator + sdkTypes.PrefixPublic
	Bech32PrefixConsAddr = Bech32MainPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixConsensus
	Bech32PrefixConsPub  = Bech32MainPrefix + sdkTypes.PrefixValidator + sdkTypes.PrefixConsensus + sdkTypes.PrefixPublic

	authzMsgExec                        = "/cosmos.authz.v1beta1.MsgExec"
	authzMsgGrant                       = "/cosmos.authz.v1beta1.MsgGrant"
	authzMsgRevoke                      = "/cosmos.authz.v1beta1.MsgRevoke"
	bankMsgSend                         = "/cosmos.bank.v1beta1.MsgSend"
	bankMsgMultiSend                    = "/cosmos.bank.v1beta1.MsgMultiSend"
	distrMsgSetWithdrawAddr             = "/cosmos.distribution.v1beta1.MsgSetWithdrawAddress"
	distrMsgWithdrawValidatorCommission = "/cosmos.distribution.v1beta1.MsgWithdrawValidatorCommission"
	distrMsgFundCommunityPool           = "/cosmos.distribution.v1beta1.MsgFundCommunityPool"
	distrMsgWithdrawDelegatorReward     = "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward"
	feegrantMsgGrantAllowance           = "/cosmos.feegrant.v1beta1.MsgGrantAllowance"
	feegrantMsgRevokeAllowance          = "/cosmos.feegrant.v1beta1.MsgRevokeAllowance"
	govMsgVoteWeighted                  = "/cosmos.gov.v1beta1.MsgVoteWeighted"
	govMsgSubmitProposal                = "/cosmos.gov.v1beta1.MsgSubmitProposal"
	govMsgDeposit                       = "/cosmos.gov.v1beta1.MsgDeposit"
	govMsgVote                          = "/cosmos.gov.v1beta1.MsgVote"
	stakingMsgEditValidator             = "/cosmos.staking.v1beta1.MsgEditValidator"
	stakingMsgDelegate                  = "/cosmos.staking.v1beta1.MsgDelegate"
	stakingMsgUndelegate                = "/cosmos.staking.v1beta1.MsgUndelegate"
	stakingMsgBeginRedelegate           = "/cosmos.staking.v1beta1.MsgBeginRedelegate"
	stakingMsgCreateValidator           = "/cosmos.staking.v1beta1.MsgCreateValidator"
	vestingMsgCreateVestingAccount      = "/cosmos.vesting.v1beta1.MsgCreateVestingAccount"
	transferMsgTransfer                 = "/ibc.applications.transfer.v1.MsgTransfer"
	liquidityMsgCreatePool              = "/tendermint.liquidity.v1beta1.MsgCreatePool"
	liquidityMsgSwapWithinBatch         = "/tendermint.liquidity.v1beta1.MsgSwapWithinBatch"
	liquidityMsgDepositWithinBatch      = "/tendermint.liquidity.v1beta1.MsgDepositWithinBatch"
	liquidityMsgWithdrawWithinBatch     = "/tendermint.liquidity.v1beta1.MsgWithdrawWithinBatch"
)
