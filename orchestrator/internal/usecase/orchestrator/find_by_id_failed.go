package orchestrator

import (
	"context"
	"orchestrator/internal/domain/orchestrator"
)

func (uc *OrchestratorUsecaseImpl) FindByIdFailed(ctx context.Context) (orchestrator.Message, error) {
	return uc.orchestratorRepository.GetByIdFailed(ctx)
}
