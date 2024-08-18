package orchestrator

import (
	"context"
	orchestratorDomain "orchestrator/internal/domain/orchestrator"
	"orchestrator/internal/provider/messaging/kafka"
	orchestratorUsecase "orchestrator/internal/usecase/orchestrator"

	"github.com/gin-gonic/gin"
)

type OrchestratorHandler interface {
	FindAllFailed(*gin.Context)
	FindAll(*gin.Context)
	Retry(*gin.Context)
}

type OrchestratorHandlerImpl struct {
	orchestratorUsecase orchestratorUsecase.OrchestratorUsecase[context.Context, orchestratorDomain.Message, error]
	orchestratorKafka   kafka.OrchestratorKafka
}

func NewOrchestratorHandler(orchestratorUsecase orchestratorUsecase.OrchestratorUsecase[context.Context, orchestratorDomain.Message, error], orchestratorKafka kafka.OrchestratorKafka) OrchestratorHandler {
	return &OrchestratorHandlerImpl{
		orchestratorUsecase: orchestratorUsecase,
		orchestratorKafka:   orchestratorKafka,
	}
}
