package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/khaldeezal/subscriptions-service/internal/config"
	"github.com/khaldeezal/subscriptions-service/internal/httpserver"
	"github.com/khaldeezal/subscriptions-service/internal/storage"
)

func main() {
	_ = godotenv.Load()

	lvl, err := zerolog.ParseLevel(getenv("LOG_LEVEL", "info"))
	if err == nil {
		zerolog.SetGlobalLevel(lvl)
	}

	cfg := config.Load()

	pool, err := storage.NewPool(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("connect db")
	}
	defer pool.Close()

	srv := httpserver.New(cfg, pool)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
