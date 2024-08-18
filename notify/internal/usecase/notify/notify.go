package notify

import "notify/internal/provider/messaging/kafka"

type NotifyUsecase struct {
	kafka kafka.OrchestratorKafka
}

func NewNotifyUsecase(kafka kafka.OrchestratorKafka) *NotifyUsecase {
	return &NotifyUsecase{
		kafka: kafka,
	}
}
