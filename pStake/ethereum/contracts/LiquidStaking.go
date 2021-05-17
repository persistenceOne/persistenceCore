package contracts

import (
	"fmt"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceCore/pStake/tendermint"
	"math/big"

	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceCore/kafka"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
)

var LiquidStaking = Contract{
	Name:    "LIQUID_STAKING",
	Address: constants.LiquidStaking,
	ABI:     abi.ABI{},
	Methods: map[string]func(kafkaState kafka.KafkaState, arguments []interface{}) error{
		constants.LiquidStakingStake:   onStake,
		constants.LiquidStakingUnStake: onUnStake,
	},
}

func onStake(kafkaState kafka.KafkaState, arguments []interface{}) error {
	fmt.Printf("Eth Stake Tx from Address: %s\n", arguments[0].(common.Address).String())
	amount := arguments[1].(*big.Int)
	stakeMsg := stakingTypes.NewMsgDelegate(tendermint.Chain.MustGetAddress(), constants.Validator1, sdkTypes.NewCoin("stake", sdkTypes.NewInt(amount.Int64())))
	response, ok, err := tendermint.Chain.SendMsg(stakeMsg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sending tx, ok: %v, code: %d, hash: %s\n", ok, response.Code, response.TxHash)
	return nil
}

func onUnStake(kafkaState kafka.KafkaState, arguments []interface{}) error {
	fmt.Printf("Eth Un-Stake Tx from Address: %s\n", arguments[0].(common.Address).String())
	amount := arguments[1].(*big.Int)
	stakeMsg := stakingTypes.NewMsgUndelegate(tendermint.Chain.MustGetAddress(), constants.Validator1, sdkTypes.NewCoin("stake", sdkTypes.NewInt(amount.Int64())))
	response, ok, err := tendermint.Chain.SendMsg(stakeMsg)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Sending tx, ok: %v, code: %d, hash: %s\n", ok, response.Code, response.TxHash)
	return nil
}
