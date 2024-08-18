package routes

import (
	"orchestrator/internal/handler/orchestrator"

	"github.com/gin-gonic/gin"
)

func (r *Routes) TransactionRoutes(rg *gin.RouterGroup, h orchestrator.OrchestratorHandler) {
	transaction := rg.Group("/transaction")
	transaction.GET("/", h.FindAll)
	transaction.GET("/failed", h.FindAllFailed)
	transaction.PATCH("/retry", h.Retry)
}
