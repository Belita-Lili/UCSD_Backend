package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var globalLogger *zap.Logger

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) *zap.Logger
}

func Init(env string, logLevel string) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Set log level
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(logLevel)); err == nil {
		config.Level.SetLevel(level)
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	globalLogger = logger
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

func With(fields ...zap.Field) *zap.Logger {
	return globalLogger.With(fields...)
}
