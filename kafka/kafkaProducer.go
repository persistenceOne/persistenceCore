/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"github.com/Shopify/sarama"

	"github.com/cosmos/cosmos-sdk/codec"
)

// NewProducer is a producer to send messages to kafka
func NewProducer(kafkaPorts []string) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducer(kafkaPorts, nil)
	if err != nil {
		panic(err)
	}

	return producer
}

// KafkaProducerDeliverMessage : delivers messages to kafka
func KafkaProducerDeliverMessage(msg KafkaMsg, topic string, producer sarama.SyncProducer, cdc *codec.LegacyAmino) error {
	kafkaStoreBytes, err := cdc.MarshalJSON(msg)

	if err != nil {
		panic(err)
	}

	sendMsg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(kafkaStoreBytes),
	}
	_, _, err = producer.SendMessage(&sendMsg)

	if err != nil {
		return err
	}

	return nil
}

// SendToKafka : handles sending message to kafka
func SendToKafka(msg KafkaMsg, kafkaState KafkaState, cdc *codec.LegacyAmino) []byte {
	Error := KafkaProducerDeliverMessage(msg, "Topic", kafkaState.Producer, cdc)
	if Error != nil {
		jsonResponse, Error := cdc.MarshalJSON(struct {
			Response string `json:"response"`
		}{Response: "Something is up with kafka server, restart rest and kafka."})
		if Error != nil {
			panic(Error)
		}

		SetTicketIDtoDB(msg.TicketID, kafkaState.KafkaDB, cdc, jsonResponse)
	} else {
		jsonResponse, err := cdc.MarshalJSON(struct {
			Error string `json:"error"`
		}{Error: "Request in process, wait and try after some time"})
		if err != nil {
			panic(err)
		}
		SetTicketIDtoDB(msg.TicketID, kafkaState.KafkaDB, cdc, jsonResponse)
	}

	jsonResponse, Error := cdc.MarshalJSON(struct {
		TicketID Ticket `json:"ticketID"`
	}{TicketID: msg.TicketID})

	if Error != nil {
		panic(Error)
	}

	return jsonResponse
}
