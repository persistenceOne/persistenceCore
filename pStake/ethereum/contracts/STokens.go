package contracts

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
)

var STokens = Contract{
	Name:    "S_TOKENS",
	Address: constants.STokensAddress,
	ABI:     abi.ABI{},
	Methods: map[string]func(arguments []interface{}) error{
		constants.STokensSetRewards:       onSetRewards,
		constants.STokensCalculateRewards: onCalculateRewards,
	},
}

func onSetRewards(arguments []interface{}) error {
	fmt.Printf("onSetRewards: %s\n", arguments[0].(common.Address).String())
	return nil
}

func onCalculateRewards(arguments []interface{}) error {
	fmt.Printf("onCalculateRewards: %s\n", arguments[0].(common.Address).String())
	return nil
}
