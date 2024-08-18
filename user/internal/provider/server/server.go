package server

import (
	"context"
	"os/signal"
	"user/internal/provider/messaging/kafka"
	"user/internal/provider/usecase/user"

	"syscall"

	"github.com/rs/zerolog/log"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	k := kafka.NewOrchestratorKafka()
	uc := user.NewUserUsecase(k)

	go uc.Start()
	go uc.ValidatePassword()
	go uc.StatusCheck()
	go uc.AccountDecision()

	<-ctx.Done()
	stop()
	log.Info().Msg("Shutting down server...")
	log.Info().Msg("Server stopped")
}
