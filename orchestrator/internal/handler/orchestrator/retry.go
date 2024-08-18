package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"orchestrator/internal/domain"

	"orchestrator/internal/domain/handler"
	"orchestrator/internal/domain/orchestrator"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

func (h *OrchestratorHandlerImpl) Retry(c *gin.Context) {
	log.Info().Msg("transaction (orchestrator) retry/rollback handler")

	response := handler.Response{
		StstusCode: 0,
		Message:    "",
		Payload:    orchestrator.Message{},
	}

	defer func() {
		c.JSON(response.StstusCode, response)
		c.Request.Body.Close()
		log.Info().Int("status_code", response.StstusCode).Msg(fmt.Sprintf("transaction (orchestrator) retry/rollback handler - %s", response.Message))
	}()

	var requestTransaction orchestrator.Message
	err := c.ShouldBindJSON(&requestTransaction)
	if err != nil {
		response.StstusCode = http.StatusBadRequest
		response.Message = err.Error()
		return
	}

	response.Payload = requestTransaction

	ctx := context.WithValue(c.Request.Context(), domain.Key("request"), requestTransaction)

	trueTransaction, err := h.orchestratorUsecase.FindByIdFailed(ctx)
	if err != nil {
		response.StstusCode = http.StatusInternalServerError
		response.Message = err.Error()
		return
	}

	if validateRetry(requestTransaction, trueTransaction) {
		response.StstusCode = http.StatusBadRequest
		response.Message = "invalid request"
		return
	}

	requestTransaction.Header.TransactionDateTime = time.Now().Local()
	requestTransaction.Header.ResponseMessage = "Manual Retry"
	requestTransaction.Header.Retries = 0

	response.Payload = requestTransaction

	payload, err := json.Marshal(requestTransaction)
	if err != nil {
		response.StstusCode = http.StatusInternalServerError
		response.Message = err.Error()
		return
	}

	writer := h.orchestratorKafka.NewProducer("topic-orchestrator")
	err = writer.WriteMessages(context.Background(),
		kafka.Message{
			Value: payload,
		},
	)

	if err != nil {
		response.StstusCode = http.StatusInternalServerError
		response.Message = err.Error()
		return
	}

	response.StstusCode = http.StatusOK
	response.Message = http.StatusText(http.StatusOK)
}

// Check the field requested field for retry
// the only can be changed: order_type, retries, response_code, request_body
func validateRetry(request, truth orchestrator.Message) bool {
	if request.Header.TransactionId != truth.Header.TransactionId {
		return true
	}

	if request.Header.TransactionDateTime != truth.Header.TransactionDateTime {
		return true
	}

	if request.Header.OrderType != truth.Header.OrderType {
		return true
	}

	if request.Header.OrderService != truth.Header.OrderService {
		return true
	}

	if request.Header.ResponseCode != truth.Header.ResponseCode {
		return true
	}

	return false
}
