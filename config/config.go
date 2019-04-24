package config

import (
	"github.com/studtool/auth-service/beans"
	"time"

	"github.com/studtool/common/config"
)

var (
	_ = func() *config.FlagVar {
		f := config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_LOG_ENV_VARS", false)
		if f.Value() {
			config.SetLogger(beans.Logger)
		}
		return f
	}()

	ServerPort = config.NewStringDefault("STUDTOOL_AUTH_SERVICE_PORT", "80")

	CorsAllowed         = config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_ALLOW_CORS", false)
	RequestsLogsEnabled = config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_LOG_REQUESTS", true)

	JwtKey     = config.NewStringDefault("STUDTOOL_JWT_KEY", "secret")
	JwtExpTime = config.NewTimeSecsDefault("STUDTOOL_JWT_EXPIRE_TIME", 5*time.Minute)

	RepositoriesEnabled = config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_ENABLE_REPOSITORIES", false)
	QueuesEnabled       = config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_ENABLE_QUEUES", false)

	StorageHost     = config.NewStringDefault("STUDTOOL_AUTH_STORAGE_HOST", "127.0.0.1")
	StoragePort     = config.NewStringDefault("STUDTOOL_AUTH_STORAGE_PORT", "5432")
	StorageDB       = config.NewStringDefault("STUDTOOL_AUTH_STORAGE_NAME", "auth")
	StorageUser     = config.NewStringDefault("STUDTOOL_AUTH_STORAGE_USER", "user")
	StoragePassword = config.NewStringDefault("STUDTOOL_AUTH_STORAGE_PASSWORD", "password")
	StorageSSL      = config.NewStringDefault("STUDTOOL_AUTH_STORAGE_SSL_MODE", "disable")

	UsersMqHost     = config.NewStringDefault("STUDTOOL_USERS_MQ_HOST", "127.0.0.1")
	UsersMqPort     = config.NewStringDefault("STUDTOOL_USERS_MQ_PORT", "5672")
	UsersMqUser     = config.NewStringDefault("STUDTOOL_USERS_MQ_USER", "user")
	UsersMqPassword = config.NewStringDefault("STUDTOOL_USERS_MQ_PASSWORD", "password")

	UsersMqConnNumRet = config.NewIntDefault("STUDTOOL_USERS_MQ_CONNECTION_NUM_RETRIES", 10)
	UsersMqConnRetItv = config.NewTimeSecsDefault("STUDTOOL_USERS_MQ_CONNECTION_RETRY_INTERVAL", 2*time.Second)

	CreatedUsersQueueName = config.NewStringDefault("STUDTOOL_CREATED_USERS_MQ_NAME", "created_users")
	DeletedUsersQueueName = config.NewStringDefault("STUDTOOL_DELETED_USERS_MQ_NAME", "deleted_users")
)
