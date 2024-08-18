package kafka

import "github.com/segmentio/kafka-go"

func (ok *OrchestratorKafkaImpl) NewConsumer(topic string) *kafka.Reader {
	config := kafka.ReaderConfig{
		Brokers:     []string{ok.broker},
		Topic:       topic,
		GroupID:     "orchestrator-group",
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: kafka.LastOffset,
	}

	return kafka.NewReader(config)
}
