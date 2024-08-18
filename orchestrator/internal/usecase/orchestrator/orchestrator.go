package orchestrator

import (
	"context"
	orchestratorDomain "orchestrator/internal/domain/orchestrator"
	internalKafka "orchestrator/internal/provider/messaging/kafka"
	orchestratorRepository "orchestrator/internal/repository/orchestrator"
)

type OrchestratorUsecase[C context.Context, T orchestratorDomain.Message, E error] interface {
	FindAll(C) ([]T, E)
	FindAllFailed(C) ([]T, E)
	FindByIdFailed(C) (T, E)
	Start()
	ConsumeFailedTransaction()
}

type OrchestratorUsecaseImpl struct {
	orchestratorRepository orchestratorRepository.OrchestratorRepository[context.Context, orchestratorDomain.Message, error]
	orchestratorKafka      internalKafka.OrchestratorKafka
}

func NewOrchestratorUsecase(orchestratorRepository orchestratorRepository.OrchestratorRepository[context.Context, orchestratorDomain.Message, error], orchestratorKafka internalKafka.OrchestratorKafka) OrchestratorUsecase[context.Context, orchestratorDomain.Message, error] {
	return &OrchestratorUsecaseImpl{
		orchestratorRepository: orchestratorRepository,
		orchestratorKafka:      orchestratorKafka,
	}
}
