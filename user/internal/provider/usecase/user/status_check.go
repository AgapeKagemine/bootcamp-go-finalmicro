package user

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"user/internal/domain/orchestrator"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

func (uc *UserUsecaseImpl) StatusCheck() {
	log.Info().Msg("Starting account status check listener...")
	reader := uc.kafka.NewConsumer("topic-user-check-status")
	defer reader.Close()

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Error reading message")
			continue
		}

		var request orchestrator.Message
		err = json.Unmarshal(msg.Value, &request)
		if err != nil {
			log.Error().Err(err).Msg("Error parsing message")
		}

		log.Info().Msg(fmt.Sprintf("Message received: %s\n", string(msg.Value)))

		var userId string
		if err == nil {
			if request.Body.(map[string]interface{})["user_id"] != nil {
				userId = request.Body.(map[string]interface{})["user_id"].(string)
			} else {
				err = fmt.Errorf("user_id is required")
			}
		}

		response := orchestrator.Message{
			Header: orchestrator.Header{
				TransactionId:       request.Header.TransactionId,
				TransactionDateTime: request.Header.TransactionDateTime,
				OrderType:           request.Header.OrderType,
				OrderService:        "user-check-status",
				Retries:             request.Header.Retries,
				ResponseCode:        http.StatusOK,
				ResponseMessage:     http.StatusText(http.StatusOK),
			},
			Body: request.Body,
		}

		if err != nil && err != fmt.Errorf("user_id is required") {
			response.Header.ResponseCode = http.StatusInternalServerError
			response.Header.ResponseMessage = http.StatusText(http.StatusInternalServerError)
		} else if userId == "" || err == fmt.Errorf("user_id is required") {
			response.Header.ResponseCode = http.StatusBadRequest
			response.Header.ResponseMessage = http.StatusText(http.StatusBadRequest)
		} else {
			if userId != "u-0001" {
				response.Body.(map[string]interface{})["status_valid"] = "invalid"
			} else {
				response.Body.(map[string]interface{})["status_valid"] = "valid"
			}
		}

		payload, err := json.Marshal(response)
		if err != nil {
			log.Error().Err(err).Msg("Error marshalling response")
		}

		producer := uc.kafka.NewProducer("topic-orchestrator")
		err = producer.WriteMessages(context.Background(),
			kafka.Message{
				Value: payload,
			},
		)

		if err != nil {
			log.Error().Err(err).Msg("Error writing message to kafka")
		}
		producer.Close()
	}
}
