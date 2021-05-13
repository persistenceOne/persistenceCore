/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"errors"
	"github.com/Shopify/sarama"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
	case EthUnbond:
		err := HandleEthUnbond(session, claim, m.KafkaConfig.Brokers, m.KafkaConfig.Denom)
		if err != nil {
			log.Printf("failed to handle EthUnbonding for topic: %v", EthUnbond)
			return err
		}
	case UnbondPool:
		err := HandleUnbondPool(session, claim, m.KafkaConfig.Brokers)
		if err != nil {
			log.Printf("failed to handle unbond pool for topic: %v", UnbondPool)
		}
	}
	return nil
}

func HandleEthUnbond(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, brokers []string, denom string) error {
	config := Config()
	producer := NewProducer(brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v", UnbondPool)
		}
	}()
	var sum = sdk.NewInt(0)
	for kafkaMsg := range claim.Messages() {
		if kafkaMsg == nil {
			return errors.New("kafka returned nil message")
		}
		var msg sdk.Msg
		err := proto.Unmarshal(kafkaMsg.Value, msg)
		if err != nil {
			log.Printf("proto failed to unmarshal")
		}
		switch txMsg := msg.(type) {
		case *bankTypes.MsgSend:
			// TODO is denom fixed?
			sum = sum.Add(txMsg.Amount.AmountOf(denom))
		default:
			log.Printf("Unexpected type found in topic: %v", EthUnbond)
		}
		err = ProducerDeliverMessage(kafkaMsg.Value, UnbondPool, producer)
		if err != nil {
			log.Printf("failed to produce message from topic %v to %v", EthUnbond, UnbondPool)
		}
		session.MarkMessage(kafkaMsg, "")
	}
	/*
		// Make a unbond msg and send TODO pick delegator and validator addresses
		unbondMsg := &stakingTypes.MsgUndelegate{
			DelegatorAddress: "",
			ValidatorAddress: "",
			Amount:           sdk.Coin{
				Denom:  denom,
				Amount: sum,
			},
		}
		msgBytes, err := proto.Marshal(sdk.Msg(unbondMsg))
		if err!= nil {
			return err
		}
		err = ProducerDeliverMessage(msgBytes, ToTendermint, producer)
		if err != nil {
			log.Printf("failed to produce message from topic %v to %v", EthUnbond, ToTendermint)
			return err
		}
	*/

	return nil
}

func HandleUnbondPool(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, brokers []string) error {
	config := Config()
	producer := NewProducer(brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v", UnbondPool)
		}
	}()
	if len(claim.Messages()) < 8 {
		// fill in till there are 8 nil msgs,
		// This is only for starting, it will get tricky if application has downtime of more than a cycle.(3days)
		// think about starting a consumer to read the messages and quit.
		for i := len(claim.Messages()); i < 8; i++ {
			sdkMsg := sdk.Msg(nil)
			msgBytes, err := proto.Marshal(sdkMsg)
			if err != nil {
				return nil
			}
			err = ProducerDeliverMessage(msgBytes, UnbondPool, producer)
			if err != nil {
				log.Printf("failed to produce message from topic %v to %v", UnbondPool, UnbondPool)
			}
		}
	} else {
		//consume

		for kafkaMsg := range claim.Messages() {
			if kafkaMsg == nil {
				return errors.New("kafka returned nil message in topic: " + UnbondPool)
			}
			var msg sdk.Msg
			err := proto.Unmarshal(kafkaMsg.Value, msg)
			if err != nil {
				return errors.New("error unmarshalling proto")
			}
			if msg == nil {
				return nil
			}
			err = ProducerDeliverMessage(kafkaMsg.Value, ToTendermint, producer)
			if err != nil {
				log.Printf("failed to produce message from topic %v to %v", UnbondPool, ToTendermint)
			}
			session.MarkMessage(kafkaMsg, "")
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
