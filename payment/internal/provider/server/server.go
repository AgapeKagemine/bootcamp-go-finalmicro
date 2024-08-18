package server

import (
	"context"
	"os/signal"
	"payment/internal/provider/messaging/kafka"
	"payment/internal/usecase/payment"

	"syscall"

	"github.com/rs/zerolog/log"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	k := kafka.NewOrchestratorKafka()
	uc := payment.NewPaymentUsecase(k)
	go uc.Start()
	<-ctx.Done()
	stop()
	log.Info().Msg("Shutting down server...")
	log.Info().Msg("Server stopped")
}
