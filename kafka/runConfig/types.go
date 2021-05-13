package runConfig

import "time"

type KafkaConfig struct {
	// Denom : staking denom,
	Denom string
	// Brokers: List of brokers to run kafka cluster
	Brokers      []string
	ToEth        TopicConsumer
	ToTendermint TopicConsumer
	// unbond pools correct value will be 8 (7+1) => if 7 max tuple of delegator/validator
	UnbondPools int
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
			BatchSize: 4,
		},
		ToTendermint: TopicConsumer{
			BatchSize: 2,
		},
		UnbondPools:        8,
		EthUnbondCycleTime: time.Duration(259200000000000),
	}
}
