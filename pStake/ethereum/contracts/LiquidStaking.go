package contracts

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/persistenceOne/persistenceCore/kafka/utils"
	"log"
	"math/big"

	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
)

var LiquidStaking = Contract{
	name:    "LIQUID_STAKING",
	address: constants.LiquidStakingAddress,
	abi:     abi.ABI{},
	methods: map[string]func(kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec, arguments []interface{}) error{
		constants.LiquidStakingStake:   onStake,
		constants.LiquidStakingUnStake: onUnStake,
	},
}

func onStake(kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	amount := arguments[1].(*big.Int)
	stakeMsg := stakingTypes.NewMsgDelegate(constants.PSTakeAddress, constants.Validator1, sdkTypes.NewCoin(constants.PSTakeDenom, sdkTypes.NewInt(amount.Int64())))
	msgBytes, err := protoCodec.MarshalInterface(sdkTypes.Msg(stakeMsg))
	err = utils.ProducerDeliverMessage(msgBytes, utils.ToTendermint, kafkaState.Producer)
	if err != nil {
		log.Print("Failed to add msg to kafka queue: ", err)
		return err
	}
	return nil
}

func onUnStake(kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec, arguments []interface{}) error {
	amount := arguments[1].(*big.Int)
	unStakeMsg := stakingTypes.NewMsgUndelegate(constants.PSTakeAddress, constants.Validator1, sdkTypes.NewCoin(constants.PSTakeDenom, sdkTypes.NewInt(amount.Int64())))
	msgBytes, err := protoCodec.MarshalInterface(sdkTypes.Msg(unStakeMsg))
	err = utils.ProducerDeliverMessage(msgBytes, utils.EthUnbond, kafkaState.Producer)
	if err != nil {
		log.Print("Failed to add msg to kafka queue: ", err)
		return err
	}
	return nil
}
