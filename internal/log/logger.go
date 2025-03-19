package log

import (
	"go.uber.org/zap"
)

type Logger struct {
	zap.Logger
}

var logger *zap.Logger

func NewLogger() *Logger {
	if logger == nil {
		logger, _ = zap.NewDevelopment()
	}
	return &Logger{*logger}
}
