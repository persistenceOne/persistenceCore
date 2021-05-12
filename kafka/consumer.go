/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"errors"
	"github.com/Shopify/sarama"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	KafkaConfig runConfig.KafkaConfig
}

var _ sarama.ConsumerGroupHandler = MsgHandler{}

func (m MsgHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (m MsgHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (m MsgHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	switch claim.Topic() {
	case ToEth:
		err := HandleTopicMsgs(session, claim, m.KafkaConfig.ToEth.BatchSize, SendBatchToEth)
		if err != nil {
			log.Printf("failed batch and handle for topic: %v", ToEth)
			return err
		}
	case ToTendermint:
		err := HandleTopicMsgs(session, claim, m.KafkaConfig.ToTendermint.BatchSize, SendBatchToTendermint)
		if err != nil {
			log.Printf("failed batch and handle for topic: %v", ToTendermint)
			return err
		}
	}
	return nil
}

// Handlers of message types
func HandleTopicMsgs(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, batchSize int, handle func([]sarama.ConsumerMessage) error) error {
	msgs := make([]sarama.ConsumerMessage, 0, batchSize)
	for {
		kafkaMsg := <-claim.Messages()
		if kafkaMsg == nil {
			return errors.New("kafka returned nil message")
		}
		log.Printf("Message topic:%q partition:%d offset:%d\n", kafkaMsg.Topic, kafkaMsg.Partition, kafkaMsg.Offset)

		ok, err := BatchAndHandle(&msgs, *kafkaMsg, handle)
		if ok && err == nil {
			session.MarkMessage(kafkaMsg, "")
			return nil
		}
		if err != nil {
			return err
		}
	}
}

// BatchAndHandle :
func BatchAndHandle(kafkaMsgs *[]sarama.ConsumerMessage, kafkaMsg sarama.ConsumerMessage,
	handle func([]sarama.ConsumerMessage) error) (bool, error) {
	*kafkaMsgs = append(*kafkaMsgs, kafkaMsg)
	if len(*kafkaMsgs) == cap(*kafkaMsgs) {
		err := handle(*kafkaMsgs)
		if err != nil {
			return false, err
		}
		*kafkaMsgs = (*kafkaMsgs)[:0]
		return true, nil
	}
	return false, nil
}

func ConvertKafkaMsgsToSDKMsg(kafkaMsgs []sarama.ConsumerMessage) ([]sdk.Msg, error) {
	var msgs []sdk.Msg
	for _, kafkaMsg := range kafkaMsgs {
		var msg sdk.Msg
		err := proto.Unmarshal(kafkaMsg.Value, msg)
		if err != nil {
			return nil, errors.New("error unmarshalling proto")
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// SendBatchToEth : Handling of msgSend
func SendBatchToEth(kafkaMsgs []sarama.ConsumerMessage) error {
	msgs, err := ConvertKafkaMsgsToSDKMsg(kafkaMsgs)
	if err != nil {
		return err
	}
	log.Printf("batched messages: %v", msgs)
	// do more with msgs.
	return nil
}

// SendBatchToTendermint :
func SendBatchToTendermint(kafkaMsgs []sarama.ConsumerMessage) error {
	msgs, err := ConvertKafkaMsgsToSDKMsg(kafkaMsgs)
	if err != nil {
		return err
	}
	log.Printf("batched messages: %v", msgs)
	// do more with messages.
	return nil
}
