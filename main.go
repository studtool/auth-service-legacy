package main

import (
	"auth-service/config"
	"auth-service/consul"
)

func main() {
	c := consul.NewClient()

	if config.ConsulClientEnabled {
		c.Register()
		defer c.Unregister()
	}

}
