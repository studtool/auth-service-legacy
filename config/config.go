package config

import (
	"auth-service/beans"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	shouldLogEnvVars = getEnvFlag("STUDTOOL_AUTH_SERVICE_SHOULD_LOG_ENV_VARS", false)

	ServerPort = getEnvStr("STUDTOOL_AUTH_SERVICE_PORT", "80")

	JwtKey = requireEnvStr("STUDTOOL_JWT_KEY")

	RepositoriesEnabled    = getEnvFlag("STUDTOOL_AUTH_SERVICE_REPOSITORIES_ENABLED", false)
	DiscoveryClientEnabled = getEnvFlag("STUDTOOL_AUTH_SERVICE_DISCOVERY_CLIENT_ENABLED", false)
	QueuesEnabled          = getEnvFlag("STUDTOOL_AUTH_SERVICE_QUEUES_ENABLED", false)

	StorageHost     = getEnvStr("STUDTOOL_AUTH_STORAGE_HOST", "127.0.0.1")
	StoragePort     = getEnvStr("STUDTOOL_AUTH_STORAGE_PORT", "5432")
	StorageDB       = getEnvStr("STUDTOOL_AUTH_STORAGE_NAME", "auth")
	StorageUser     = getEnvStr("STUDTOOL_AUTH_STORAGE_USER", "user")
	StoragePassword = getEnvStr("STUDTOOL_AUTH_STORAGE_PASSWORD", "password")
	StorageSSL      = getEnvStr("STUDTOOL_AUTH_STORAGE_SSL_MODE", "disable")

	UsersMqHost     = getEnvStr("STUDTOOL_USERS_MQ_HOST", "127.0.0.1")
	UsersMqPort     = getEnvStr("STUDTOOL_USERS_MQ_PORT", "5672")
	UsersMqUser     = getEnvStr("STUDTOOL_USERS_MQ_USER", "user")
	UsersMqPassword = getEnvStr("STUDTOOL_USERS_MQ_PASSWORD", "password")

	UsersMqConnNumRet = getEnvInt("STUDTOOL_USERS_MQ_CONNECTION_NUM_RETRIES", 10)
	UsersMqConnRetItv = getEnvTimeSec("STUDTOOL_USERS_MQ_CONNECTION_RETRY_INTERVAL", 2*time.Second)

	CreatedUsersQueueName = getEnvStr("STUDTOOL_CREATED_USERS_QUEUE_NAME", "created_users")
	DeletedUsersQueueName = getEnvStr("STUDTOOL_DELETED_USERS_QUEUE_NAME", "deleted_users")

	DiscoveryServiceAddress = getEnvStr("STUDTOOL_DISCOVERY_SERVICE_ADDRESS", "127.0.0.1:8500")
	HealthCheckTimeout      = getEnvTimeSec("STUDTOOL_AUTH_SERVICE_HEALTH_CHECK_TIMEOUT", 10*time.Second)
)

func getEnvStr(name string, defVal string) string {
	v := os.Getenv(name)
	if v == "" {
		v = defVal
	}

	if shouldLogEnvVars {
		beans.Logger.Infof("%s=%s", name, v)
	}

	return v
}

func requireEnvStr(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic(fmt.Sprintf("config error: %s required", name))
	}
	if shouldLogEnvVars {
		beans.Logger.Infof("%s=%s", name, v)
	}
	return v
}

func getEnvFlag(name string, defVal bool) (val bool) {
	v := os.Getenv(name)
	if v == "" {
		val = defVal
	} else if v == "true" {
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
	} else {
		var err error
		val, err = strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("%s: integer expected", name))
		}
	}

	if shouldLogEnvVars {
		beans.Logger.Infof("%s=%v", name, val)
	}

	return val
}

func getEnvTimeSec(name string, defVal time.Duration) (val time.Duration) {
	v := os.Getenv(name)
	if v == "" {
		val = defVal
	} else {
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
	}

	if shouldLogEnvVars {
		beans.Logger.Infof("%s=%v", name, val)
	}

	return val
}
