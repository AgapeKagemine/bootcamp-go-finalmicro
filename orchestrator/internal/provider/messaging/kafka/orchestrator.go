package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"orchestrator/internal/domain"
	"orchestrator/internal/domain/orchestrator"
	"orchestrator/internal/provider/database"
	orchestratorRepository "orchestrator/internal/repository/orchestrator"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

func Start() {
	reader := NewConsumer("topic-orchestrator")
	defer reader.Close()

	db, err := database.NewDB()
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Error connecting to the database: %v\n", err))
	}

	repo := orchestratorRepository.NewOrchestratorRepository(db)

	go ConsumeFailedTransaction()

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Error while reading message: %v\n", err))
			continue
		}

		var request orchestrator.Message
		err = json.Unmarshal(message.Value, &request)
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Error parsing message: %v\n", err))
			request.Header.ResponseMessage = err.Error()
			request.Header.ResponseCode = http.StatusInternalServerError
		}

		log.Info().Msg(fmt.Sprintf("Message received: %s\n", string(message.Value)))

		responseBytes, err := json.Marshal(request)
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Error marshalling message: %v\n", err))
			request.Header.ResponseMessage = err.Error()
			request.Header.ResponseCode = http.StatusInternalServerError
		}

		ct := context.WithValue(context.Background(), domain.Key("type"), request.Header.OrderType)
		ctx := context.WithValue(ct, domain.Key("service"), request.Header.OrderService)

		orchestrate, err := repo.GetConfig(ctx)
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Error getting orchestration config: %v\n", err))
			request.Header.ResponseMessage = err.Error()
			request.Header.ResponseCode = http.StatusInternalServerError
		}

		writer := NewProducer("topic-" + orchestrate.Topic)

		if orchestrate.OrderService == "START" {
			err = writer.WriteMessages(context.Background(),
				kafka.Message{
					Key:   []byte(request.Header.TransactionId),
					Value: responseBytes,
				},
			)

			if err != nil {
				log.Error().Err(err).Msg(fmt.Sprintf("Error writing message to topic-%s: %v\n", orchestrate.Topic, err))
				request.Header.ResponseMessage = err.Error()
				request.Header.ResponseCode = http.StatusInternalServerError
			} else {
				log.Info().Msg(fmt.Sprintf("Message sent to topic-%s: %s\n", orchestrate.Topic, string(responseBytes)))
			}
		} else if request.Header.ResponseCode != http.StatusOK {
			if request.Header.Retries < 3 {
				request.Header.Retries = request.Header.Retries + 1
				request.Header.ResponseCode = http.StatusOK
				request.Header.ResponseMessage = http.StatusText(http.StatusOK)

				responseBytes, err := json.Marshal(request)
				if err != nil {
					log.Error().Err(err).Msg(fmt.Sprintf("Error marshalling message: %v\n", err))
					request.Header.ResponseMessage = err.Error()
					request.Header.ResponseCode = http.StatusInternalServerError
				}

				writer2 := NewProducer("topic-" + request.Header.OrderService)
				err = writer2.WriteMessages(context.Background(),
					kafka.Message{
						Key:   []byte(request.Header.TransactionId),
						Value: responseBytes,
					},
				)

				if err != nil {
					log.Error().Err(err).Msg(fmt.Sprintf("Error writing message to topic-%s: %v\n", orchestrate.Topic, err))
					request.Header.ResponseMessage = err.Error()
					request.Header.ResponseCode = http.StatusInternalServerError
				} else {
					log.Info().Msg(fmt.Sprintf("Message sent to %s: %s\n", writer2.Topic, string(responseBytes)))
				}

				writer2.Close()
			} else {
				writer2 := NewProducer("topic-retry-step")

				err = writer2.WriteMessages(context.Background(),
					kafka.Message{
						Key:   []byte(request.Header.TransactionId),
						Value: responseBytes,
					},
				)

				responseBytes, err := json.Marshal(request)
				if err != nil {
					log.Error().Err(err).Msg(fmt.Sprintf("Error marshalling message: %v\n", err))
					request.Header.ResponseMessage = err.Error()
					request.Header.ResponseCode = http.StatusInternalServerError
				}

				if err != nil {
					log.Error().Err(err).Msg(fmt.Sprintf("Error writing message to topic-%s: %v\n", orchestrate.Topic, err))
					request.Header.ResponseMessage = err.Error()
					request.Header.ResponseCode = http.StatusInternalServerError
				} else {
					log.Info().Msg(fmt.Sprintf("Message sent to %s: %s\n", writer2.Topic, string(responseBytes)))
				}

				writer2.Close()
			}
		} else if orchestrate.Topic == "FINISH" {
			request.Header.ResponseCode = http.StatusOK
			request.Header.ResponseMessage = http.StatusText(http.StatusOK)
			log.Info().Msg(fmt.Sprintf("Order completed: %s\n", string(message.Value)))
		} else {
			if orchestrate.OrderService == request.Header.OrderService {

				err = writer.WriteMessages(context.Background(),
					kafka.Message{
						Key:   []byte(request.Header.TransactionId),
						Value: responseBytes,
					},
				)

				if err != nil {
					log.Error().Err(err).Msg(fmt.Sprintf("Error writing message to topic-%s: %v\n", orchestrate.Topic, err))
					request.Header.ResponseMessage = err.Error()
					request.Header.ResponseCode = http.StatusInternalServerError
				} else {
					log.Info().Msg(fmt.Sprintf("Message sent to topic-%s: %s\n", orchestrate.Topic, string(responseBytes)))
				}

			}
		}
		ctxs := context.WithValue(ctx, domain.Key("request"), request)
		repo.Save(ctxs)
		writer.Close()
	}
}
