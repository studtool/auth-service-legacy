package consul

import (
	"auth-service/config"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-uuid"
	"os"
	"strconv"
)

type Client struct {
	srvId string
	agent *api.Agent
}

func NewClient() *Client {
	cfg := api.DefaultConfig()

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	agent := client.Agent()

	id, _ := uuid.GenerateUUID()

	return &Client{
		srvId: id,
		agent: agent,
	}
}

func (c *Client) Register() {
	err := c.agent.ServiceRegister(&api.AgentServiceRegistration{
		Name: "auth",
		ID:   c.srvId,
		Address: func() string {
			v, _ := os.Hostname()
			return v
		}(),
		Port: func() int {
			v, err := strconv.Atoi(config.ServerPort)
			if err != nil {
				panic(err)
			}
			return v
		}(),
	})
	if err != nil {
		panic(err)
	}
}

func (c *Client) Unregister() {
	if err := c.agent.ServiceDeregister(c.srvId); err != nil {
		panic(err)
	}
}
