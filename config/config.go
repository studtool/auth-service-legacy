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
	ServerPort = func() string {
		v := os.Getenv("STUDTOOL_AUTH_SERVICE_PORT")
		if v == "" {
			return "8080"
		}
		return v
	}()

	ConsulClientEnabled = func() bool {
		v := os.Getenv("STUDTOOL_AUTH_SERVICE_ENABLE_DISCOVERY_CLIENT")
		if v == "true" {
			return true
		}
		return false
	}()

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
