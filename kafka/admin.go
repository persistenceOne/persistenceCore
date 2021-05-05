/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"github.com/Shopify/sarama"
)

// KafkaAdmin : is admin to create topics
func KafkaAdmin(kafkaPorts []string) sarama.ClusterAdmin {
	config := sarama.NewConfig()
	config.Version = sarama.V2_7_0_0 // hardcoded

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
