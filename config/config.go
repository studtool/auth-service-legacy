package config

import "os"

var (
	ServerPort = func() string {
		v := os.Getenv("STUDTOOL_AUTH_SERVICE_PORT")
		if v == "" {
			return "8080"
		}
		return v
	}()
)
