package logs

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/studtool/common/types"
	"github.com/studtool/common/utils/process"
)

type RequestLogger struct {
	host string
	pid  int64

	componentName    string
	componentVersion string

	logger *logrus.Logger
}

type RequestLoggerParams struct {
	ComponentName    string
	ComponentVersion string
}

func NewRequestLogger(params RequestLoggerParams) Logger {
	return &RequestLogger{
		host: process.GetHost(),
		pid:  process.GetPid(),

		componentName:    params.ComponentName,
		componentVersion: params.ComponentVersion,

		logger: func() *logrus.Logger {
			log := logrus.StandardLogger()
			log.SetFormatter(&logrus.JSONFormatter{})
			return log
		}(),
	}
}

type RequestParams struct {
	Method      string
	Path        string
	Status      int
	Type        string
	UserID      types.ID
	IP          string
	UserAgent   string
	RequestTime time.Duration
}

const (
	RequestMessage     = "request handled"
	WrongMethodMessage = "method with format should not be called"
)

func (log *RequestLogger) Debug(args ...interface{}) {
	log.logger.WithFields(log.makeLogFields(args...)).Debug(RequestMessage)
}

func (log *RequestLogger) Debugf(format string, args ...interface{}) {
	panic(WrongMethodMessage)
}

func (log *RequestLogger) Info(args ...interface{}) {
	log.logger.WithFields(log.makeLogFields(args...)).Info(RequestMessage)
}

func (log *RequestLogger) Infof(format string, args ...interface{}) {
	panic(WrongMethodMessage)
}

func (log *RequestLogger) Warning(args ...interface{}) {
	log.logger.WithFields(log.makeLogFields(args...)).Warning(RequestMessage)
}

func (log *RequestLogger) Warningf(format string, args ...interface{}) {
	panic(WrongMethodMessage)
}

func (log *RequestLogger) Error(args ...interface{}) {
	log.logger.WithFields(log.makeLogFields(args...)).Error(RequestMessage)
}

func (log *RequestLogger) Errorf(format string, args ...interface{}) {
	panic(WrongMethodMessage)
}

func (log *RequestLogger) Fatal(args ...interface{}) {
	log.logger.WithFields(log.makeLogFields(args...)).Fatal(RequestMessage)
}

func (log *RequestLogger) Fatalf(format string, args ...interface{}) {
	panic(WrongMethodMessage)
}

func (log *RequestLogger) makeLogFields(args ...interface{}) logrus.Fields {
	if len(args) != 1 {
		panic(args)
	}

	p, ok := args[0].(*RequestParams)
	if !ok {
		panic(args)
	}

	return logrus.Fields{
		"host":        log.host,
		"pid":         log.pid,
		"component":   log.componentName,
		"version":     log.componentVersion,
		"method":      p.Method,
		"path":        p.Path,
		"status":      p.Status,
		"type":        p.Type,
		"userID":      p.UserID,
		"IP":          p.IP,
		"User-Agent":  p.UserAgent,
		"requestTime": p.RequestTime,
	}
}
