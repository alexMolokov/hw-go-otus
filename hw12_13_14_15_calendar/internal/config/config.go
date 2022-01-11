package config

import (
	"context"
	"time"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

type Config struct {
	Logger   LoggerConf
	DriverDB string `config:"driverDb"`
	DB       DBConf
	HTTP     HTTPConf
}

type LoggerConf struct {
	Level    string `config:"logger-level"`
	Encoding string `config:"logger-encoding"`
	Output   string `config:"logger-output"`
}

type DBConf struct {
	Driver            string `config:"db-driver"`
	Host              string `config:"db-host"`
	Port              int    `config:"db-port"`
	Name              string `config:"db-name"`
	User              string `config:"db-user"`
	Password          string `config:"db-password"`
	MaxConnectionPool int    `config:"db-maxConnectionPool"`
	SslMode           string `config:"db-sslMode"`
}

type HTTPConf struct {
	Host string `config:"http-host"`
	Port int    `config:"http-port"`
}

func NewConfig(fileName string) (*Config, error) {
	loader := confita.NewLoader(
		file.NewBackend(fileName),
	)

	cfg := defaultConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := loader.Load(ctx, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		Logger: LoggerConf{
			Level: "DEBUG",
		},
		DB: DBConf{
			Driver:            "postgres",
			Host:              "localhost",
			Port:              5432,
			Name:              "calendar",
			User:              "user",
			Password:          "password",
			MaxConnectionPool: 2,
			SslMode:           "disable",
		},
		DriverDB: "memory",
		HTTP: HTTPConf{
			Host: "localhost",
			Port: 8083,
		},
	}
}
