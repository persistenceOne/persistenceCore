package ethereum

import (
	"context"
	"github.com/persistenceOne/persistenceCore/pStake/ethereum/contracts"
	"log"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func handleBlock(client *ethclient.Client, ctx *context.Context, block *types.Block) {
	for _, transaction := range block.Transactions() {
		if transaction.To() != nil {
			switch transaction.To().String() {
			case contracts.STokens.Address:
				handleTransaction(client, ctx, transaction, &contracts.STokens)
			default:

			}
		}
	}
}

func handleTransaction(client *ethclient.Client, ctx *context.Context, transaction *types.Transaction, contract contracts.ContractI) {
	receipt, err := client.TransactionReceipt(*ctx, transaction.Hash())
	if err != nil {
		log.Fatalf("Error while fetching receipt of tx %s: %s", transaction.Hash().String(), err.Error())
	}

	if receipt.Status == 1 {
		method, arguments, err := contract.GetMethodAndArguments(transaction.Data())
		if err != nil {
			log.Fatalf("Error in getting method and arguments of %s,: %s\n", contract.GetName(), err.Error())
		}

		if processFunc, ok := contract.GetMethods()[method.RawName]; ok {
			err = processFunc(arguments)
			if err != nil {
				log.Fatalf("Error in processing arguments of contarct %s and method  %s,: %s\n", contract.GetName(), method.RawName, err.Error())
			}
		}
	}
}
