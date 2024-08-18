package user

import "user/internal/provider/messaging/kafka"

type UserUsecase interface {
	Start()
	AccountDecision()
	ValidatePassword()
	StatusCheck()
}

type UserUsecaseImpl struct {
	kafka kafka.OrchestratorKafka
}

func NewUserUsecase(kafka kafka.OrchestratorKafka) UserUsecase {
	return &UserUsecaseImpl{
		kafka: kafka,
	}
}
