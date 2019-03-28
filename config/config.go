package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	ServerPort = getEnvStr("STUDTOOL_AUTH_SERVICE_PORT", "80")

	StorageHost       = getEnvStr("STUDTOOL_AUTH_STORAGE_HOST", "127.0.0.1")
	StoragePort       = getEnvStr("STUDTOOL_AUTH_STORAGE_PORT", "5432")
	StorageDB         = getEnvStr("STUDTOOL_AUTH_STORAGE_NAME", "auth")
	StorageUser       = getEnvStr("STUDTOOL_AUTH_STORAGE_USER", "user")
	StoragePassword   = getEnvStr("STUDTOOL_AUTH_STORAGE_PASSWORD", "password")
	StorageSSL        = getEnvStr("STUDTOOL_AUTH_STORAGE_SSL_MODE", "disable")
	ShouldInitStorage = getEnvFlag("STUDTOOL_AUTH_STORAGE_SHOULD_INIT", false)

	MessageQueueHost     = getEnvStr("STUDTOOL_MQ_HOST", "127.0.0.1")
	MessageQueuePort     = getEnvStr("STUDTOOL_MQ_PORT", "5672")
	MessageQueueUser     = getEnvStr("STUDTOOL_MQ_USER", "user")
	MessageQueuePassword = getEnvStr("STUDTOOL_MQ_PASSWORD", "password")

	UsersMqConnNumRet = getEnvInt("STUDTOOL_USERS_MQ_CONNECTION_NUM_RETRIES", 10)
	UsersMqConnRetItv = getEnvTimeSec("STUDTOOL_USERS_MQ_CONNECTION_RETRIES_INTERVAL", 2*time.Second)

	CreatedUsersQueueName = getEnvStr("STUDTOOL_CREATED_USERS_QUEUE_NAME", "created_users")
	DeletedUsersQueueName = getEnvStr("STUDTOOL_DELETED_USERS_QUEUE_NAME", "deleted_users")

	ConsulClientEnabled = getEnvFlag("STUDTOOL_AUTH_SERVICE_DISCOVERY_CLIENT_ENABLED", false)

	ConsulAddress = getEnvStr("STUDTOOL_SERVICE_DISCOVERY_ADDRESS", "127.0.0.1:8500")

	HealthCheckTimeout = getEnvTimeSec("STUDTOOL_AUTH_SERVICE_HEALTH_CHECK_TIMEOUT", 10*time.Second)
)

func getEnvStr(name string, defaultValue string) string {
	v := os.Getenv(name)
	if v == "" {
		return defaultValue
	}
	return v
}

func getEnvFlag(name string, defaultValue bool) bool {
	v := os.Getenv(name)
	if v == "true" {
		return true
	} else if v == "false" {
		return false
	}
	return defaultValue
}

func getEnvInt(name string, defVal int) int {
	v := os.Getenv(name)
	if v == "" {
		return defVal
	}

	val, err := strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf("%s: integer expected", name))
	}

	return val
}

func getEnvTimeSec(name string, defVal time.Duration) time.Duration {
	v := os.Getenv(name)
	if v == "" {
		return defVal
	}

	if v[len(v)-1] == 's' {
		v = v[:len(v)-1]

		t, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}

		return time.Duration(t) * time.Second
	}

	panic(fmt.Sprintf("%s: time expected. default: %fs", name, defVal.Seconds()))
}
