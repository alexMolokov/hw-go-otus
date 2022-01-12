package internalhttp

import (
	"context"
	"errors"
	"net/http"
	"time"
)

type Server struct {
	App        Application
	Logger     Logger
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
	s := &Server{
		App:    app,
		Logger: logger,
	}

	router := s.NewRouter()
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.loggingMiddleware(router),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s
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
