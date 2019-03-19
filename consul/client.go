package consul

import (
	"github.com/hashicorp/consul/api"
)

type Client struct {
	client *api.Client
}

func NewClient() *Client {
	config := api.DefaultConfig()

	client, err := api.NewClient(config)
	if err != nil {
		panic(err)
	}

	return &Client{
		client: client,
	}
}

func (c *Client) Register() {
	//TODO
}

func (c *Client) Unregister() {
	//TODO
}
