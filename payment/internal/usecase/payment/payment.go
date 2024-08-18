package payment

import "payment/internal/provider/messaging/kafka"

type PaymentUsecase interface {
	Start()
}

type PaymentUsecaseImpl struct {
	kafka kafka.OrchestratorKafka
}

func NewPaymentUsecase(kafka kafka.OrchestratorKafka) PaymentUsecase {
	return &PaymentUsecaseImpl{
		kafka: kafka,
	}
}
