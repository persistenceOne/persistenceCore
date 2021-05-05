package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/relayer/relayer"
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
			go kafkaRoutine(ports)
			run(initClientCtx, args[0], timeout)
			return nil
		},
	}
	pStakeCommand.Flags().String(constants.FlagTimeOut, "10s", "timeout time for connecting to rpc")
	pStakeCommand.Flags().String("ports", "localhost:9092", "ports kafka brokers are running on, --ports 192.100.10.10:443,192.100.10.11:443")

	return pStakeCommand
}

// kafkaRoutine: starts kafka in a separate goRoutine, consumers will each start in different go routines
// no need to store any db, producers and consumers are inside kafkaState struct.
// use kafka.KafkaProducerDeliverMessage() -> to produce message
// use kafka.KafkaTopicConsumer -> to consume messages.
func kafkaRoutine(ports string) {
	portsList := strings.Split(ports, ",")
	_ = kafka.NewKafkaState(portsList)

	time.Sleep(1000000000)

	go consumeMsgSend()
	// go consume other messages

	fmt.Println("started consumers")
}
func consumeMsgSend() {
	for {
		//consume logic here.
		time.Sleep(kafka.SleepRoutine)
	}
}

func run(initClientCtx client.Context, chainConfigJsonPath, timeout string) {
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
				go handleEncodeTx(initClientCtx, txEvent.Data.(tmTypes.EventDataTx).Tx)
			}
		case blockEvent := <-blockEvents:
			fmt.Println(blockEvent.Data.(tmTypes.EventDataNewBlock).Block.Height)
		}
	}

}

func handleEncodeTx(initClientCtx client.Context, encodedTx []byte) {
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
			fmt.Println(msg.String())
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
