package config

import (
	"time"

	"github.com/studtool/common/config"

	"github.com/studtool/auth-service/beans"
)

var (
	_ = func() *cconfig.FlagVar {
		f := cconfig.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_LOG_ENV_VARS", false)
		if f.Value() {
			cconfig.SetLogger(beans.Logger)
		}
		return f
	}()

	ServerPort = cconfig.NewStringDefault("STUDTOOL_AUTH_SERVICE_PORT", "80")

	CorsAllowed         = cconfig.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_ALLOW_CORS", false)
	RequestsLogsEnabled = cconfig.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_LOG_REQUESTS", true)

	JwtKey     = cconfig.NewStringDefault("STUDTOOL_JWT_KEY", "secret")
	JwtExpTime = cconfig.NewTimeSecsDefault("STUDTOOL_JWT_EXPIRE_TIME", 5*time.Minute)

	RepositoriesEnabled = cconfig.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_ENABLE_REPOSITORIES", false)
	QueuesEnabled       = cconfig.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_ENABLE_QUEUES", false)

	StorageHost     = cconfig.NewStringDefault("STUDTOOL_AUTH_STORAGE_HOST", "127.0.0.1")
	StoragePort     = cconfig.NewStringDefault("STUDTOOL_AUTH_STORAGE_PORT", "5432")
	StorageDB       = cconfig.NewStringDefault("STUDTOOL_AUTH_STORAGE_NAME", "auth")
	StorageUser     = cconfig.NewStringDefault("STUDTOOL_AUTH_STORAGE_USER", "user")
	StoragePassword = cconfig.NewStringDefault("STUDTOOL_AUTH_STORAGE_PASSWORD", "password")
	StorageSSL      = cconfig.NewStringDefault("STUDTOOL_AUTH_STORAGE_SSL_MODE", "disable")

	MqHost     = cconfig.NewStringDefault("STUDTOOL_MQ_HOST", "127.0.0.1")
	MqPort     = cconfig.NewStringDefault("STUDTOOL_MQ_PORT", "5672")
	MqUser     = cconfig.NewStringDefault("STUDTOOL_MQ_USER", "user")
	MqPassword = cconfig.NewStringDefault("STUDTOOL_MQ_PASSWORD", "password")

	UsersMqConnNumRet = cconfig.NewIntDefault("STUDTOOL_USERS_MQ_CONNECTION_NUM_RETRIES", 10)
	UsersMqConnRetItv = cconfig.NewTimeSecsDefault("STUDTOOL_USERS_MQ_CONNECTION_RETRY_INTERVAL", 2*time.Second)
)
