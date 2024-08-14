package routes

import (
	"order/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewRoutes() *gin.Engine {
	route := gin.New()

	api := route.Group("/api")
	api.POST("/order", handler.Order)

	return route
}
