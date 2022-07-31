/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceCore contributors
 SPDX-License-Identifier: Apache-2.0
*/

package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Name             = "PersistenceCore"
	Bech32MainPrefix = "persistence"
	UpgradeName      = "v3"
	CoinType         = 750
	Purpose          = 44

	Bech32PrefixAccAddr  = Bech32MainPrefix
	Bech32PrefixAccPub   = Bech32MainPrefix + sdk.PrefixPublic
	Bech32PrefixValAddr  = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator
	Bech32PrefixValPub   = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	Bech32PrefixConsAddr = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus
	Bech32PrefixConsPub  = Bech32MainPrefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic

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
)
