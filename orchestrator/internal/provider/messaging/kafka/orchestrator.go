package kafka

import (
	"context"
	orchestratorDomain "orchestrator/internal/domain/orchestrator"
	orchestratorRepository "orchestrator/internal/repository/orchestrator"
)

const (
	kafkaBroker = "localhost:29092"
)

type OrchestratorKafka struct {
	repo       orchestratorRepository.OrchestratorRepository[context.Context, orchestratorDomain.Message, error]
	maxRetries int
	broker     string
}

func NewOrchestratorKafka(repo orchestratorRepository.OrchestratorRepository[context.Context, orchestratorDomain.Message, error]) *OrchestratorKafka {
	return &OrchestratorKafka{
		repo:       repo,
		maxRetries: 3,
		broker:     kafkaBroker,
	}
}
