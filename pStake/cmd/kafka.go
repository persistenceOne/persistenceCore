package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/persistenceOne/persistenceCore/kafka"
	"github.com/persistenceOne/persistenceCore/kafka/runConfig"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"path/filepath"
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
	cmd.PersistentFlags().String(kafka.FlagKafkaHome, kafka.DefaultKafkaHome, "The kafka config file directory")
	return cmd
}

func InitCmd() *cobra.Command {
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Init kafka config file",
		RunE: func(cmd *cobra.Command, args []string) error {

			config := runConfig.NewKafkaConfig()

			var buf bytes.Buffer
			encoder := toml.NewEncoder(&buf)
			if err := encoder.Encode(config); err != nil {
				panic(err)
			}

			homeDir, err := cmd.Flags().GetString(kafka.FlagKafkaHome)
			if err != nil {
				panic(err)
			}
			if err := ioutil.WriteFile(filepath.Join(homeDir, "kafkaConfig.toml"), buf.Bytes(), 0644); err != nil {
				panic(err)
			}

			return nil
		},
	}
	return initCmd
}

// kafkaClose: closes all kafka connections
func kafkaClose(kafkaState kafka.KafkaState) func() {
	return func() {
		fmt.Println("closing all kafka clients.")
		if err := kafkaState.Producer.Close(); err != nil {
			log.Print("Error in closing producer:", err)
		}
		if err := kafkaState.ConsumerGroup.Close(); err != nil {
			log.Print("Error in closing partition:", err)
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
func kafkaRoutine(kafkaState kafka.KafkaState) {
	kafkaConfig := runConfig.KafkaConfig{}

	_, err := toml.DecodeFile(filepath.Join(kafkaState.HomeDir, "kafkaConfig.toml"), &kafkaConfig)
	if err != nil {
		log.Printf("Error decoding kafkaConfig file: %v", err)
	}
	ctx := context.Background()

	go consumeMsgs(ctx, kafkaState, kafkaConfig)
	go consumeUnbondings(ctx, kafkaState, kafkaConfig)
	// go consume other messages

	fmt.Println("started consumers")
}

func consumeMsgs(ctx context.Context, state kafka.KafkaState, kafkaConfig runConfig.KafkaConfig) {
	consumerGroup := state.ConsumerGroup
	for {
		handler := kafka.MsgHandler{KafkaConfig: kafkaConfig}
		err := consumerGroup.Consume(ctx, []string{kafka.ToEth, kafka.ToTendermint}, handler)
		if err != nil {
			log.Println("Error in consumer group.Consume", err)
		}
	}
}
func consumeUnbondings(ctx context.Context, state kafka.KafkaState, kafkaConfig runConfig.KafkaConfig) {
	consumerGroup := state.ConsumerGroup
	for {
		handler := kafka.MsgHandler{KafkaConfig: kafkaConfig}
		err := consumerGroup.Consume(ctx, []string{kafka.EthUnbond}, handler)
		if err != nil {
			log.Println("Error in consumer group.Consume for EthUnbond ", err)
		}
		err = consumerGroup.Consume(ctx, []string{kafka.UnbondPool}, handler)
		if err != nil {
			log.Println("Error in consumer group.Consume for UnbondPool", err)
		}
	}
}
