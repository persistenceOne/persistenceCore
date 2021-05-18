/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"github.com/Shopify/sarama"
)

// KafkaAdmin : is admin to create topics
func KafkaAdmin(kafkaPorts []string, config *sarama.Config) sarama.ClusterAdmin {
	admin, Error := sarama.NewClusterAdmin(kafkaPorts, config)
	if Error != nil {
		panic(Error)
	}

	return admin
}

// TopicsInit : is needed to initialise topics
func TopicsInit(admin sarama.ClusterAdmin, topic string) {
	err := admin.CreateTopic(topic, &topicDetail, true)
	if err != nil {
		panic(err)
	}
}
