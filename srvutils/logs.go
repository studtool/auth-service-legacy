package srvutils

import (
	"github.com/studtool/common/logs"
	"github.com/studtool/common/rft"

	"github.com/studtool/auth-service/config"
)

func MakeRawLogger(_ interface{}) logs.Logger {
	return logs.NewRawLogger()
}

func MakeStructLogger(v interface{}) logs.Logger {
	return logs.NewStructLogger(
		logs.StructLoggerParams{
			ComponentName:     config.ComponentName,
			ComponentVersion:  config.ComponentVersion,
			StructWithPkgName: rft.StructName(v),
		},
	)
}

func MakeReflectLogger(_ interface{}) logs.Logger {
	return logs.NewReflectLogger()
}

func MakeRequestLogger(_ interface{}) logs.Logger {
	return logs.NewRequestLogger(
		logs.RequestLoggerParams{
			ComponentName:    config.ComponentName,
			ComponentVersion: config.ComponentVersion,
		},
	)
}
