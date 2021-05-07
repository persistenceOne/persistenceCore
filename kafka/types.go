/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"github.com/Shopify/sarama"
	dbm "github.com/tendermint/tm-db"
)

// Ticket : is a type that implements string
type Ticket string

//// KafkaMsg : is a store that can be stored in kafka queues
//type KafkaMsg struct {
//	Msg      sdk.Msg `json:"msg"`
//	TicketID Ticket  `json:"ticketID"`
//}
//
//// NewKafkaMsgFromRest : makes a msg to send to kafka queue
//func NewKafkaMsgFromRest(msg sdk.Msg, ticketID Ticket) KafkaMsg {
//	return KafkaMsg{
//		Msg:      msg,
//		TicketID: ticketID,
//	}
//}

// TicketIDResponse : is a json structure to send TicketID to user
type TicketIDResponse struct {
	TicketID Ticket `json:"ticketID" valid:"required~ticketID is mandatory,length(20)~ticketID length should be 20" `
}

// KafkaState : is a struct showing the state of kafka
type KafkaState struct {
	KafkaDB   *dbm.GoLevelDB
	Admin     sarama.ClusterAdmin
	Consumer  sarama.Consumer
	Consumers map[string]sarama.PartitionConsumer
	Producer  sarama.SyncProducer
	Topics    []string
}

// NewKafkaState : returns a kafka state
func NewKafkaState(kafkaPorts []string) KafkaState {
	kafkaDB, _ := dbm.NewGoLevelDB("KafkaDB", DefaultCLIHome)
	admin := KafkaAdmin(kafkaPorts)
	adminTopics, err := admin.ListTopics()
	if err != nil {
		panic(err)
	}
	//create topics if not present
	for _, topic := range Topics {
		if _, ok := adminTopics[topic]; !ok {
			TopicsInit(admin, topic)
		}
	}
	producer := NewProducer(kafkaPorts)
	consumer := NewConsumer(kafkaPorts)
	var consumers = make(map[string]sarama.PartitionConsumer)
	for _, topic := range Topics {
		partitionConsumer := PartitionConsumers(consumer, topic)
		consumers[topic] = partitionConsumer
	}

	return KafkaState{
		KafkaDB:   kafkaDB,
		Admin:     admin,
		Consumer:  consumer,
		Consumers: consumers,
		Producer:  producer,
		Topics:    Topics,
	}
}
