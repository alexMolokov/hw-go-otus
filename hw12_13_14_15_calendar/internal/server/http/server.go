package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	App        Application
	httpServer *http.Server
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warning(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type Application interface{}

func NewServer(logger Logger, app Application, addr string) *Server {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "ok")
		},
	))

	return &Server{
		App: app,
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      loggingMiddleware(logger, mux),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	err := s.httpServer.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
