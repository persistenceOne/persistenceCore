package kafka

import "github.com/Shopify/sarama"

func Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_7_0_0                 // hardcoded
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 3                    // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	return config
}
