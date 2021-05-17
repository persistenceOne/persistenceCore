package contracts

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceCore/kafka"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
)

var STokens = Contract{
	name:    "S_TOKENS",
	address: constants.STokensAddress,
	abi:     abi.ABI{},
	methods: map[string]func(kafkaState kafka.KafkaState, protoCodec *codec.ProtoCodec, arguments []interface{}) error{
		constants.STokensSetRewards:       onSetRewards,
		constants.STokensCalculateRewards: onCalculateRewards,
	},
}

func onSetRewards(_ kafka.KafkaState, _ *codec.ProtoCodec, arguments []interface{}) error {
	fmt.Printf("onSetRewards: %s\n", arguments[0].(common.Address).String())
	return nil
}

func onCalculateRewards(_ kafka.KafkaState, _ *codec.ProtoCodec, arguments []interface{}) error {
	fmt.Printf("onCalculateRewards: %s\n", arguments[0].(common.Address).String())
	return nil
}
