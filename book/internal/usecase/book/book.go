package book

import (
	"book/internal/provider/messaging/kafka"
)

type BookUsecase struct {
	kafka kafka.OrchestratorKafka
}

func NewBookUsecase(kafka kafka.OrchestratorKafka) *BookUsecase {
	return &BookUsecase{
		kafka: kafka,
	}
}
