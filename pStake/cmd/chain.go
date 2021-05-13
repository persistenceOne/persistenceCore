package cmd

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/persistenceOne/persistenceCore/kafka"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	"github.com/persistenceOne/persistenceCore/pStake/cosmos"
	"github.com/persistenceOne/persistenceCore/pStake/ethereum"
	"github.com/spf13/cobra"
	tmTypes "github.com/tendermint/tendermint/types"
	"log"
	"strings"
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

			ethereumEndPoint, err := cmd.Flags().GetString(constants.FlagEthereumEndPoint)
			if err != nil {
				log.Fatalln(err)
			}

			ports, err := cmd.Flags().GetString("ports")
			fmt.Println(ports, err)
			kafkaHome, err := cmd.Flags().GetString(kafka.FlagKafkaHome)

			protoCodec := codec.NewProtoCodec(initClientCtx.InterfaceRegistry)

			if err != nil {
				return err
			}
			portsList := strings.Split(ports, ",")
			kafkaState := kafka.NewKafkaState(portsList, kafkaHome)
			go kafkaRoutine(kafkaState, protoCodec)
			server.TrapSignal(kafkaClose(kafkaState))

			log.Println("Starting to listen ethereum....")
			go ethereum.StartListening(ethereumEndPoint)

			run(initClientCtx, args[0], timeout, homePath, coinType, args[1], kafkaState)

			return nil
		},
	}
	pStakeCommand.Flags().String(constants.FlagTimeOut, "10s", "timeout time for connecting to rpc")
	pStakeCommand.Flags().Uint32(constants.FlagCoinType, 118, "coin type for wallet")
	pStakeCommand.Flags().String(constants.FlagHome, "./pStake", "home for pStake")
	pStakeCommand.Flags().String(constants.FlagEthereumEndPoint, "wss://goerli.infura.io/ws/v3/e2549c9ec9764e46a7768cc7619a1939", "ethereum node to connect")
	pStakeCommand.Flags().String("ports", "localhost:9092", "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")
	pStakeCommand.Flags().String(kafka.FlagKafkaHome, kafka.DefaultKafkaHome, "The kafka config file directory")
	return pStakeCommand
}

func run(initClientCtx client.Context, chainConfigJsonPath, timeout, homePath string, coinType uint32, mnemonics string, kafkaState kafka.KafkaState) {
	err := cosmos.InitializeAndStartChain(chainConfigJsonPath, timeout, homePath, coinType, mnemonics)
	if err != nil {
		log.Fatalf("Error while intiializing and starting chain: %s\n", err.Error())
	}

	log.Println("Starting to listen cosmos txs....")
	txEvents, txCancel, err := cosmos.Chain.Subscribe(constants.TxEvents)
	if err != nil {
		log.Fatalf("Error while subscribing to tx events: %s\n", err.Error())
	}
	defer txCancel()

	log.Println("Starting to listen cosmos blocks....")
	blockEvents, blockCancel, err := cosmos.Chain.Subscribe(constants.BlockEvents)
	if err != nil {
		log.Fatalf("Error while subscribing to block events: %s\n", err.Error())
	}
	defer blockCancel()

	for {
		select {
		case txEvent := <-txEvents:
			go cosmos.HandleTxEvent(&initClientCtx, txEvent.Data.(tmTypes.EventDataTx), kafkaState)
		case blockEvent := <-blockEvents:
			go cosmos.HandleNewBlock(cosmos.Chain, &initClientCtx, blockEvent.Data.(tmTypes.EventDataNewBlock), kafkaState)
		}
	}

}

// TODO for Eth events =>
func ethEvents() {
	// msg delegate => convert to MsgDelegate and push to ToEth queue
	// msg unbond =>convert to MsgSend and push to EthUnbond queue

}
