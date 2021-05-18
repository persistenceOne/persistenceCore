package tendermint

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
	"github.com/persistenceOne/persistenceCore/kafka/utils"
	"log"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/persistenceOne/persistenceCore/pStake/status"
)

func StartListening(initClientCtx client.Context, chain *relayer.Chain, kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec, sleepDuration time.Duration) {
	ctx := context.Background()

	for {
		abciInfo, err := chain.Client.ABCIInfo(ctx)
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

			txSearchResult, err := chain.Client.TxSearch(ctx, fmt.Sprintf("tx.height=%d", processHeight), true, nil, nil, "asc")
			if err != nil {
				log.Println(err)
				time.Sleep(sleepDuration)
				continue
			}

			err = handleTxSearchResult(initClientCtx, txSearchResult, kafkaState, protoCodec)
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
