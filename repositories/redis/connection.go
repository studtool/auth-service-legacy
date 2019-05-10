package redis

import (
	"github.com/go-redis/redis"

	"github.com/studtool/common/consts"
)

type Connection struct {
	client *redis.Client
}

func NewConnection() *Connection {
	return &Connection{
		client: redis.NewClient(&redis.Options{
			Addr:     "127.0.0.1:6379", //TODO
			Password: consts.EmptyString,
			DB:       0,
		}),
	}
}

func (conn *Connection) Open() (err error) {
	return nil //TODO
}

func (conn *Connection) Close() (err error) {
	return nil //TODO
}
