package runconfig

import "time"

type KafkaConfig struct {
	// Denom : staking denom,
	Denom string
	// Brokers: List of brokers to run kafka cluster
	Brokers      []string
	ToEth        TopicConsumer
	ToTendermint TopicConsumer
	// Time for each unbonding transactions 3 days => input nano-seconds 259200000000000
	EthUnbondCycleTime time.Duration
}

type TopicConsumer struct {
	BatchSize int
}

func NewKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Denom:   "stake",
		Brokers: []string{"localhost:9092"},
		ToEth: TopicConsumer{
			BatchSize: 2,
		},
		ToTendermint: TopicConsumer{
			BatchSize: 3,
		},
		EthUnbondCycleTime: time.Duration(259200000000000),
	}
}
