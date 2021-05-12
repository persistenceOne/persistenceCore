package runConfig

type KafkaConfig struct {
	// Brokers: List of brokers to run kafka cluster
	Brokers      []string
	ToEth        TopicConsumer
	ToTendermint TopicConsumer
}

type TopicConsumer struct {
	BatchSize int
}

func NewKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers: []string{"localhost:9092"},
		ToEth: TopicConsumer{
			BatchSize: 4,
		},
		ToTendermint: TopicConsumer{
			BatchSize: 2,
		},
	}
}
