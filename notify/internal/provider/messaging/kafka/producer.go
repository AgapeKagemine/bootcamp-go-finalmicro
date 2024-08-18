package kafka

import (
	"github.com/segmentio/kafka-go"
)

func (ok *OrchestratorKafkaImpl) NewProducer(topic string) *kafka.Writer {
	config := kafka.WriterConfig{
		Brokers: []string{ok.broker},
		Topic:   topic,
	}
	return kafka.NewWriter(config)
}
