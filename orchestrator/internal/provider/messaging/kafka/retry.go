package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"orchestrator/internal/domain"
	"orchestrator/internal/domain/orchestrator"
	"orchestrator/internal/provider/database"
	orchestratorRepository "orchestrator/internal/repository/orchestrator"

	"github.com/rs/zerolog/log"
)

func ConsumeFailedTransaction() {
	reader := NewConsumer("topic-retry-step")
	defer reader.Close()

	db, err := database.NewDB()
	if err != nil {
		log.Error().Err(err).Msg(fmt.Sprintf("Error connecting to the database: %v\n", err))
	}

	repo := orchestratorRepository.NewOrchestratorRepository(db)

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Error while reading message: %v\n", err))
		}

		var request orchestrator.Message
		if err := json.Unmarshal(message.Value, &request); err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Error parsing message: %v\n", err))
			continue
		}

		log.Info().Msg(fmt.Sprintf("Message received: %s\n", string(message.Value)))

		ctx := context.WithValue(context.Background(), domain.Key("request"), request)

		err = repo.SaveFailed(ctx)
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Error saving: %v\n", err))
			continue
		}
	}
}
