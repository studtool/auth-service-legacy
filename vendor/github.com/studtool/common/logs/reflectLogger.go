package logs

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"

	"github.com/studtool/common/utils/process"
)

type ReflectLogger struct {
	pid  int64
	host string

	logger *logrus.Logger
}

func NewReflectLogger() Logger {
	return &ReflectLogger{
		pid:  process.GetPid(),
		host: process.GetHost(),

		logger: func() *logrus.Logger {
			log := logrus.StandardLogger()
			log.SetFormatter(&logrus.JSONFormatter{})
			return log
		}(),
	}
}

func (log *ReflectLogger) Debug(args ...interface{}) {
	log.logger.WithFields(log.getCallerInfo()).Debug(args...)
}

func (log *ReflectLogger) Debugf(format string, args ...interface{}) {
	log.logger.WithFields(log.getCallerInfo()).Debugf(format, args...)
}

func (log *ReflectLogger) Info(args ...interface{}) {
	log.logger.WithFields(log.getCallerInfo()).Info(args...)
}

func (log *ReflectLogger) Infof(format string, args ...interface{}) {
	log.logger.WithFields(log.getCallerInfo()).Infof(format, args...)
}

func (log *ReflectLogger) Warning(args ...interface{}) {
	log.logger.WithFields(log.getCallerInfo()).Warn(args...)
}

func (log *ReflectLogger) Warningf(format string, args ...interface{}) {
	log.logger.WithFields(log.getCallerInfo()).Warningf(format, args...)
}

func (log *ReflectLogger) Error(args ...interface{}) {
	log.logger.WithFields(log.getCallerInfo()).Error(args...)
}

func (log *ReflectLogger) Errorf(format string, args ...interface{}) {
	log.logger.WithFields(log.getCallerInfo()).Errorf(format, args...)
}

func (log *ReflectLogger) Fatal(args ...interface{}) {
	log.logger.WithFields(log.getCallerInfo()).Fatal(args...)
}

func (log *ReflectLogger) Fatalf(format string, args ...interface{}) {
	log.logger.WithFields(log.getCallerInfo()).Fatalf(format, args...)
}

func (log *ReflectLogger) getCallerInfo() logrus.Fields {
	return logrus.Fields{
		"host":  log.host,
		"pid":   log.pid,
		"stack": debug.Stack(),
	}
}
