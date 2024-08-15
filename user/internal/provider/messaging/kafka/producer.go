package kafka

import (
	"github.com/segmentio/kafka-go"
)

func NewProducer(topic string) *kafka.Writer {
	config := kafka.WriterConfig{
		Brokers: []string{"localhost:29092"},
		Topic:   topic,
	}
	return kafka.NewWriter(config)
}
