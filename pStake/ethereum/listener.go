package ethereum

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	//"strings"
)

func StartListening(ethereumEndPoint string) {
	client, err := ethclient.Dial(ethereumEndPoint)
	if err != nil {
		log.Fatalf("Error while dialing to eth node %s: %s\n", ethereumEndPoint, err.Error())
	}
	ctx := context.Background()

	block, err := client.BlockByNumber(ctx, big.NewInt(4784014))
	if err != nil {
		log.Fatalln(err)
	}

	handleBlock(client, &ctx, block)

	//msg, err := block.Transactions()[0].AsMessage(types.NewEIP155Signer(block.Transactions()[0].ChainId()))
	//fmt.Println(msg.From())

	//height, err :=  client.BlockNumber(ctx)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//query := ethereum.FilterQuery{
	//	FromBlock: nil,
	//	ToBlock:   nil,
	//	Topics: [][]common.Hash,
	//	Addresses: []common.Address{address},
	//}

	var logs = make(chan types.Log)
	headers := make(chan *types.Header)

	//sub, err := client.SubscribeFilterLogs(ctx, query, logs)
	sub, err := client.SubscribeNewHead(ctx, headers)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Printf("Logs subscription error: %s\n", err.Error())
			break
		case l := <-logs:
			fmt.Println("new log", l.Address)
		case header := <-headers:
			go getAndHandleBlock(client, &ctx, header)
		}
	}
}

func getAndHandleBlock(client *ethclient.Client, ctx *context.Context, header *types.Header) {
	block, err := client.BlockByHash(*ctx, header.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Eth new block: %s\n", header.Number.String())
	handleBlock(client, ctx, block)
}
