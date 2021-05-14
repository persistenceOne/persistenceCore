package ethereum

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceCore/pStake/status"
	"log"
	"math/big"
	"time"
	//"strings"
)

func StartListening(ethereumEndPoint string, sleepDuration time.Duration) {
	client, err := ethclient.Dial(ethereumEndPoint)
	if err != nil {
		log.Fatalf("Error while dialing to eth node %s: %s\n", ethereumEndPoint, err.Error())
	}
	ctx := context.Background()

	for {
		latestEthHeight, err := client.BlockNumber(ctx)
		if err != nil {
			log.Printf("Error while fetching latest block height: %s\n", err.Error())
			time.Sleep(sleepDuration)
			continue
		}

		ethStatus, err := status.GetEthereumStatus()
		if err != nil {
			panic(err)
		}

		if latestEthHeight > uint64(ethStatus.LastCheckHeight) {
			processHeight := big.NewInt(ethStatus.LastCheckHeight + 1)
			fmt.Printf("ETH: %d\n", processHeight)

			block, err := client.BlockByNumber(ctx, processHeight)
			if err != nil {
				log.Println(err)
				time.Sleep(sleepDuration)
				continue
			}

			err = handleBlock(client, &ctx, block)
			if err != nil {
				panic(err)
			}

			err = status.SetEthereumStatus(processHeight.Int64())
			if err != nil {
				panic(err)
			}
		}
		time.Sleep(sleepDuration)
	}
}
