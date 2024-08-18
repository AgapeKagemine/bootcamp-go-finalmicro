package orchestrator

import (
	"context"
	"orchestrator/internal/domain/orchestrator"
)

func (uc *OrchestratorUsecaseImpl) FindAll(ctx context.Context) ([]orchestrator.Message, error) {
	return uc.orchestratorRepository.GetAll(ctx)
}
