package orchestrator

import (
	"context"
	"orchestrator/internal/domain/orchestrator"
)

func (h *OrchestratorUsecaseImpl) FindAllFailed(ctx context.Context) ([]orchestrator.Message, error) {
	return h.orchestratorRepository.GetAllFailed(ctx)
}
