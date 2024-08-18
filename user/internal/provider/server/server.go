package server

import (
	"context"
	"os/signal"
	"user/internal/provider/messaging/kafka"

	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// Start the Kafka consumer in the background
	go func() {
		// Add a short delay to allow the server to start
		time.Sleep(2 * time.Second)
		kafka.Start()
	}()

	<-ctx.Done()
	stop()
	log.Info().Msg("Shutting down server...")
	log.Info().Msg("Server stopped")
}
