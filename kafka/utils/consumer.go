/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"github.com/Shopify/sarama"
)

func NewConsumerGroup(kafkaPorts []string, groupID string, config *sarama.Config) sarama.ConsumerGroup {
	consumerGroup, Error := sarama.NewConsumerGroup(kafkaPorts, groupID, config)
	if Error != nil {
		panic(Error)
	}
	return consumerGroup
}
