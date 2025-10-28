package httpserver

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	port string
	http *http.Server
}

func New(port string, handler http.Handler) *Server {
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
	return &Server{
		port: port,
		http: srv,
	}
}

func (s *Server) Start() error {
	fmt.Printf("listening on :%s\n", s.port)
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
