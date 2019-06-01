package beans

import (
	"github.com/studtool/common/logs"
)

var (
	logger = logs.NewRawLogger()
)

func Logger() logs.Logger {
	return logger
}
