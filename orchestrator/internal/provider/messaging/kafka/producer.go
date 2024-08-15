package kafka

import (
	"github.com/segmentio/kafka-go"
)

func NewProducer(topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:29092"},
		Topic:   topic,
	})
}
