/*
 Copyright [2019] - [2021], PERSISTENCE TECHNOLOGIES PTE. LTD. and the persistenceSDK contributors
 SPDX-License-Identifier: Apache-2.0
*/

package utils

import (
	"github.com/Shopify/sarama"
	"os"
)

// These are the config parameters for running kafka admins and producers and consumers. Declared very minimal
var replicaAssignment = map[int32][]int32{}
var configEntries = map[string]*string{}

// DefaultKafkaHome : is the home path
var DefaultKafkaHome = os.ExpandEnv("$HOME/.kafka")

var FlagKafkaHome = "kafka-home"

// topicDetail : configs only required for admin to create topics if not present.
var topicDetail = sarama.TopicDetail{
	NumPartitions:     1,
	ReplicationFactor: 1,
	ReplicaAssignment: replicaAssignment,
	ConfigEntries:     configEntries,
}

// Consumer groups
const GroupToEth = "group-to-ethereum"
const GroupToTendermint = "group-to-tendermint"
const GroupEthUnbond = "group-ethereum-unbond"
const GroupMsgSend = "group-msg-send"
const GroupMsgDelegate = "group-msg-delegate"
const GroupMsgUnbond = "group-msg-unbond"

var Groups = []string{GroupToEth, GroupToTendermint, GroupEthUnbond,
	GroupMsgSend, GroupMsgDelegate, GroupMsgUnbond,
}

//Topics
const ToEth = "to-ethereum"
const ToTendermint = "to-tendermint"
const MsgSend = "msg-send"          //priority3
const MsgDelegate = "msg-delegate"  //priority2
const MsgUnbond = "msg-unbond"      //priority1
const EthUnbond = "ethereum-unbond" //flushes every 3 days

// Topics : is list of topics
var Topics = []string{
	ToEth, ToTendermint, EthUnbond,
	MsgSend, MsgDelegate, MsgUnbond,
}
