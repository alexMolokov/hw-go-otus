package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/app"
	configApp "github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/config"
	"github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/server/http"
	storagefactory "github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/storage/factory"
)

var configFile string

func init() {
	flag.StringVar(
		&configFile,
		"config",
		"C:\\Users\\Alexey\\go\\src\\hw-go-otus\\hw12_13_14_15_calendar\\configs\\config.json",
		"Path to configuration file",
	)
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cfg, err := configApp.NewConfig(configFile)
	if err != nil {
		fmt.Printf("Can't load config: %v", err)
		os.Exit(1)
	}

	logg, err := logger.New(&cfg.Logger)
	if err != nil {
		fmt.Printf("Can't create logger: %v", err)
		os.Exit(1)
	}

	st, err := storagefactory.NewStorage(cfg)
	if err != nil {
		fmt.Printf("Can't create pool connect to storage: %v", err)
		os.Exit(1)
	}

	calendar := app.New(logg, st)
	defer calendar.Close()

	server := internalhttp.NewServer(logg, calendar, fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port))

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		logg.Info("calendar is running...")

		if err := server.Start(); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	if err := server.Stop(ctx); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	} else {
		logg.Info("calendar is stopped")
	}
}
