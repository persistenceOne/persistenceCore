package eth

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceCore/pStake/sTokens"
	"io/ioutil"
	"log"
	"math/big"
	"strings"
	//"strings"
)

func StartListening() {
	infura := "wss://goerli.infura.io/ws/v3/e2549c9ec9764e46a7768cc7619a1939"
	//address := common.HexToAddress("0x925d092d9ff6c95eab70ee5a23c77f355c67f46d")
	cl, err := ethclient.Dial(infura)
	if err != nil {
		log.Fatalln(err)
	}
	ctx := context.Background()

	block, err := cl.BlockByNumber(ctx, big.NewInt(4772102))
	if err != nil {
		log.Fatalln(err)
	}

	liquidStakingByte, err := ioutil.ReadFile(sTokens.STokensABI)
	if err != nil {
		log.Fatal(err)
	}

	contractABI, err := abi.JSON(strings.NewReader(string(liquidStakingByte)))
	if err != nil {
		log.Fatal("Unable to decode abi:  " + err.Error())
	}

	fmt.Println(contractABI.Receive.Name)
	for _, tx := range block.Transactions() {
		if tx.To().String() == "0xB1ab2F588Fe8198D1A8459d9C1d9457a47a811C7" {
			fmt.Println(tx.Hash())
			txData := hex.EncodeToString(tx.Data())
			fmt.Println(txData)

			//decodedSig, err := hex.DecodeString(txData)
			//if err != nil {
			//	log.Fatal(err)
			//}
			//
			//method, err := sTokens.(decodedSig)
			//if err != nil {
			//	log.Fatal(err)
			//}
			//
			//decodedData, err := hex.DecodeString(txData[8:])
			//if err != nil {
			//	log.Fatal(err)
			//}

			//type FunctionInputs struct {
			//	Relay        string // *big.Int for uint256 for example
			//	UnstakeDelay *big.Int
			//}

			//var data sTokens.STokens

			// unpack method inputs
			//a, err := method.Inputs.Unpack(decodedData)
			//if err != nil {
			//	log.Fatal(err)
			//}
			//fmt.Println(a)
			//fmt.Println(data)

			receipt, err := cl.TransactionReceipt(ctx, tx.Hash())
			if err != nil {
				log.Println("error decoding data: " + err.Error())
			} else {
				fmt.Println(receipt.Status)
			}
		}
	}

	//fmt.Println(block.Transactions()[0].Hash())
	//fmt.Println(block.Transactions()[0].To())
	//msg, err := block.Transactions()[0].AsMessage(types.NewEIP155Signer(block.Transactions()[0].ChainId()))
	//fmt.Println(msg.From())

	//height, err :=  cl.BlockNumber(ctx)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//query := ethereum.FilterQuery{
	//	FromBlock: nil,
	//	ToBlock:   nil,
	//	// Topics [][]common.Hash
	//	Addresses: []common.Address{address},
	//}

	//var logs = make(chan types.Log)
	//headers := make(chan *types.Header)

	//sub, err := cl.SubscribeFilterLogs(ctx, query, logs)
	//sub, err := cl.SubscribeNewHead(ctx, headers)
	//if err != nil {
	//	log.Fatalln(err)
	//}

	//for {
	//	select {
	//	case err := <-sub.Err():
	//		fmt.Println("Logs subscription error", err.Error())
	//		break
	//	case l := <-logs:
	//		fmt.Println("new log", l.Address)
	//	case header := <-headers:
	//		block, err := cl.BlockByHash(ctx, header.Hash())
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		fmt.Println(block.Number())
	//		fmt.Println(block.Transactions().Len())
	//	}
	//}

}

func decodeTxParams(abi abi.ABI, v map[string]interface{}, data []byte) (map[string]interface{}, error) {
	m, err := abi.MethodById(data[:4])
	if err != nil {
		return map[string]interface{}{}, err
	}
	if err := m.Inputs.UnpackIntoMap(v, data[4:]); err != nil {
		return map[string]interface{}{}, err
	}
	return v, nil
}
