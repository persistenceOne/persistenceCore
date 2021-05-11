/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package rest

// KafkaConsumerMessages : messages to consume 5 second delay
//func KafkaConsumerMessages(codec *codec.LegacyAmino, kafkaState kafka.KafkaState) {
//	quit := make(chan bool)
//
//	var ticketIDList []kafka.Ticket
//
//	var msgList []sdkTypes.Msg
//
//	go func() {
//		for {
//			select {
//			case <-quit:
//				return
//			default:
//				kafkaMsg := kafka.KafkaTopicConsumer("Topic", kafkaState.ConsumerGroup, codec)
//				if kafkaMsg.Msg != nil {
//					ticketIDList = append(ticketIDList, kafkaMsg.TicketID)
//					msgList = append(msgList, kafkaMsg.Msg)
//				}
//			}
//		}
//	}()
//
//	time.Sleep(kafka.SleepTimer)
//	quit <- true
//
//	if len(msgList) == 0 {
//		return
//	}
//
//	output, err := squash()
//	if err != nil {
//		jsonError, e := codec.MarshalJSON(struct {
//			Error string `json:"error"`
//		}{Error: err.Error()})
//		if e != nil {
//			panic(err)
//		}
//
//		for _, ticketID := range ticketIDList {
//			kafka.AddResponseToDB(ticketID, jsonError, kafkaState.KafkaDB, codec)
//		}
//
//		return
//	}
//
//	for _, ticketID := range ticketIDList {
//		kafka.AddResponseToDB(ticketID, output, kafkaState.KafkaDB, codec)
//	}
//}
