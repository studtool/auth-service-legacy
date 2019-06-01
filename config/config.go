package config

import (
	"time"

	"github.com/studtool/common/config"

	"github.com/studtool/auth-service/beans"
)

var (
	ComponentName    = "auth-service"
	ComponentVersion = "v0.0.1"

	_ = func() *config.FlagVar {
		f := config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_LOG_ENV_VARS", false)
		if f.Value() {
			config.SetLogger(beans.Logger())
		}
		return f
	}()

	ServerPort = config.NewIntDefault("STUDTOOL_AUTH_SERVICE_PORT", 80)

	CorsAllowed         = config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_ALLOW_CORS", false)
	RequestsLogsEnabled = config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_LOG_REQUESTS", true)

	JwtKey         = config.NewStringDefault("STUDTOOL_JWT_KEY", "secret")
	JwtValidPeriod = config.NewTimeDefault("STUDTOOL_JWT_VALID_PERIOD", 5*time.Minute)

	VerificationRequired = config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_CHECK_ACCOUNT_VERIFIED_ON_SIGN_IN", true)

	RepositoriesEnabled = config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_ENABLE_REPOSITORIES", false)
	QueuesEnabled       = config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_ENABLE_QUEUES", false)

	AuthStorageHost     = config.NewStringDefault("STUDTOOL_AUTH_STORAGE_HOST", "127.0.0.1")
	AuthStoragePort     = config.NewIntDefault("STUDTOOL_AUTH_STORAGE_PORT", 5432)
	AuthStorageName     = config.NewStringDefault("STUDTOOL_AUTH_STORAGE_NAME", "auth")
	AuthStorageUser     = config.NewStringDefault("STUDTOOL_AUTH_STORAGE_USER", "user")
	AuthStoragePassword = config.NewStringDefault("STUDTOOL_AUTH_STORAGE_PASSWORD", "password")
	AuthStorageSSL      = config.NewStringDefault("STUDTOOL_AUTH_STORAGE_SSL_MODE", "disable")

	TokensStorageHost = config.NewStringDefault("STUDTOOL_TOKENS_STORAGE_HOST", "127.0.0.1")
	TokensStoragePort = config.NewIntDefault("STUDTOOL_TOKENS_STORAGE_PORT", 6379)

	MqHost     = config.NewStringDefault("STUDTOOL_MQ_HOST", "127.0.0.1")
	MqPort     = config.NewIntDefault("STUDTOOL_MQ_PORT", 5672)
	MqUser     = config.NewStringDefault("STUDTOOL_MQ_USER", "user")
	MqPassword = config.NewStringDefault("STUDTOOL_MQ_PASSWORD", "password")

	MqConnNumRet = config.NewIntDefault("STUDTOOL_MQ_CONNECTION_NUM_RETRIES", 10)
	MqConnRetItv = config.NewTimeDefault("STUDTOOL_MQ_CONNECTION_RETRY_INTERVAL", 2*time.Second)
)
