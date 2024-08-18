package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"orchestrator/internal/domain"
	orchestratorDomain "orchestrator/internal/domain/orchestrator"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

func (uc *OrchestratorUsecaseImpl) Start() {
	log.Info().Msg("Starting orchestrator event listener...")

	reader := uc.orchestratorKafka.NewConsumer("topic-orchestrator")
	defer reader.Close()

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Error while reading message: %v\n", err))
			continue
		}

		var request orchestratorDomain.Message
		err = json.Unmarshal(message.Value, &request)
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Error parsing message: %v\n", err))
			request.Header.ResponseMessage = err.Error()
			request.Header.ResponseCode = http.StatusInternalServerError
		}

		log.Info().Msg(fmt.Sprintf("Message received: %s\n", string(message.Value)))

		ct := context.WithValue(context.Background(), domain.Key("type"), request.Header.OrderType)
		ctx := context.WithValue(ct, domain.Key("service"), request.Header.OrderService)

		orchestrate, err := uc.orchestratorRepository.GetConfig(ctx)
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Error getting orchestration config: %v\n", err))
			request.Header.ResponseMessage = err.Error()
			request.Header.ResponseCode = http.StatusInternalServerError
		}

		writer := &kafka.Writer{}
		var responseBytes []byte

		ctxs := context.WithValue(ctx, domain.Key("request"), request)
		uc.orchestratorRepository.Save(ctxs)

		if request.Header.ResponseCode > http.StatusCreated {
			if (request.Header.ResponseCode == http.StatusBadRequest || request.Header.Retries >= uc.orchestratorKafka.GetMaxRetries()) && request.Header.ResponseMessage != "Manual Retry" {
				if orchestrate.RollbackTopic != "FINISH" {
					writer = uc.orchestratorKafka.NewProducer("topic-orchestrator")
					orchestrate.Topic = "orchestrator"
					request.Header.OrderService = orchestrate.RollbackTopic
					request.Header.ResponseCode = http.StatusOK
					request.Header.ResponseMessage = "ROLLBACK"
				} else {
					writer = uc.orchestratorKafka.NewProducer("topic-retry-step")
					orchestrate.Topic = "retry-step"
				}
			} else {
				request.Header.Retries = request.Header.Retries + 1
				request.Header.ResponseMessage = "AUTO RETRYING"
				request.Header.ResponseCode = http.StatusOK

				// Simple backoff
				time.Sleep(time.Duration(request.Header.Retries) * time.Second)

				writer = uc.orchestratorKafka.NewProducer("topic-" + request.Header.OrderService)
				orchestrate.Topic = request.Header.OrderService
			}
		} else if request.Header.OrderService == "FINISH" {
			request.Header.ResponseCode = http.StatusOK
			request.Header.ResponseMessage = http.StatusText(http.StatusOK)
			log.Info().Msg(fmt.Sprintf("Order completed: %s\n", string(message.Value)))
		} else {
			writer = uc.orchestratorKafka.NewProducer("topic-" + orchestrate.Topic)
		}

		if orchestrate.Topic != "FINISH" || request.Header.ResponseCode > http.StatusOK {
			responseBytes, err = json.Marshal(request)
			if err != nil {
				log.Error().Err(err).Msg(fmt.Sprintf("Error marshalling message: %v\n", err))
				request.Header.ResponseMessage = err.Error()
				request.Header.ResponseCode = http.StatusInternalServerError
			}

			err = writer.WriteMessages(context.Background(),
				kafka.Message{
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
		writer.Close()
	}
}
