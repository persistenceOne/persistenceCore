package runConfig

type KafkaConfig struct {
	// Brokers: List of brokers to run kafka cluster
	Brokers        []string
	MsgSendForward TopicConsumer
	MsgSendRevert  TopicConsumer
}

type TopicConsumer struct {
	BatchSize int
}

func NewKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers: []string{"localhost:9092"},
		MsgSendForward: TopicConsumer{
			BatchSize: 2,
		},
		MsgSendRevert: TopicConsumer{
			BatchSize: 1,
		},
	}
}
