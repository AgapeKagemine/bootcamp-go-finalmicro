package routes

import (
	"orchestrator/internal/handler/orchestrator"

	"github.com/gin-gonic/gin"
)

func NewRoutes(h orchestrator.OrchestratorHandler) *gin.Engine {
	route := gin.New()

	api := route.Group("/api")
	transaction := api.Group("/transaction")
	transaction.GET("/", h.FindAll)
	transaction.GET("/failed", h.FindAllFailed)
	transaction.PATCH("/retry", h.Retry)

	return route
}
