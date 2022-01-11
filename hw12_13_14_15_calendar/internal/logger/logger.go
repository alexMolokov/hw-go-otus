package logger

import (
	configApp "github.com/alexMolokov/hw-go-otus/hw12_13_14_15_calendar/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zp *zap.SugaredLogger
}

func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func New(c *configApp.LoggerConf) (l *Logger) {
	defer func() {
		l.zp.Sync()
	}()

	cfg := zap.NewProductionConfig()

	cfg.Level = zap.NewAtomicLevelAt(getLevel(c.Level))
	cfg.Encoding = c.Encoding
	cfg.OutputPaths = []string{c.Output}
	cfg.ErrorOutputPaths = []string{c.Output}
	cfg.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	logger, _ := cfg.Build()

	return &Logger{
		zp: logger.Sugar(),
	}
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.zp.Debugf(msg, args...)
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.zp.Infof(msg, args...)
}

func (l *Logger) Warning(msg string, args ...interface{}) {
	l.zp.Warnf(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.zp.Errorf(msg, args...)
}
