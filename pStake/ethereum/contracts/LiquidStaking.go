package contracts

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceCore/pStake/tendermint"
	"log"
	"math/big"

	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/persistenceOne/persistenceCore/kafka"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
)

var LiquidStaking = Contract{
	Name:    "LIQUID_STAKING",
	Address: constants.LiquidStaking,
	ABI:     abi.ABI{},
	Methods: map[string]func(kafkaState kafka.KafkaState, protoCodec *codec.ProtoCodec, arguments []interface{}) error{
		constants.LiquidStakingStake:   onStake,
		constants.LiquidStakingUnStake: onUnStake,
	},
}

func onStake(kafkaState kafka.KafkaState, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	amount := arguments[1].(*big.Int)
	stakeMsg := stakingTypes.NewMsgDelegate(tendermint.Chain.MustGetAddress(), constants.Validator1, sdkTypes.NewCoin(constants.Denom, sdkTypes.NewInt(amount.Int64())))
	msgBytes, err := protoCodec.MarshalInterface(stakeMsg)
	err = kafka.ProducerDeliverMessage(msgBytes, kafka.ToTendermint, kafkaState.Producer)
	if err != nil {
		log.Print("Failed to add msg to kafka queue: ", err)
		return err
	}
	//response, ok, err := tendermint.Chain.SendMsg(stakeMsg)
	//if err != nil {
	//	panic(err)
	//}
	//log.Printf("Sending tx, ok: %v, code: %d, hash: %s\n", ok, response.Code, response.TxHash)
	return nil
}

func onUnStake(kafkaState kafka.KafkaState, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	amount := arguments[1].(*big.Int)
	unStakeMsg := stakingTypes.NewMsgUndelegate(tendermint.Chain.MustGetAddress(), constants.Validator1, sdkTypes.NewCoin(constants.Denom, sdkTypes.NewInt(amount.Int64())))
	msgBytes, err := protoCodec.MarshalInterface(unStakeMsg)
	err = kafka.ProducerDeliverMessage(msgBytes, kafka.ToTendermint, kafkaState.Producer)
	if err != nil {
		log.Print("Failed to add msg to kafka queue: ", err)
		return err
	}
	//response, ok, err := tendermint.Chain.SendMsg(unstakeMsg)
	//if err != nil {
	//	panic(err)
	//}
	//log.Printf("Sending tx, ok: %v, code: %d, hash: %s\n", ok, response.Code, response.TxHash)
	return nil
}
