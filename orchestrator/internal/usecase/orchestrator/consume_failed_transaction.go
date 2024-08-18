package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	orchestratorDomain "orchestrator/internal/domain/orchestrator"

	"github.com/rs/zerolog/log"
)

// Should Dead Letter Message... but I don't know what to put
func (uc *OrchestratorUsecaseImpl) ConsumeFailedTransaction() {
	log.Info().Msg("Starting retry step consumer...")
	reader := uc.orchestratorKafka.NewConsumer("topic-retry-step")
	defer reader.Close()

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Error while reading message: %v\n", err))
		}

		var request orchestratorDomain.Message
		if err := json.Unmarshal(message.Value, &request); err != nil {
			log.Error().Err(err).Msg(fmt.Sprintf("Error parsing message: %v\n", err))
			continue
		}

		log.Info().Msg(fmt.Sprintf("Message received: %s\n", string(message.Value)))

		// ctx := context.WithValue(context.Background(), domain.Key("request"), request)

		// err = repo.SaveFailed(ctx)
		// if err != nil {
		// 	log.Error().Err(err).Msg(fmt.Sprintf("Error saving: %v\n", err))
		// 	continue
		// }
	}
}
