/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"errors"
	"github.com/Shopify/sarama"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/golang/protobuf/proto"
	"github.com/persistenceOne/persistenceCore/kafka/runConfig"
	"log"
)

func NewConsumerGroup(kafkaPorts []string, groupID string, config *sarama.Config) sarama.ConsumerGroup {
	consumerGroup, Error := sarama.NewConsumerGroup(kafkaPorts, groupID, config)
	if Error != nil {
		panic(Error)
	}
	return consumerGroup
}

type MsgHandler struct {
	runConfig.KafkaConfig
}

var _ sarama.ConsumerGroupHandler = MsgHandler{}

func (m MsgHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (m MsgHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (m MsgHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	msgSendForward := make([]sarama.ConsumerMessage, 0, m.MsgSendForward.BatchSize)
	msgSendRevert := make([]sarama.ConsumerMessage, 0, m.MsgSendRevert.BatchSize)
	for {
		kafkaMsg := <-claim.Messages()
		if kafkaMsg == nil {
			return errors.New("kafka returned nil message")
		}
		log.Printf("Message topic:%q partition:%d offset:%d\n", kafkaMsg.Topic, kafkaMsg.Partition, kafkaMsg.Offset)

		switch topic := kafkaMsg.Topic; topic {
		case MsgSendForward:
			ok := BatchAndHandle(&msgSendForward, *kafkaMsg, HandleMsgSendForward)
			if ok {
				session.MarkMessage(kafkaMsg, "")
			}
		case MsgSendRevert:
			ok := BatchAndHandle(&msgSendRevert, *kafkaMsg, HandleMsgSendRevert)
			if ok {
				session.MarkMessage(kafkaMsg, "")
			}
		}
	}
}

// Handlers of message types

// BatchMsgSendForward :
func BatchAndHandle(kafkaMsgs *[]sarama.ConsumerMessage, kafkaMsg sarama.ConsumerMessage,
	handle func([]sarama.ConsumerMessage) error) bool {
	*kafkaMsgs = append(*kafkaMsgs, kafkaMsg)
	if len(*kafkaMsgs) == cap(*kafkaMsgs) {
		err := handle(*kafkaMsgs)
		if err != nil {
			log.Printf("error in handling msgsendForward: %v", err)
			return false
		}
		*kafkaMsgs = (*kafkaMsgs)[:0]
		return true
	}
	return false
}

func ConvertKafkaMsgsToMsgSend(kafkaMsgs []sarama.ConsumerMessage) ([]banktypes.MsgSend, error) {
	var msgs []banktypes.MsgSend
	for _, kafkaMsg := range kafkaMsgs {
		var msg = banktypes.MsgSend{}
		err := proto.Unmarshal(kafkaMsg.Value, &msg)
		if err != nil {
			return nil, errors.New("error unmarshalling proto")
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// HandleMsgSendForward : Handling of msgSend
func HandleMsgSendForward(kafkaMsgs []sarama.ConsumerMessage) error {
	msgs, err := ConvertKafkaMsgsToMsgSend(kafkaMsgs)
	if err != nil {
		return err
	}
	log.Printf("batched messages: %v", msgs)
	// do more with msgs.
	return nil
}

func HandleMsgSendRevert(kafkaMsgs []sarama.ConsumerMessage) error {
	msgs, err := ConvertKafkaMsgsToMsgSend(kafkaMsgs)
	if err != nil {
		return err
	}
	log.Printf("batched messages: %v", msgs)

	return nil
}
