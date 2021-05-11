package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/persistenceOne/persistenceCore/kafka"
	"github.com/persistenceOne/persistenceCore/kafka/runConfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
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

			homeDir, err := cmd.Flags().GetString(flags.FlagHome)
			if err != nil {
				panic(err)
			}
			if err := ioutil.WriteFile(filepath.Join(homeDir, "config", "kafkaConfig.toml"), buf.Bytes(), 0644); err != nil {
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
	go consumeMsgs(kafkaState)
	// go consume other messages

	fmt.Println("started consumers")
}

func consumeMsgs(state kafka.KafkaState) {
	kafkaConfig := runConfig.KafkaConfig{}

	homeDir := viper.GetString(flags.FlagHome)
	_, err := toml.DecodeFile(filepath.Join(homeDir, "config", "kafkaConfig.toml"), &kafkaConfig)
	if err != nil {
		log.Printf("Error decoding kafkaConfig file: %v", err)
	}
	consumerGroup := state.ConsumerGroup
	ctx := context.Background()

	for {
		handler := kafka.MsgHandler{KafkaConfig: kafkaConfig}
		err := consumerGroup.Consume(ctx, kafka.Topics, handler)
		if err != nil {
			log.Println("Error in consumer group.Consume", err)
		}
		time.Sleep(kafka.SleepRoutine)
	}
}
