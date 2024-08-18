package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"payment/internal/domain/orchestrator"
	"payment/internal/domain/payment"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

const paymentService = "http://127.0.0.1:8090/balance"

func (uc *PaymentUsecaseImpl) Start() {
	log.Info().Msg("Starting payment event listener...")

	producer := uc.kafka.NewProducer("topic-orchestrator")
	reader := uc.kafka.NewConsumer("topic-payment-deduct-balance")

	defer func() {
		producer.Close()
		reader.Close()
	}()

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

		paymentRequest := payment.Request{}
		paymentQuantity := 0.0
		if err == nil {
			paymentRequest.UserId = request.Body.(map[string]interface{})["user_id"].(string)
			if request.Body.(map[string]interface{})["quantity"] != nil {
				paymentQuantity = request.Body.(map[string]interface{})["quantity"].(float64)
			} else {
				err = fmt.Errorf("quantity is required")
			}
		}

		paymentResponse := payment.Response{
			Message: -1,
		}
		if err == nil {
			paymentPayload, err := json.Marshal(paymentRequest)
			if err != nil {
				log.Error().Err(err).Msg("Error marshalling payment request")
			}

			res, err := http.Post(paymentService, "application/json", bytes.NewBuffer(paymentPayload))
			if err != nil {
				log.Error().Err(err).Msg("Error sending payment request")
			}

			if err == nil {
				err = json.NewDecoder(res.Body).Decode(&paymentResponse)
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
				OrderService:        "payment-deduct-balance",
				Retries:             request.Header.Retries,
				ResponseCode:        http.StatusOK,
				ResponseMessage:     http.StatusText(http.StatusOK),
			},
			Body: request.Body,
		}

		if err != nil && err != fmt.Errorf("quantity is required") {
			response.Header.ResponseCode = http.StatusInternalServerError
			response.Header.ResponseMessage = http.StatusText(http.StatusInternalServerError)
		} else if paymentResponse.Message < paymentQuantity*2000 || err == fmt.Errorf("quantity is required") {
			response.Header.ResponseCode = http.StatusBadRequest
			response.Header.ResponseMessage = http.StatusText(http.StatusBadRequest)
		}

		payload, err := json.Marshal(response)
		if err != nil {
			log.Error().Err(err).Msg("Error marshalling response")
		}

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
