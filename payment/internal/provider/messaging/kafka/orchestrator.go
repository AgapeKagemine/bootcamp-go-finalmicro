package kafka

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"payment/internal/domain/orchestrator"
	"payment/internal/domain/payment"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

const paymentService = "http://127.0.0.1:8090/balance"

func Start() {
	producer := NewProducer("topic-orchestrator")
	reader := NewConsumer("topic-payment-deduct-balance")

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
			continue
		}

		response := orchestrator.Message{
			Header: orchestrator.Header{
				TransactionId:       request.Header.TransactionId,
				TransactionDateTime: request.Header.TransactionDateTime,
				OrderType:           request.Header.OrderType,
				OrderService:        "payment-deduct-balance",
				Retries:             request.Header.Retries,
				ResponseCode:        http.StatusOK,
				ResponseMessage:     http.StatusText(http.StatusOK),
			},
			Body: request.Body,
		}

		paymentRequest := payment.Request{}
		paymentQuantity := 0.0

		paymentRequest.UserId = request.Body.(map[string]interface{})["user_id"].(string)
		paymentQuantity = request.Body.(map[string]interface{})["quantity"].(float64)

		paymentPayload, err := json.Marshal(paymentRequest)
		if err != nil {
			log.Error().Err(err).Msg("Error marshalling payment request")
		}

		res, err := http.Post(paymentService, "application/json", bytes.NewBuffer(paymentPayload))
		if err != nil {
			log.Error().Err(err).Msg("Error sending payment request")
		}

		byteRes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Error().Err(err).Msg("Error reading payment response")
		}

		log.Info().Msg(string(byteRes))

		var paymentResponse payment.Response
		err = json.Unmarshal(byteRes, &paymentResponse)
		if err != nil {
			log.Error().Err(err).Msg("Error decoding payment response")
		}

		if paymentResponse.Message < paymentQuantity*2000 && err == nil {
			response.Header.ResponseCode = http.StatusBadRequest
			response.Header.ResponseMessage = http.StatusText(http.StatusBadRequest)
		}

		if err != nil {
			response.Header.ResponseCode = http.StatusInternalServerError
			response.Header.ResponseMessage = http.StatusText(http.StatusInternalServerError)
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
		}
	}
}
