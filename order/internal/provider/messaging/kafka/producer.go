package kafka

import (
	"github.com/segmentio/kafka-go"
)

const broker = "localhost:29092"

func NewProducer(topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})
}
