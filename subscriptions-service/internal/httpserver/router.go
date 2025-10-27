package httpserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"

	"github.com/khaldeezal/subscriptions-service/internal/config"
	"github.com/khaldeezal/subscriptions-service/internal/subscriptions"
)

type Server struct {
	cfg  config.Config
	pool *pgxpool.Pool
	http *http.Server
}

func New(cfg config.Config, pool *pgxpool.Pool) *Server {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(hlog.NewHandler(log.Logger))
	r.Use(hlog.RequestIDHandler("req_id", "Request-Id"))
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Heartbeat("/healthz"))

	r.Route("/api/v1", func(r chi.Router) {
		h := subscriptions.NewHandler(pool)

		r.Mount("/subscriptions", h.Routes()) // /subscriptions/ → попадёт в r.Post("/", ...)

		// Fallback — принимаем оба варианта коллекционных путей на верхнем уровне:
		r.Post("/subscriptions", h.Create)  // без слэша
		r.Post("/subscriptions/", h.Create) // со слэшем
		r.Get("/subscriptions", h.List)
		r.Get("/subscriptions/", h.List)

		r.Get("/cost", h.GetTotalCost)
	})
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}
	return &Server{cfg: cfg, pool: pool, http: srv}
}

func (s *Server) Start() error {
	fmt.Printf("listening on :%s\n", s.cfg.Port)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
