package config

import (
	"time"

	"github.com/studtool/common/config"
)

var (
	ServerPort = config.NewStringDefault("STUDTOOL_AUTH_SERVICE_PORT", "80")

	JwtKey     = config.NewStringDefault("STUDTOOL_JWT_KEY", "secret")
	JwtExpTime = config.NewTimeSecsDefault("STUDTOOL_JWT_EXP_TIME", 5*time.Minute)

	RepositoriesEnabled    = config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_REPOSITORIES_ENABLED", false)
	DiscoveryClientEnabled = config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_DISCOVERY_CLIENT_ENABLED", false)
	QueuesEnabled          = config.NewFlagDefault("STUDTOOL_AUTH_SERVICE_QUEUES_ENABLED", false)

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

	CreatedUsersQueueName = config.NewStringDefault("STUDTOOL_CREATED_USERS_QUEUE_NAME", "created_users")
	DeletedUsersQueueName = config.NewStringDefault("STUDTOOL_DELETED_USERS_QUEUE_NAME", "deleted_users")

	DiscoveryServiceAddress = config.NewStringDefault("STUDTOOL_DISCOVERY_SERVICE_ADDRESS", "127.0.0.1:8500")
	HealthCheckTimeout      = config.NewTimeSecsDefault("STUDTOOL_AUTH_SERVICE_HEALTH_CHECK_TIMEOUT", 10*time.Second)
)
