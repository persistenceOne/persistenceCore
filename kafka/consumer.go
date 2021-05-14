/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"errors"
	"github.com/Shopify/sarama"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankTypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
	ProtoCodec  *codec.ProtoCodec
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
		err := m.HandleTopicMsgs(session, claim, m.KafkaConfig.ToEth.BatchSize, SendBatchToEth)
		if err != nil {
			log.Printf("failed batch and handle for topic: %v with error %v", ToEth, err)
			return err
		}
	case ToTendermint:
		err := m.HandleTopicMsgs(session, claim, m.KafkaConfig.ToTendermint.BatchSize, SendBatchToTendermint)
		if err != nil {
			log.Printf("failed batch and handle for topic: %v with error %v", ToTendermint, err)
			return err
		}
	case EthUnbond:
		err := m.HandleEthUnbond(session, claim)
		if err != nil {
			log.Printf("failed to handle EthUnbonding for topic: %v", EthUnbond)
			return err
		}
	}
	return nil
}

func (m MsgHandler) HandleEthUnbond(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	config := Config()
	producer := NewProducer(m.KafkaConfig.Brokers, config)
	defer func() {
		err := producer.Close()
		if err != nil {
			log.Printf("failed to close producer in topic: %v", EthUnbond)
		}
	}()
	var sum = sdk.NewInt(0)
	for kafkaMsg := range claim.Messages() {
		if kafkaMsg == nil {
			return errors.New("kafka returned nil message")
		}
		var msg sdk.Msg
		err := m.ProtoCodec.UnmarshalInterface(kafkaMsg.Value, &msg)
		if err != nil {
			log.Printf("proto failed to unmarshal")
		}
		switch txMsg := msg.(type) {
		case *bankTypes.MsgSend:
			sum = sum.Add(txMsg.Amount.AmountOf(m.KafkaConfig.Denom))
		default:
			log.Printf("Unexpected type found in topic: %v", EthUnbond)
		}
		// should we mark offset after summing all? more like a defer function, will be tricky if kafkaMsg is nil
		session.MarkMessage(kafkaMsg, "")
	}

	/*
		// Make a unbond msg and send TODO pick delegator and validator addresses
		unbondMsg := &stakingTypes.MsgUndelegate{
			DelegatorAddress: "",
			ValidatorAddress: "",
			Amount:           sdk.Coin{
				Denom:  m.KafkaConfig.Denom,
				Amount: sum,
			},
		}
		msgBytes, err := m.ProtoCodec.MarshalInterface(sdk.Msg(unbondMsg))
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

// Handlers of message types
func (m MsgHandler) HandleTopicMsgs(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim, batchSize int,
	handle func([]sarama.ConsumerMessage, *codec.ProtoCodec) error) error {
	msgs := make([]sarama.ConsumerMessage, 0, batchSize)
	for {
		kafkaMsg := <-claim.Messages()
		if kafkaMsg == nil {
			return errors.New("kafka returned nil message")
		}
		log.Printf("Message topic:%q partition:%d offset:%d\n", kafkaMsg.Topic, kafkaMsg.Partition, kafkaMsg.Offset)

		ok, err := BatchAndHandle(&msgs, *kafkaMsg, m.ProtoCodec, handle)
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
	protoCodec *codec.ProtoCodec, handle func([]sarama.ConsumerMessage, *codec.ProtoCodec) error) (bool, error) {
	*kafkaMsgs = append(*kafkaMsgs, kafkaMsg)
	if len(*kafkaMsgs) == cap(*kafkaMsgs) {
		err := handle(*kafkaMsgs, protoCodec)
		if err != nil {
			return false, err
		}
		*kafkaMsgs = (*kafkaMsgs)[:0]
		return true, nil
	}
	return false, nil
}

func ConvertKafkaMsgsToSDKMsg(kafkaMsgs []sarama.ConsumerMessage, protoCodec *codec.ProtoCodec) ([]sdk.Msg, error) {
	var msgs []sdk.Msg
	for _, kafkaMsg := range kafkaMsgs {
		var msg sdk.Msg
		err := protoCodec.UnmarshalInterface(kafkaMsg.Value, &msg)
		if err != nil {
			return nil, err
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}

// SendBatchToEth : Handling of msgSend
func SendBatchToEth(kafkaMsgs []sarama.ConsumerMessage, protoCodec *codec.ProtoCodec) error {
	msgs, err := ConvertKafkaMsgsToSDKMsg(kafkaMsgs, protoCodec)
	if err != nil {
		return err
	}
	log.Printf("batched messages to send to ETH: %v", msgs)
	// TODO: do more with msgs.
	return nil
}

// SendBatchToTendermint :
func SendBatchToTendermint(kafkaMsgs []sarama.ConsumerMessage, protoCodec *codec.ProtoCodec) error {
	msgs, err := ConvertKafkaMsgsToSDKMsg(kafkaMsgs, protoCodec)
	if err != nil {
		return err
	}
	log.Printf("batched messages to send to Tendermint: %v", msgs)
	//TODO: do more with messages.
	return nil
}