package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/khaldeezal/subscriptions-service/internal/config"
	"github.com/khaldeezal/subscriptions-service/internal/controller"
	"github.com/khaldeezal/subscriptions-service/internal/httpserver"
	"github.com/khaldeezal/subscriptions-service/internal/repository"
	"github.com/khaldeezal/subscriptions-service/internal/service"
)

func Run() error {
	_ = godotenv.Load()
	cfg := config.Load()

	lvl, err := zerolog.ParseLevel(cfg.LogLvl)
	if err == nil {
		zerolog.SetGlobalLevel(lvl)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Heartbeat("/healthz"))

	subscriptionRepository, err := repository.NewSubscriptionRepository(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}

	defer subscriptionRepository.ClosePool()

	subscriptionService := service.NewSubscriptionService(subscriptionRepository)
	subscriptionController := controller.NewController(subscriptionService)

	handler := subscriptionController.InitRoutes(r)

	srv := httpserver.New(cfg.Port, handler)

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	_ = srv.Shutdown(ctx)
	return nil
}
