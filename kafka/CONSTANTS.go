/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package kafka

import (
	"os"
	"time"

	"github.com/Shopify/sarama"
)

// SleepTimer : the time the kafka messages are to be taken in
const SleepTimer = time.Duration(1000000000)

// SleepRoutine : the time the kafka messages are to be taken in
const SleepRoutine = time.Duration(2500000000)

// These are the config parameters for running kafka admins and producers and consumers. Declared very minimal
var replicaAssignment = map[int32][]int32{}
var configEntries = map[string]*string{}

// DefaultCLIHome : is the home path
var DefaultCLIHome = os.ExpandEnv("$HOME/.kafka")

const partition = int32(0)
const offset = int64(0)

// topicDetail : configs
var topicDetail = sarama.TopicDetail{
	NumPartitions:     1,
	ReplicationFactor: 1,
	ReplicaAssignment: replicaAssignment,
	ConfigEntries:     configEntries,
}

// Topics : is list of topics
var Topics = []string{
	"Topic",
}
