package logs

import (
	"github.com/sirupsen/logrus"
	"github.com/studtool/common/utils/process"
)

type StructLogger struct {
	logger *logrus.Logger
	fields logrus.Fields
}

type StructLoggerParams struct {
	ComponentName     string
	ComponentVersion  string
	StructWithPkgName string
}

func NewStructLogger(params StructLoggerParams) Logger {
	return &StructLogger{
		logger: func() *logrus.Logger {
			log := logrus.StandardLogger()
			log.SetFormatter(&logrus.JSONFormatter{})
			return log
		}(),
		fields: logrus.Fields{
			"pid":       process.GetPid(),
			"host":      process.GetHost(),
			"component": params.ComponentName,
			"version":   params.ComponentVersion,
			"struct":    params.StructWithPkgName,
		},
	}
}

func (log *StructLogger) Debug(args ...interface{}) {
	log.logger.WithFields(log.fields).Debug(args...)
}

func (log *StructLogger) Debugf(format string, args ...interface{}) {
	log.logger.WithFields(log.fields).Debugf(format, args...)
}

func (log *StructLogger) Info(args ...interface{}) {
	log.logger.WithFields(log.fields).Info(args...)
}

func (log *StructLogger) Infof(format string, args ...interface{}) {
	log.logger.WithFields(log.fields).Infof(format, args...)
}

func (log *StructLogger) Warning(args ...interface{}) {
	log.logger.WithFields(log.fields).Warn(args...)
}

func (log *StructLogger) Warningf(format string, args ...interface{}) {
	log.logger.WithFields(log.fields).Warningf(format, args...)
}

func (log *StructLogger) Error(args ...interface{}) {
	log.logger.WithFields(log.fields).Error(args...)
}

func (log *StructLogger) Errorf(format string, args ...interface{}) {
	log.logger.WithFields(log.fields).Errorf(format, args...)
}

func (log *StructLogger) Fatal(args ...interface{}) {
	log.logger.WithFields(log.fields).Fatal(args...)
}

func (log *StructLogger) Fatalf(format string, args ...interface{}) {
	log.logger.WithFields(log.fields).Fatalf(format, args...)
}
