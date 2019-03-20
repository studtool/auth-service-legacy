package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ServerPort = getEnvVar("STUDTOOL_AUTH_SERVICE_PORT", "80")

	StorageHost       = getEnvVar("STUDTOOL_AUTH_STORAGE_HOST", "127.0.0.1")
	StoragePort       = getEnvVar("STUDTOOL_AUTH_STORAGE_PORT", "5432")
	StorageDB         = getEnvVar("STUDTOOL_AUTH_STORAGE_NAME", "auth")
	StorageUser       = getEnvVar("STUDTOOL_AUTH_STORAGE_USER", "user")
	StoragePassword   = getEnvVar("STUDTOOL_AUTH_STORAGE_PASSWORD", "password")
	StorageSSL        = getEnvVar("STUDTOOL_AUTH_STORAGE_SSL_MODE", "disable")
	ShouldInitStorage = getEnvFlag("STUDTOOL_AUTH_STORAGE_SHOULD_INIT", false)

	ConsulClientEnabled = getEnvFlag("STUDTOOL_AUTH_SERVICE_DISCOVERY_CLIENT_ENABLED", false)

	ConsulAddress = func() string {
		v := os.Getenv("STUDTOOL_SERVICE_DISCOVERY_ADDRESS")
		if v == "" {
			return "127.0.0.1:8500"
		}
		return v
	}()

	HealthCheckTimeout = func() time.Duration {
		const ev = "STUDTOOL_AUTH_SERVICE_HEALTH_CHECK_TIMEOUT"

		v := os.Getenv(ev)
		if v == "" {
			return 10 * time.Second
		}

		if strings.Contains(v, "s") {
			v = v[:len(v)-1]

			t, err := strconv.Atoi(v)
			if err != nil {
				panic(err)
			}

			return time.Duration(t) * time.Second
		}

		panic(fmt.Sprintf("invalid %s", ev))
	}()

	Logger = logrus.StandardLogger()
)

func getEnvVar(name string, defaultValue string) string {
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
