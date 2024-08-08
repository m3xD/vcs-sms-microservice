package util

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() *zap.Logger {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{
		"stdout",
		"app.log",
	}

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logs, err := config.Build()

	if err != nil {
		panic(err)
	}
	return logs
}
