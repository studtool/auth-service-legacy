package config

import (
	"auth-service/beans"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	ServerPort = getEnvStr("STUDTOOL_AUTH_SERVICE_PORT", "80", true)

	StorageHost       = getEnvStr("STUDTOOL_AUTH_STORAGE_HOST", "127.0.0.1", true)
	StoragePort       = getEnvStr("STUDTOOL_AUTH_STORAGE_PORT", "5432", true)
	StorageDB         = getEnvStr("STUDTOOL_AUTH_STORAGE_NAME", "auth", true)
	StorageUser       = getEnvStr("STUDTOOL_AUTH_STORAGE_USER", "user", false)
	StoragePassword   = getEnvStr("STUDTOOL_AUTH_STORAGE_PASSWORD", "password", false)
	StorageSSL        = getEnvStr("STUDTOOL_AUTH_STORAGE_SSL_MODE", "disable", true)
	ShouldInitStorage = getEnvFlag("STUDTOOL_AUTH_STORAGE_SHOULD_INIT", false)

	StorageConnNumRet = getEnvInt("STUDTOOL_AUTH_STORAGE_CONNECTION_NUM_RETRIES", 10)
	StorageConnRetItv = getEnvTimeSec("STUDTOOL_AUTH_STORAGE_CONNECTION_RETRY_INTERVAL", 2*time.Second)

	UsersMqHost     = getEnvStr("STUDTOOL_USERS_MQ_HOST", "127.0.0.1", true)
	UsersMqPort     = getEnvStr("STUDTOOL_USERS_MQ_PORT", "5672", true)
	UsersMqUser     = getEnvStr("STUDTOOL_USERS_MQ_USER", "user", false)
	UsersMqPassword = getEnvStr("STUDTOOL_USERS_MQ_PASSWORD", "password", false)

	UsersMqConnNumRet = getEnvInt("STUDTOOL_USERS_MQ_CONNECTION_NUM_RETRIES", 10)
	UsersMqConnRetItv = getEnvTimeSec("STUDTOOL_USERS_MQ_CONNECTION_RETRY_INTERVAL", 2*time.Second)

	CreatedUsersQueueName = getEnvStr("STUDTOOL_CREATED_USERS_QUEUE_NAME", "created_users", true)
	DeletedUsersQueueName = getEnvStr("STUDTOOL_DELETED_USERS_QUEUE_NAME", "deleted_users", true)

	DiscoveryServiceAddress = getEnvStr("STUDTOOL_DISCOVERY_SERVICE_ADDRESS", "127.0.0.1:8500", true)
	DiscoveryClientEnabled  = getEnvFlag("STUDTOOL_AUTH_SERVICE_DISCOVERY_CLIENT_ENABLED", false)
	HealthCheckTimeout      = getEnvTimeSec("STUDTOOL_AUTH_SERVICE_HEALTH_CHECK_TIMEOUT", 10*time.Second)
)

func getEnvStr(name string, defVal string, withLog bool) string {
	v := os.Getenv(name)
	if v == "" {
		v = defVal
	}
	if withLog {
		beans.Logger.Infof("%s=%s", name, v)
	}
	return v
}

func getEnvFlag(name string, defVal bool) (val bool) {
	v := os.Getenv(name)
	if v == "" {
		val = defVal
	}
	if v == "true" {
		val = true
	} else if v == "false" {
		val = false
	}
	beans.Logger.Infof("%s=%v", name, val)
	return val
}

func getEnvInt(name string, defVal int) (val int) {
	v := os.Getenv(name)
	if v == "" {
		val = defVal
	}

	var err error
	val, err = strconv.Atoi(v)
	if err != nil {
		panic(fmt.Sprintf("%s: integer expected", name))
	}

	beans.Logger.Infof("%s=%v", name, val)
	return val
}

func getEnvTimeSec(name string, defVal time.Duration) (val time.Duration) {
	v := os.Getenv(name)
	if v == "" {
		val = defVal
	}

	if v[len(v)-1] == 's' {
		v = v[:len(v)-1]

		t, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}

		val = time.Duration(t) * time.Second
	} else {
		panic(fmt.Sprintf("%s: time expected. default: %fs", name, defVal.Seconds()))
	}

	beans.Logger.Infof("%s=%v", name, val)
	return val
}
