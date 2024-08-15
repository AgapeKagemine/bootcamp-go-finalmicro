package kafka

import (
	"book/internal/domain/book"
	"book/internal/domain/orchestrator"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

const bookService = "http://127.0.0.1:8090/book"

func Start() {
	producer := NewProducer("topic-orchestrator")
	reader := NewConsumer("topic-book-validate")

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
				OrderService:        "book-validate",
				Retries:             request.Header.Retries,
				ResponseCode:        http.StatusOK,
				ResponseMessage:     http.StatusText(http.StatusOK),
			},
			Body: request.Body,
		}

		bookRequest := &book.Request{}

		bookRequest.BookId = request.Body.(map[string]interface{})["book_id"].(string)
		bookQuantity := request.Body.(map[string]interface{})["quantity"].(float64)

		bookPayload, err := json.Marshal(bookRequest)
		if err != nil {
			log.Error().Err(err).Msg("Error marshalling book request")
		}

		res, err := http.Post(bookService, "application/json", bytes.NewBuffer(bookPayload))
		if err != nil {
			log.Error().Err(err).Msg("Error sending book request")
		}

		byteRes, err := io.ReadAll(res.Body)
		if err != nil {
			log.Error().Err(err).Msg("Error reading book response")
		}

		log.Info().Msg(string(byteRes))

		var bookResponse book.Response
		err = json.Unmarshal(byteRes, &bookResponse)
		if err != nil {
			log.Error().Err(err).Msg("Error decoding book response")
		}

		log.Info().Msg(bookRequest.BookId)

		if err != nil {
			response.Header.ResponseCode = http.StatusInternalServerError
			response.Header.ResponseMessage = http.StatusText(http.StatusInternalServerError)
		}

		if (bookResponse.Message != "valid" || bookQuantity > 5) && err == nil {
			response.Header.ResponseCode = http.StatusBadRequest
			response.Header.ResponseMessage = http.StatusText(http.StatusBadRequest)
		}

		payload, err := json.Marshal(response)
		if err != nil {
			log.Error().Err(err).Msg("Error marshalling response")
			continue
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
