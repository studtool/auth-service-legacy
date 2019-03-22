package main

import (
	"auth-service/api"
	"auth-service/beans"
	"auth-service/config"
	"auth-service/consul"
	"os"
	"os/signal"
)

func main() {
	if config.ConsulClientEnabled {
		c := consul.NewClient()

		if err := c.Register(); err != nil {
			beans.Logger.Fatal(err)
			return
		}
		defer func() {
			err := c.Unregister()
			if err != nil {
				beans.Logger.Fatal(err)
			}
		}()
	}

	srv := api.NewServer()

	go srv.Run()
	defer srv.Shutdown()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	<-ch
}
