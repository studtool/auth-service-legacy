package main

import (
	"auth-service/api"
	"auth-service/config"
	"auth-service/consul"
	"os"
	"os/signal"
)

func main() {
	if config.ConsulClientEnabled {
		c := consul.NewClient()

		c.Register()
		defer c.Unregister()
	}

	srv := api.NewServer()

	go srv.Run()
	defer srv.Shutdown()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	<-ch
}
