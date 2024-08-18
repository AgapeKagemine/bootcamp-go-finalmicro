package orchestrator

import (
	"context"
	domain "orchestrator/internal/domain/orchestrator"
	"orchestrator/internal/provider/messaging/kafka"
	orchestratorRepo "orchestrator/internal/repository/orchestrator"

	"github.com/gin-gonic/gin"
)

type OrchestratorHandler interface {
	FindAllFailed(*gin.Context)
	FindAll(*gin.Context)
	Retry(*gin.Context)
}

type OrchestratorHandlerImpl struct {
	orchestratorRepo  orchestratorRepo.OrchestratorRepository[context.Context, domain.Message, error]
	orchestratorKafka kafka.OrchestratorKafka
}

func NewOrchestratorHandler(orchestratorRepo orchestratorRepo.OrchestratorRepository[context.Context, domain.Message, error], orchestratorKafka kafka.OrchestratorKafka) OrchestratorHandler {
	return &OrchestratorHandlerImpl{
		orchestratorRepo:  orchestratorRepo,
		orchestratorKafka: orchestratorKafka,
	}
}
