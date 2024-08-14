package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"order/internal/domain/orchestrator"
	internalKafka "order/internal/provider/messaging/kafka"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

const (
	orderService = "START"
)

func Order(ctx *gin.Context) {
	rawRequest := make(map[string]interface{})
	err := ctx.ShouldBindJSON(&rawRequest)
	if err != nil {
		log.Error().Err(err).Msg("Invalid request payload")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	writer := internalKafka.NewProducer("orchestrator-topic")
	defer writer.Close()

	request := make(map[string]interface{})

	for key, value := range rawRequest {
		if key == "order_type" || key == "order_service" {
			continue
		}
		request[key] = value
	}

	trxId, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate transaction id")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate transaction id"})
		return
	}

	if request["order_type"] == nil {
		log.Error().Err(err).Msg("Invalid request payload")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid request payload"})
		return
	}

	orchestRequest := orchestrator.Message{
		Header: orchestrator.Header{
			TransactionId:       trxId.String(),
			TransactionDateTime: time.Now().Local(),
			OrderType:           rawRequest["order_type"].(string),
			OrderService:        orderService,
			Retries:             0,
			ResponseCode:        http.StatusCreated,
			ResponseMessage:     http.StatusText(http.StatusCreated),
		},
		Body: request,
	}

	payload, err := json.Marshal(orchestRequest)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal request")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
		return
	}

	message := kafka.Message{
		Value: payload,
	}

	err = writer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to produce message")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to produce Kafka message"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Order placed successfully", "Order": orchestRequest})
}
