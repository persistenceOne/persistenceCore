package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/relayer/relayer"
	"github.com/golang/protobuf/proto"
	"github.com/persistenceOne/persistenceCore/kafka"
	"github.com/persistenceOne/persistenceCore/pStake/constants"
	"github.com/spf13/cobra"
	tmservice "github.com/tendermint/tendermint/libs/service"
	tmTypes "github.com/tendermint/tendermint/types"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func GetCmd(initClientCtx client.Context) *cobra.Command {
	pStakeCommand := &cobra.Command{
		Use:   "pStake [path_to_chain_json]",
		Short: "Start pStake",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			timeout, err := cmd.Flags().GetString(constants.FlagTimeOut)
			if err != nil {
				log.Fatalln(err)
			}
			ports, err := cmd.Flags().GetString("ports")
			fmt.Println(ports, err)
			if err != nil {
				return err
			}
			portsList := strings.Split(ports, ",")
			kafkaState := kafka.NewKafkaState(portsList)
			go kafkaRoutine(kafkaState)
			run(initClientCtx, args[0], timeout, kafkaState)
			return nil
		},
	}
	pStakeCommand.Flags().String(constants.FlagTimeOut, "10s", "timeout time for connecting to rpc")
	pStakeCommand.Flags().String("ports", "localhost:9092", "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")

	return pStakeCommand
}

// kafkaRoutine: starts kafka in a separate goRoutine, consumers will each start in different go routines
// no need to store any db, producers and consumers are inside kafkaState struct.
// use kafka.ProducerDeliverMessage() -> to produce message
// use kafka.TopicConsumer -> to consume messages.
func kafkaRoutine(kafkaState kafka.KafkaState) {
	go consumeMsgSend(kafkaState)
	// go consume other messages

	fmt.Println("started consumers")
}
func consumeMsgSend(state kafka.KafkaState) {
	for {
		//consume logic here.
		var msgs []banktypes.MsgSend
		for i := 0; i < kafka.BatchSize; {
			bz, _ := kafka.TopicConsumer(kafka.MsgSendForward, state.Consumers)
			fmt.Println("message received from kafka", bz)
			if bz != nil {
				var msg = banktypes.MsgSend{}
				err := proto.Unmarshal(bz, &msg)
				if err != nil {
					panic(err)
				}
				msgs = append(msgs, msg)
				i++
			} else {
				time.Sleep(kafka.SleepTimer)
			}

		}
		fmt.Println("batch the messages: ", msgs)
		time.Sleep(kafka.SleepRoutine)
	}
}

func run(initClientCtx client.Context, chainConfigJsonPath, timeout string, kafkaState kafka.KafkaState) {
	chain, err := fileInputAdd(chainConfigJsonPath)
	to, err := time.ParseDuration(timeout)
	if err != nil {
		log.Fatalf("Error while parsing timeout: %w", err)
	}
	homePath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error while getting current directory: %w", err)
	}
	err = chain.Init(homePath, to, nil, true)
	if err != nil {
		log.Fatalln(err.Error())
	}
	if err = chain.Start(); err != nil {
		if err != tmservice.ErrAlreadyStarted {
			chain.Error(err)
			return
		}
	}

	txxEvents, txCancel, err := chain.Subscribe(constants.TxEvents)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer txCancel()

	blockEvents, blockCancel, err := chain.Subscribe(constants.BlockEvents)
	if err != nil {
		chain.Error(err)
		return
	}
	defer blockCancel()

	for {
		select {
		case txEvent := <-txxEvents:
			if txEvent.Data.(tmTypes.EventDataTx).Result.Code == 0 {
				go handleEncodeTx(initClientCtx, txEvent.Data.(tmTypes.EventDataTx).Tx, kafkaState)
			}
		case blockEvent := <-blockEvents:
			fmt.Println(blockEvent.Data.(tmTypes.EventDataNewBlock).Block.Height)
		}
	}

}

func handleEncodeTx(initClientCtx client.Context, encodedTx []byte, kafkaState kafka.KafkaState) {
	// Should be used if encodedTx is string
	//decodedTx, err := base64.StdEncoding.DecodeString(encodedTx)
	//if err != nil {
	//	log.Fatalln(err.Error())
	//}

	txInterface, err := initClientCtx.TxConfig.TxDecoder()(encodedTx)
	if err != nil {
		log.Fatalln(err.Error())
	}

	tx, ok := txInterface.(signing.Tx)
	if !ok {
		log.Fatalln("Unable to parse tx")
	}

	fmt.Printf("Memo: %s", tx.GetMemo())

	for _, msg := range tx.GetMsgs() {
		switch msg := msg.(type) {
		case *banktypes.MsgSend:
			if true {
				// produce to send queue
				msgBytes, err := proto.Marshal(msg)

				if err != nil {
					panic(err)
				}
				err = kafka.ProducerDeliverMessage(msgBytes, kafka.MsgSendForward, kafkaState.Producer)
				if err != nil {
					log.Print("Failed to add msg to kafka queue: ", err)
				}
				fmt.Println("Produced to kafka: ", msg.String())
			} else {
				// reversal queue
			}
		default:

		}
	}
}

func fileInputAdd(file string) (*relayer.Chain, error) {
	// If the user passes in a file, attempt to read the chain config from that file
	c := &relayer.Chain{}
	if _, err := os.Stat(file); err != nil {
		return c, err
	}

	byt, err := ioutil.ReadFile(file)
	if err != nil {
		return c, err
	}

	if err = json.Unmarshal(byt, c); err != nil {
		return c, err
	}

	return c, nil
}
