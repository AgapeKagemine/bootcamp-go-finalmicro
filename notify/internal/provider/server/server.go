package server

import (
	"context"
	"notify/internal/provider/messaging/kafka"
	"notify/internal/usecase/notify"
	"os/signal"

	"syscall"

	"github.com/rs/zerolog/log"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	k := kafka.NewOrchestratorKafka()
	uc := notify.NewNotifyUsecase(k)
	go uc.Start()
	<-ctx.Done()
	stop()
	log.Info().Msg("Shutting down server...")
	log.Info().Msg("Server stopped")
}
