/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"github.com/Shopify/sarama"
)

// NewProducer is a producer to send messages to kafka
func NewProducer(kafkaPorts []string, config *sarama.Config) sarama.SyncProducer {
	producer, err := sarama.NewSyncProducer(kafkaPorts, config)
	if err != nil {
		panic(err)
	}

	return producer
}

// ProducerDeliverMessage : delivers messages to kafka
func ProducerDeliverMessage(msgBytes []byte, topic string, producer sarama.SyncProducer) error {

	sendMsg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(msgBytes),
	}
	_, _, err := producer.SendMessage(&sendMsg)

	if err != nil {
		return err
	}

	return nil
}

func ProducerDeliverMessages(msgBytes [][]byte, topic string, producer sarama.SyncProducer) error {
	var sendMsgs []*sarama.ProducerMessage
	for _, msgByte := range msgBytes {
		sendMsg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.ByteEncoder(msgByte),
		}
		sendMsgs = append(sendMsgs, sendMsg)
	}
	err := producer.SendMessages(sendMsgs)
	if err != nil {
		return err
	}
	return nil

}

// SendToKafka : handles sending message to kafka
//func SendToKafka(msg KafkaMsg, kafkaState KafkaState, cdc *codec.LegacyAmino) []byte {
//	Error := ProducerDeliverMessage(msg, "Topic", kafkaState.Producer, cdc)
//	if Error != nil {
//		jsonResponse, Error := cdc.MarshalJSON(struct {
//			Response string `json:"response"`
//		}{Response: "Something is up with kafka server, restart rest and kafka."})
//		if Error != nil {
//			panic(Error)
//		}
//
//		SetTicketIDtoDB(msg.TicketID, kafkaState.KafkaDB, cdc, jsonResponse)
//	} else {
//		jsonResponse, err := cdc.MarshalJSON(struct {
//			Error string `json:"error"`
//		}{Error: "Request in process, wait and try after some time"})
//		if err != nil {
//			panic(err)
//		}
//		SetTicketIDtoDB(msg.TicketID, kafkaState.KafkaDB, cdc, jsonResponse)
//	}
//
//	jsonResponse, Error := cdc.MarshalJSON(struct {
//		TicketID Ticket `json:"ticketID"`
//	}{TicketID: msg.TicketID})
//
//	if Error != nil {
//		panic(Error)
//	}
//
//	return jsonResponse
//}
