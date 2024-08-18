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

func (uc *UserUsecaseImpl) AccountDecision() {
	log.Info().Msg("Starting account check decision listener...")
	reader := uc.kafka.NewConsumer("topic-user-decision")
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

		var password_valid string
		var status_valid string
		if err == nil {
			if request.Body.(map[string]interface{})["password_valid"] != nil {
				password_valid = request.Body.(map[string]interface{})["password_valid"].(string)
			} else {
				err = fmt.Errorf("password_valid is not generated")
			}
			if request.Body.(map[string]interface{})["status_valid"] != nil {
				status_valid = request.Body.(map[string]interface{})["status_valid"].(string)
			} else {
				err = fmt.Errorf("status_valid is not generated")
			}
		}

		response := orchestrator.Message{
			Header: orchestrator.Header{
				TransactionId:       request.Header.TransactionId,
				TransactionDateTime: request.Header.TransactionDateTime,
				OrderType:           request.Header.OrderType,
				OrderService:        "user-decision",
				Retries:             request.Header.Retries,
				ResponseCode:        http.StatusOK,
				ResponseMessage:     http.StatusText(http.StatusOK),
			},
			Body: request.Body,
		}

		if err != nil && err != fmt.Errorf("password_valid is not generated") && err != fmt.Errorf("status_valid is not generated") {
			response.Header.ResponseCode = http.StatusInternalServerError
			response.Header.ResponseMessage = http.StatusText(http.StatusInternalServerError)
		} else if err == fmt.Errorf("password_valid is not generated") || err == fmt.Errorf("status_valid is not generated") || password_valid == "" || status_valid == "" {
			response.Header.ResponseCode = http.StatusBadRequest
			response.Header.ResponseMessage = http.StatusText(http.StatusBadRequest)
		} else if status_valid != "valid" || password_valid != "valid" {
			response.Header.ResponseCode = http.StatusContinue
			response.Header.ResponseMessage = http.StatusText(http.StatusContinue)
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
