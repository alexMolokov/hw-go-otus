package app

import (
	st "github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warning(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type Storage interface {
	CreateEvent(event st.Event) (st.Event, error)
	Close() error
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) CreateEvent(event st.Event) (st.Event, error) {
	return a.Storage.CreateEvent(event)
}

func (a *App) Close() error {
	return a.Storage.Close()
}
