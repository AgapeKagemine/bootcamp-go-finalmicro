package kafka

import "github.com/segmentio/kafka-go"

type OrchestratorKafka interface {
	NewConsumer(string) *kafka.Reader
	NewProducer(string) *kafka.Writer
	GetMaxRetries() int
}

type OrchestratorKafkaImpl struct {
	MaxRetries int
	broker     string
}

const (
	kafkaBroker = "localhost:29092"
)

func NewOrchestratorKafka(retries int) OrchestratorKafka {
	return &OrchestratorKafkaImpl{
		MaxRetries: retries,
		broker:     kafkaBroker,
	}
}
