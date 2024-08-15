package server

import (
	"context"
	"fmt"
	"net/http"
	"order/internal/config"
	"order/internal/provider/routes"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	address = "127.0.0.1"
	port    = 8060
)

func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	serverConfig := config.ServerConfig{
		Address: address,
		Port:    port,
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", serverConfig.Address, serverConfig.Port),
		Handler: routes.NewRoutes(),
	}

	go func() {
		log.Info().Msg(fmt.Sprintf("Starting server on port %d...", serverConfig.Port))

		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error().Err(err).Msg("Error starting server")
		}
	}()

	<-ctx.Done()
	stop()
	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Error shutting down server")
	}

	log.Info().Msg("HTTP server stopped")
}
