package kafka

import "github.com/segmentio/kafka-go"

type OrchestratorKafka interface {
	NewConsumer(string) *kafka.Reader
	NewProducer(string) *kafka.Writer
}

type OrchestratorKafkaImpl struct {
	broker string
}

const kafkaBroker = "localhost:29092"

func NewOrchestratorKafka() OrchestratorKafka {
	return &OrchestratorKafkaImpl{
		broker: kafkaBroker,
	}
}