package kafka

import (
	"github.com/segmentio/kafka-go"
)

func NewConsumer(topic string) *kafka.Reader {
	config := kafka.ReaderConfig{
		Brokers:     []string{"localhost:29092"},
		Topic:       topic,
		GroupID:     "book-consumer-group",
		MinBytes:    10e3,
		MaxBytes:    10e6,
		StartOffset: kafka.LastOffset,
	}
	return kafka.NewReader(config)
}
