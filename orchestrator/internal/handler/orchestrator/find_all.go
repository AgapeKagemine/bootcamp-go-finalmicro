package orchestrator

import (
	"fmt"
	"net/http"
	"orchestrator/internal/domain/handler"
	"orchestrator/internal/domain/orchestrator"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (h *OrchestratorHandlerImpl) FindAll(c *gin.Context) {
	log.Info().Msg("transaction (orchestrator) find all handler")

	response := handler.Response{
		StstusCode: 0,
		Message:    "",
		Payload:    []orchestrator.Message{},
	}

	defer func() {
		c.JSON(response.StstusCode, response)
		c.Request.Body.Close()
		log.Info().Int("status_code", response.StstusCode).Msg(fmt.Sprintf("transaction (orchestrator) find all - %s", response.Message))
	}()

	trx, err := h.orchestratorUsecase.FindAll(c)
	if err != nil {
		response.StstusCode = http.StatusInternalServerError
		response.Message = err.Error()
		return
	}

	response.StstusCode = http.StatusOK
	response.Message = http.StatusText(http.StatusOK)
	response.Payload = trx
}
