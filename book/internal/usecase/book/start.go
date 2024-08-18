package book

import (
	"book/internal/domain/book"
	"book/internal/domain/orchestrator"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

const bookService = "http://127.0.0.1:8090/book"

func (uc *BookUsecase) Start() {
	log.Info().Msg("Starting book event listener...")
	reader := uc.kafka.NewConsumer("topic-book-validate")
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

		bookRequest := book.Request{}
		bookQuantity := 0.0
		if err == nil {
			if request.Body.(map[string]interface{})["book_id"] != nil {
				bookRequest.BookId = request.Body.(map[string]interface{})["book_id"].(string)
			} else {
				err = fmt.Errorf("book_id is required")
			}
			if request.Body.(map[string]interface{})["quantity"] != nil {
				bookQuantity = request.Body.(map[string]interface{})["quantity"].(float64)
			} else {
				err = fmt.Errorf("quantity is required")
			}
		}

		var bookResponse book.Response
		if err == nil {
			bookPayload, err := json.Marshal(bookRequest)
			if err != nil {
				log.Error().Err(err).Msg("Error marshalling book request")
			}

			res, err := http.Post(bookService, "application/json", bytes.NewBuffer(bookPayload))
			if err != nil {
				log.Error().Err(err).Msg("Error sending book request")
			}

			if err == nil {
				err = json.NewDecoder(res.Body).Decode(&bookResponse)
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
				OrderService:        "book-validate",
				Retries:             request.Header.Retries,
				ResponseCode:        http.StatusOK,
				ResponseMessage:     http.StatusText(http.StatusOK),
			},
			Body: request.Body,
		}

		log.Error().Err(err).Msg("Error validating book")

		if err != nil && err.Error() != "book_id is required" && err.Error() != "quantity is required" {
			response.Header.ResponseCode = http.StatusInternalServerError
			response.Header.ResponseMessage = http.StatusText(http.StatusInternalServerError)
		} else if (bookResponse.Message != "valid" || bookQuantity < 1 || bookQuantity > 5) || err == fmt.Errorf("book_id is required") || err == fmt.Errorf("quantity is required") {
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
