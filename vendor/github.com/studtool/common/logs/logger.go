package logs

import (
	"fmt"
	"runtime"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func NewLogger() *Logger {
	return &Logger{
		logger: func() *logrus.Logger {
			log := logrus.StandardLogger()
			log.SetFormatter(&logrus.JSONFormatter{})
			return log
		}(),
	}
}

func (log *Logger) Debug(args ...interface{}) {
	log.logger.WithFields(log.callerInfo()).Debug(args...)
}

func (log *Logger) Info(args ...interface{}) {
	log.logger.WithFields(log.callerInfo()).Info(args...)
}

func (log *Logger) Warning(args ...interface{}) {
	log.logger.WithFields(log.callerInfo()).Warn(args...)
}

func (log *Logger) Error(args ...interface{}) {
	log.logger.WithFields(log.callerInfo()).Error(args...)
}

func (log *Logger) Fatal(args ...interface{}) {
	log.logger.WithFields(log.callerInfo()).Fatal(args...)
}

const (
	callerStackDepth = 3
)

func (log *Logger) callerInfo() logrus.Fields {
	fpcs := make([]uintptr, 1)

	n := runtime.Callers(callerStackDepth, fpcs)
	if n == 0 {
		return nil
	}

	caller := runtime.FuncForPC(fpcs[0] - 1)
	if caller == nil {
		return nil
	}

	name := caller.Name()
	file, line := caller.FileLine(fpcs[0] - 1)

	return logrus.Fields{
		"func": name,
		"file": fmt.Sprintf("%s:%d", file, line),
	}
}
