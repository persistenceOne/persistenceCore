package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/relayer/relayer"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/persistenceOne/persistenceCore/kafka/handler"
	"github.com/persistenceOne/persistenceCore/kafka/runconfig"
	"github.com/persistenceOne/persistenceCore/kafka/utils"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func KafkaCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "kafka",
		Short:                      "kafka commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 1,
		RunE:                       func(cmd *cobra.Command, args []string) error { return errors.New("expect a subcommand") },
	}
	cmd.AddCommand(InitCmd())
	cmd.PersistentFlags().String(utils.FlagKafkaHome, utils.DefaultKafkaHome, "The kafka config file directory")
	return cmd
}

func InitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Init kafka config file",
		RunE: func(cmd *cobra.Command, args []string) error {

			config := runconfig.NewKafkaConfig()

			var buf bytes.Buffer
			encoder := toml.NewEncoder(&buf)
			if err := encoder.Encode(config); err != nil {
				panic(err)
			}

			homeDir, err := cmd.Flags().GetString(utils.FlagKafkaHome)
			if err != nil {
				panic(err)
			}
			if err = os.MkdirAll(homeDir, os.ModePerm); err != nil {
				panic(err)
			}
			if err := ioutil.WriteFile(filepath.Join(homeDir, "kafkaConfig.toml"), buf.Bytes(), 0644); err != nil {
				panic(err)
			}
			log.Println("generated config file at ", filepath.Join(homeDir, "kafkaConfig.toml"))

			return nil
		},
	}
	return initCmd
}

// kafkaClose: closes all kafka connections
func kafkaClose(kafkaState utils.KafkaState) func() {
	return func() {
		fmt.Println("closing all kafka clients.")
		if err := kafkaState.Producer.Close(); err != nil {
			log.Print("Error in closing producer:", err)
		}
		for _, consumerGroup := range kafkaState.ConsumerGroup {
			if err := consumerGroup.Close(); err != nil {
				log.Print("Error in closing partition:", err)
			}
		}
		if err := kafkaState.Admin.Close(); err != nil {
			log.Print("Error in closing admin:", err)
		}

	}
}

// kafkaRoutine: starts kafka in a separate goRoutine, consumers will each start in different go routines
// no need to store any db, producers and consumers are inside kafkaState struct.
// use kafka.ProducerDeliverMessage() -> to produce message
// use kafka.TopicConsumer -> to consume messages.
func kafkaRoutine(kafkaState utils.KafkaState, protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client) {
	kafkaConfig := runconfig.KafkaConfig{}

	_, err := toml.DecodeFile(filepath.Join(kafkaState.HomeDir, "kafkaConfig.toml"), &kafkaConfig)
	if err != nil {
		log.Printf("Error decoding kafkaConfig file: %v", err)
	}
	ctx := context.Background()

	go consumeToEthMsgs(ctx, kafkaState, kafkaConfig, protoCodec, chain, ethereumClient)
	go consumeUnbondings(ctx, kafkaState, kafkaConfig, protoCodec, chain, ethereumClient)
	go consumeToTendermintMessages(ctx, kafkaState, kafkaConfig, protoCodec, chain, ethereumClient)
	// go consume other messages

	fmt.Println("started consumers")
}

func consumeToEthMsgs(ctx context.Context, state utils.KafkaState, kafkaConfig runconfig.KafkaConfig,
	protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client) {
	consumerGroup := state.ConsumerGroup[utils.GroupToEth]
	for {
		msgHandler := handler.MsgHandler{KafkaConfig: kafkaConfig, ProtoCodec: protoCodec,
			Chain: chain, EthClient: ethereumClient, Count: 0}
		err := consumerGroup.Consume(ctx, []string{utils.ToEth}, msgHandler)
		if err != nil {
			log.Println("Error in consumer group.Consume", err)
		}
		time.Sleep(time.Duration(1000000000))
	}
}

func consumeToTendermintMessages(ctx context.Context, state utils.KafkaState, kafkaConfig runconfig.KafkaConfig,
	protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client) {
	groupMsgUnbond := state.ConsumerGroup[utils.GroupMsgUnbond]
	groupMsgDelegate := state.ConsumerGroup[utils.GroupMsgDelegate]
	groupMsgSend := state.ConsumerGroup[utils.GroupMsgSend]
	groupMsgToTendermint := state.ConsumerGroup[utils.GroupToTendermint]
	for {
		msgHandler := handler.MsgHandler{KafkaConfig: kafkaConfig, ProtoCodec: protoCodec,
			Chain: chain, EthClient: ethereumClient, Count: 0}
		err := groupMsgUnbond.Consume(ctx, []string{utils.MsgUnbond}, msgHandler)
		if err != nil {
			log.Println("Error in consumer group.Consume for MsgUnbond", err)
		}
		err = groupMsgDelegate.Consume(ctx, []string{utils.MsgDelegate}, msgHandler)
		if err != nil {
			log.Println("Error in consumer group.Consume", err)
		}
		err = groupMsgSend.Consume(ctx, []string{utils.MsgSend}, msgHandler)
		if err != nil {
			log.Println("Error in consumer group.Consume", err)
		}
		err = groupMsgToTendermint.Consume(ctx, []string{utils.ToTendermint}, msgHandler)
		if err != nil {
			log.Println("Error in consumer group.Consume", err)
		}
		time.Sleep(time.Duration(1000000000))
	}
}

func consumeUnbondings(ctx context.Context, state utils.KafkaState, kafkaConfig runconfig.KafkaConfig,
	protoCodec *codec.ProtoCodec, chain *relayer.Chain, ethereumClient *ethclient.Client) {
	ethUnbondConsumerGroup := state.ConsumerGroup[utils.GroupEthUnbond]
	for {
		msgHandler := handler.MsgHandler{KafkaConfig: kafkaConfig, ProtoCodec: protoCodec,
			Chain: chain, EthClient: ethereumClient, Count: 0}
		err := ethUnbondConsumerGroup.Consume(ctx, []string{utils.EthUnbond}, msgHandler)
		if err != nil {
			log.Println("Error in consumer group.Consume for EthUnbond ", err)
		}
		time.Sleep(kafkaConfig.EthUnbondCycleTime)
	}
}
