package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"user/internal/domain/orchestrator"
	"user/internal/domain/user"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

const userService = "http://127.0.0.1:8090/user"

func (uc *UserUsecaseImpl) Start() {
	log.Info().Msg("Starting user event listener...")
	reader := uc.kafka.NewConsumer("topic-user-validate")
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

		userRequest := user.Request{}
		if err == nil {
			if request.Body.(map[string]interface{})["user_id"] != nil {
				userRequest.UserId = request.Body.(map[string]interface{})["user_id"].(string)
			} else {
				err = fmt.Errorf("user_id is required")
			}
		}

		var userResponse user.Response
		if err == nil {
			userPayload, err := json.Marshal(userRequest)
			if err != nil {
				log.Error().Err(err).Msg("Error marshalling user request")
			}

			res, err := http.Post(userService, "application/json", bytes.NewBuffer(userPayload))
			if err != nil {
				log.Error().Err(err).Msg("Error sending user request")
			}

			if err == nil {
				err = json.NewDecoder(res.Body).Decode(&userResponse)
				if err == io.EOF {
					err = nil
				}
			}
		}

		response := orchestrator.Message{
			Header: orchestrator.Header{
				TransactionId:       request.Header.TransactionId,
				TransactionDateTime: request.Header.TransactionDateTime,
				OrderType:           request.Header.OrderType,
				OrderService:        "user-validate",
				Retries:             request.Header.Retries,
				ResponseCode:        http.StatusOK,
				ResponseMessage:     http.StatusText(http.StatusOK),
			},
			Body: request.Body,
		}

		if err != nil && err != fmt.Errorf("user_id is required") {
			response.Header.ResponseCode = http.StatusInternalServerError
			response.Header.ResponseMessage = http.StatusText(http.StatusInternalServerError)
		} else if userResponse.Message != "valid" || err == fmt.Errorf("user_id is required") {
			response.Header.ResponseCode = http.StatusBadRequest
			response.Header.ResponseMessage = http.StatusText(http.StatusBadRequest)
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
