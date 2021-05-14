package cmd

import (
	"context"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/persistenceOne/persistenceCore/kafka"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	"github.com/persistenceOne/persistenceCore/pStake/ethereum"
	"github.com/persistenceOne/persistenceCore/pStake/status"
	"github.com/persistenceOne/persistenceCore/pStake/tendermint"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

func GetCmd(initClientCtx client.Context) *cobra.Command {
	pStakeCommand := &cobra.Command{
		Use:   "pStake [path_to_chain_json] [mnemonics]",
		Short: "Start pStake",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {

			timeout, err := cmd.Flags().GetString(constants.FlagTimeOut)
			if err != nil {
				log.Fatalln(err)
			}

			coinType, err := cmd.Flags().GetUint32(constants.FlagCoinType)
			if err != nil {
				log.Fatalln(err)
			}

			homePath, err := cmd.Flags().GetString(constants.FlagHome)
			if err != nil {
				log.Fatalln(err)
			}

			cosmosSleepTime, err := cmd.Flags().GetInt(constants.FlagCosmosSleepTime)
			if err != nil {
				log.Fatalln(err)
			}
			cosmosSleepDuration := time.Duration(cosmosSleepTime) * time.Millisecond

			ethereumEndPoint, err := cmd.Flags().GetString(constants.FlagEthereumEndPoint)
			if err != nil {
				log.Fatalln(err)
			}

			ethSleepTime, err := cmd.Flags().GetInt(constants.FlagEthereumSleepTime)
			if err != nil {
				log.Fatalln(err)
			}
			ethSleepDuration := time.Duration(ethSleepTime) * time.Millisecond

			ports, err := cmd.Flags().GetString("ports")
			if err != nil {
				log.Fatalln(err)
			}

			kafkaHome, err := cmd.Flags().GetString(kafka.FlagKafkaHome)
			if err != nil {
				log.Fatalln(err)
			}

			db, err := status.InitializeDB(homePath+"/db", 5732, 4789782)
			if err != nil {
				log.Fatalln(err)
			}
			defer db.Close()

			protoCodec := codec.NewProtoCodec(initClientCtx.InterfaceRegistry)
			portsList := strings.Split(ports, ",")
			kafkaState := kafka.NewKafkaState(portsList, kafkaHome)
			go kafkaRoutine(kafkaState, protoCodec)
			server.TrapSignal(kafkaClose(kafkaState))

			log.Println("Starting to listen ethereum....")
			go ethereum.StartListening(ethereumEndPoint, ethSleepDuration)

			run(initClientCtx.WithHomeDir(homePath), args[0], timeout, homePath, coinType, args[1], kafkaState, cosmosSleepDuration)

			return nil
		},
	}
	pStakeCommand.Flags().String(constants.FlagTimeOut, "10s", "timeout time for connecting to rpc")
	pStakeCommand.Flags().Uint32(constants.FlagCoinType, 118, "coin type for wallet")
	pStakeCommand.Flags().String(constants.FlagHome, "./pStake", "home for pStake")
	pStakeCommand.Flags().String(constants.FlagEthereumEndPoint, "wss://goerli.infura.io/ws/v3/e2549c9ec9764e46a7768cc7619a1939", "ethereum node to connect")
	pStakeCommand.Flags().String("ports", "localhost:9092", "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")
	pStakeCommand.Flags().String(kafka.FlagKafkaHome, kafka.DefaultKafkaHome, "The kafka config file directory")
	pStakeCommand.Flags().Int(constants.FlagCosmosSleepTime, 3000, "sleep time between block checking for cosmos in ms (default 3000 ms)")
	pStakeCommand.Flags().Int(constants.FlagEthereumSleepTime, 4000, "sleep time between block checking for cosmos in ms (default 4000 ms)")
	return pStakeCommand
}

func run(initClientCtx client.Context, chainConfigJsonPath, timeout, homePath string, coinType uint32, mnemonics string, kafkaState kafka.KafkaState, sleepDuration time.Duration) {
	err := tendermint.InitializeAndStartChain(chainConfigJsonPath, timeout, homePath, coinType, mnemonics)
	if err != nil {
		log.Fatalf("Error while intiializing and starting chain: %s\n", err.Error())
	}

	ctx := context.Background()

	for {
		abciInfo, err := tendermint.Chain.Client.ABCIInfo(ctx)
		if err != nil {
			log.Printf("Error while fetching tendermint abci info: %s\n", err.Error())
			time.Sleep(sleepDuration)
			continue
		}
		fmt.Printf("TM new block: %d\n", abciInfo.Response.LastBlockHeight)

		cosmosStatus, err := status.GetCosmosStatus()
		if err != nil {
			panic(err)
		}

		if abciInfo.Response.LastBlockHeight > cosmosStatus.LastCheckHeight {
			processHeight := cosmosStatus.LastCheckHeight + 1
			fmt.Printf("Processing TM: %d\n", processHeight)

			txSearchResult, err := tendermint.Chain.Client.TxSearch(ctx, fmt.Sprintf("tx.height=%d", processHeight), true, nil, nil, "asc")
			if err != nil {
				log.Println(err)
				time.Sleep(sleepDuration)
				continue
			}

			err = tendermint.HandleTxSearchResult(initClientCtx, txSearchResult, kafkaState)
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

// TODO for Eth events =>
func ethEvents() {
	// msg delegate => convert to MsgDelegate and push to ToEth queue
	// msg unbond =>convert to MsgSend and push to EthUnbond queue

}
