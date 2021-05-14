package tendermint

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/persistenceOne/persistenceCore/kafka"
	"github.com/persistenceOne/persistenceCore/pStake/status"
)

func StartListening(initClientCtx client.Context, chainConfigJsonPath, timeout, homePath string, coinType uint32, mnemonics string, kafkaState kafka.KafkaState, sleepDuration time.Duration) {
	err := InitializeAndStartChain(chainConfigJsonPath, timeout, homePath, coinType, mnemonics)
	if err != nil {
		log.Fatalf("Error while intiializing and starting chain: %s\n", err.Error())
	}

	ctx := context.Background()

	for {
		abciInfo, err := Chain.Client.ABCIInfo(ctx)
		if err != nil {
			log.Printf("Error while fetching tendermint abci info: %s\n", err.Error())
			time.Sleep(sleepDuration)
			continue
		}

		cosmosStatus, err := status.GetCosmosStatus()
		if err != nil {
			panic(err)
		}

		if abciInfo.Response.LastBlockHeight > cosmosStatus.LastCheckHeight {
			processHeight := cosmosStatus.LastCheckHeight + 1
			fmt.Printf("TM: %d\n", processHeight)

			txSearchResult, err := Chain.Client.TxSearch(ctx, fmt.Sprintf("tx.height=%d", processHeight), true, nil, nil, "asc")
			if err != nil {
				log.Println(err)
				time.Sleep(sleepDuration)
				continue
			}

			err = handleTxSearchResult(initClientCtx, txSearchResult, kafkaState)
			if err != nil {
				panic(err)
			}

			err = status.SetCosmosStatus(processHeight)
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(sleepDuration)
	}
}
