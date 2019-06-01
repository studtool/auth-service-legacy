package logs

import (
	"github.com/sirupsen/logrus"
)

type RawLogger struct {
	logger *logrus.Logger
}

func NewRawLogger() Logger {
	return &RawLogger{
		logger: logrus.StandardLogger(),
	}
}

func (log *RawLogger) Debug(args ...interface{}) {
	log.logger.Debug(args...)
}

func (log *RawLogger) Debugf(format string, args ...interface{}) {
	log.logger.Debugf(format, args...)
}

func (log *RawLogger) Info(args ...interface{}) {
	log.logger.Info(args...)
}

func (log *RawLogger) Infof(format string, args ...interface{}) {
	log.logger.Infof(format, args...)
}

func (log *RawLogger) Warning(args ...interface{}) {
	log.logger.Warn(args...)
}

func (log *RawLogger) Warningf(format string, args ...interface{}) {
	log.logger.Warningf(format, args...)
}

func (log *RawLogger) Error(args ...interface{}) {
	log.logger.Error(args...)
}

func (log *RawLogger) Errorf(format string, args ...interface{}) {
	log.logger.Errorf(format, args...)
}

func (log *RawLogger) Fatal(args ...interface{}) {
	log.logger.Fatal(args...)
}

func (log *RawLogger) Fatalf(format string, args ...interface{}) {
	log.logger.Fatalf(format, args...)
}
