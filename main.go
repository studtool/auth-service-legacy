package main

import (
	"auth-service/consul"
)

func main() {
	c := consul.NewClient()

	c.Register()
	defer c.Unregister()
}
