package routes

import (
	"orchestrator/internal/handler/orchestrator"

	"github.com/gin-gonic/gin"
)

type Routes struct {
	Routes *gin.Engine
}

func NewRoutes(h orchestrator.OrchestratorHandler) *Routes {
	route := &Routes{
		Routes: gin.New(),
	}

	api := route.Routes.Group("/api")
	route.TransactionRoutes(api, h)

	return route
}
