/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
)

// NewConsumer : is a consumer which is needed to create child consumers to consume topics
func NewConsumer(kafkaPorts []string) sarama.Consumer {
	config := sarama.NewConfig()

	consumer, Error := sarama.NewConsumer(kafkaPorts, config)
	if Error != nil {
		panic(Error)
	}

	return consumer
}

// PartitionConsumers : is a child consumer
func PartitionConsumers(consumer sarama.Consumer, topic string) sarama.PartitionConsumer {
	// partition and offset defined in CONSTANTS.go
	fmt.Println(consumer.Topics())
	partitionConsumer, Error := consumer.ConsumePartition(topic, partition, offset)
	if Error != nil {
		panic(Error)
	}

	return partitionConsumer
}

// TopicConsumer : Takes a consumer and makes it consume a topic message at a time
func TopicConsumer(topic string, consumers map[string]sarama.PartitionConsumer) ([]byte, error) {
	partitionConsumer := consumers[topic]

	if len(partitionConsumer.Messages()) == 0 {
		return nil, errors.New("No Msgs")
	}

	kafkaMsg := <-partitionConsumer.Messages()
	return kafkaMsg.Value, nil
}
