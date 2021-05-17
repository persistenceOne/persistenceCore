package contracts

import (
	"encoding/hex"
	"github.com/persistenceOne/persistenceCore/kafka"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type ContractI interface {
	GetName() string
	GetAddress() string
	GetABI() abi.ABI
	SetABI(contractABIString string)
	GetMethods() map[string]func(kafkaState kafka.KafkaState, arguments []interface{}) error
	GetMethodAndArguments(inputData []byte) (*abi.Method, []interface{}, error)
}

type Contract struct {
	Name    string
	Address string
	ABI     abi.ABI
	Methods map[string]func(kafkaState kafka.KafkaState, arguments []interface{}) error
}

var _ ContractI = &Contract{}

func (contract *Contract) GetName() string {
	return contract.Name
}

func (contract *Contract) GetAddress() string {
	return contract.Address
}

func (contract *Contract) GetABI() abi.ABI {
	return contract.ABI
}

func (contract *Contract) GetMethods() map[string]func(kafkaState kafka.KafkaState, arguments []interface{}) error {
	return contract.Methods
}

func (contract *Contract) SetABI(contractABIString string) {
	contractABI, err := abi.JSON(strings.NewReader(contractABIString))
	if err != nil {
		log.Fatalln("Unable to decode abi:  " + err.Error())
	}
	contract.ABI = contractABI
}

func (contract *Contract) GetMethodAndArguments(inputData []byte) (*abi.Method, []interface{}, error) {
	txData := hex.EncodeToString(inputData)
	if txData[:2] == "0x" {
		txData = txData[2:]
	}

	decodedSig, err := hex.DecodeString(txData[:8])
	if err != nil {
		log.Fatalf("Unable decode method ID (decodeSig) of %s: %s\n", contract.Name, err.Error())
	}

	method, err := contract.ABI.MethodById(decodedSig)
	if err != nil {
		log.Fatalf("Unable to fetch method of %s: %s\n", contract.Name, err.Error())
	}

	decodedData, err := hex.DecodeString(txData[8:])
	if err != nil {
		log.Fatalf("Unable to decode input data of %s: %s\n", contract.Name, err.Error())
	}

	arguments, err := method.Inputs.Unpack(decodedData)
	return method, arguments, err
}
