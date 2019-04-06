package discovery

import (
	"auth-service/beans"
	"auth-service/config"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	"os"
	"strconv"
	"time"
)

type Client struct {
	srvId   string
	srvName string
	srvTTL  time.Duration
	agent   *api.Agent
}

func NewClient() *Client {
	cfg := api.DefaultConfig()
	cfg.Address = config.DiscoveryServiceAddress.Value()

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	srvName := "auth-service"
	srvId := fmt.Sprintf("%s:%s",
		srvName, func() string {
			v, _ := uuid.GenerateUUID()
			return v
		}(),
	)

	return &Client{
		srvId:   srvId,
		srvName: srvName,
		srvTTL:  config.HealthCheckTimeout.Value(),
		agent:   client.Agent(),
	}
}

func (c *Client) Register() error {
	err := c.agent.ServiceRegister(&api.AgentServiceRegistration{
		ID:   c.srvId,
		Name: c.srvName,
		Address: func() string {
			v, _ := os.Hostname()
			return v
		}(),
		Port: func() int {
			v, err := strconv.Atoi(config.ServerPort.Value())
			if err != nil {
				panic(err)
			}
			return v
		}(),
		Check: &api.AgentServiceCheck{
			CheckID: "StateCheck",
			TTL:     config.HealthCheckTimeout.Value().String(),
		},
	})
	if err != nil {
		return err
	}

	beans.Logger.Infof("ServiceID: %s", c.srvId)
	beans.Logger.Infof("Consul connection: %s", config.DiscoveryServiceAddress.Value())

	go c.UpdateTTL()

	return nil
}

func (c *Client) Unregister() error {
	if err := c.agent.ServiceDeregister(c.srvId); err != nil {
		return err
	}

	beans.Logger.Infof("Consul connection closed")

	return nil
}

func (c *Client) UpdateTTL() {
	ticker := time.NewTicker(c.srvTTL / 2)
	for range ticker.C {
		_ = c.agent.PassTTL("StateCheck", "UP")
	}
}
