package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"user/internal/domain/orchestrator"
	"user/internal/domain/user"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

const userService = "http://127.0.0.1:8090/user"

func Start() {
	producer := NewProducer("topic-orchestrator")
	reader := NewConsumer("topic-user-validate")

	defer func() {
		producer.Close()
		reader.Close()
	}()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Error reading message")
			continue
		}

		var request orchestrator.Message
		if err := json.Unmarshal(m.Value, &request); err != nil {
			log.Error().Err(err).Msg("Error parsing message")
		}

		userRequest := user.Request{}

		for k, v := range request.Body.(map[string]interface{}) {
			if k == "user_id" {
				userRequest.UserId = v.(string)
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

		userPayload, err := json.Marshal(userRequest)
		if err != nil {
			log.Error().Err(err).Msg("Error marshalling user request")
		}

		res, err := http.Post(userService, "application/json", bytes.NewBuffer(userPayload))
		if err != nil {
			log.Error().Err(err).Msg("Error sending user request")
		}

		byteRes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Error().Err(err).Msg("Error reading user response")
		}

		var userResponse user.Response
		err = json.Unmarshal(byteRes, &userResponse)
		if err != nil {
			log.Error().Err(err).Msg("Error decoding user response")
		}

		if err != nil {
			response.Header.ResponseCode = http.StatusInternalServerError
			response.Header.ResponseMessage = http.StatusText(http.StatusInternalServerError)
		}

		if userResponse.Message != "valid" && err == nil {
			response.Header.ResponseCode = http.StatusBadRequest
			response.Header.ResponseMessage = http.StatusText(http.StatusBadRequest)
		}

		payload, err := json.Marshal(response)
		if err != nil {
			log.Error().Err(err).Msg("Error marshalling response")
		}

		// Produce the response message to topic_0
		err = producer.WriteMessages(context.Background(),
			kafka.Message{
				Value: payload,
			},
		)

		if err != nil {
			log.Error().Err(err).Msg("Error writing message to kafka")
			continue
		}
	}
}
