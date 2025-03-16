package log

import "log"

type Logger struct {
	Logger log.Logger
}

func (l *Logger) Info(str string, args ...interface{}) {
}
