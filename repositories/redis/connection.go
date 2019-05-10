package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/studtool/auth-service/config"

	"github.com/studtool/common/consts"
)

type Connection struct {
	client *redis.Client
}

func NewConnection() *Connection {
	return &Connection{
		client: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%d",
				config.TokensStorageHost.Value(), config.TokensStoragePort.Value(),
			),
			Password: consts.EmptyString,
		}),
	}
}

func (conn *Connection) Open() (err error) {
	return nil //TODO
}

func (conn *Connection) Close() (err error) {
	return nil //TODO
}
