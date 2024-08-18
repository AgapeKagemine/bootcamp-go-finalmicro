package server

import (
	"book/internal/provider/messaging/kafka"
	"book/internal/usecase/book"
	"context"
	"os/signal"

	"syscall"

	"github.com/rs/zerolog/log"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	k := kafka.NewOrchestratorKafka()
	uc := book.NewBookUsecase(k)
	go uc.Start()
	go uc.Rollback()
	<-ctx.Done()
	stop()
	log.Info().Msg("Shutting down server...")
	log.Info().Msg("Server stopped")
}
