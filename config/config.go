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
			cconfig.SetLogger(beans.Logger())
		}
		return f
	}()

	ServerPort = cconfig.NewStringDefault("STUDTOOL_AUTH_SERVICE_PORT", "80")

	CorsAllowed         = cconfig.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_ALLOW_CORS", false)
	RequestsLogsEnabled = cconfig.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_LOG_REQUESTS", true)

	JwtKey         = cconfig.NewStringDefault("STUDTOOL_JWT_KEY", "secret")
	JwtValidPeriod = cconfig.NewTimeDefault("STUDTOOL_JWT_VALID_PERIOD", 5*time.Minute)

	VerificationRequired = cconfig.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_CHECK_ACCOUNT_VERIFIED_ON_SIGN_IN", true)

	RepositoriesEnabled = cconfig.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_ENABLE_REPOSITORIES", false)
	QueuesEnabled       = cconfig.NewFlagDefault("STUDTOOL_AUTH_SERVICE_SHOULD_ENABLE_QUEUES", false)

	AuthStorageHost     = cconfig.NewStringDefault("STUDTOOL_AUTH_STORAGE_HOST", "127.0.0.1")
	AuthStoragePort     = cconfig.NewIntDefault("STUDTOOL_AUTH_STORAGE_PORT", 5432)
	AuthStorageName     = cconfig.NewStringDefault("STUDTOOL_AUTH_STORAGE_NAME", "auth")
	AuthStorageUser     = cconfig.NewStringDefault("STUDTOOL_AUTH_STORAGE_USER", "user")
	AuthStoragePassword = cconfig.NewStringDefault("STUDTOOL_AUTH_STORAGE_PASSWORD", "password")
	AuthStorageSSL      = cconfig.NewStringDefault("STUDTOOL_AUTH_STORAGE_SSL_MODE", "disable")

	TokensStorageHost = cconfig.NewStringDefault("STUDTOOL_TOKENS_STORAGE_HOST", "127.0.0.1")
	TokensStoragePort = cconfig.NewIntDefault("STUDTOOL_TOKENS_STORAGE_PORT", 6379)

	MqHost     = cconfig.NewStringDefault("STUDTOOL_MQ_HOST", "127.0.0.1")
	MqPort     = cconfig.NewIntDefault("STUDTOOL_MQ_PORT", 5672)
	MqUser     = cconfig.NewStringDefault("STUDTOOL_MQ_USER", "user")
	MqPassword = cconfig.NewStringDefault("STUDTOOL_MQ_PASSWORD", "password")

	MqConnNumRet = cconfig.NewIntDefault("STUDTOOL_MQ_CONNECTION_NUM_RETRIES", 10)
	MqConnRetItv = cconfig.NewTimeDefault("STUDTOOL_MQ_CONNECTION_RETRY_INTERVAL", 2*time.Second)
)
