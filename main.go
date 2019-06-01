package main

import (
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/dig"

	"github.com/studtool/common/logs"
	"github.com/studtool/common/utils/assertions"

	"github.com/studtool/auth-service/api"
	"github.com/studtool/auth-service/config"
	"github.com/studtool/auth-service/messages"
	"github.com/studtool/auth-service/repositories"
	"github.com/studtool/auth-service/repositories/postgres"
	"github.com/studtool/auth-service/repositories/redis"
)

func main() {
	c := dig.New()
	logger := logs.NewReflectLogger()

	if config.RepositoriesEnabled.Value() {
		assertions.AssertOk(c.Provide(postgres.NewConnection))
		assertions.AssertOk(c.Invoke(func(conn *postgres.Connection) {
			if err := conn.Open(); err != nil {
				logger.Fatal(err)
			}
		}))
		defer func() {
			assertions.AssertOk(c.Invoke(func(conn *postgres.Connection) {
				if err := conn.Close(); err != nil {
					logger.Fatal(err)
				}
			}))
		}()

		assertions.AssertOk(c.Provide(
			postgres.NewProfilesRepository,
			dig.As(new(repositories.ProfilesRepository)),
		))
		assertions.AssertOk(c.Provide(
			postgres.NewSessionsRepository,
			dig.As(new(repositories.SessionsRepository)),
		))

		assertions.AssertOk(c.Provide(redis.NewConnection))
		assertions.AssertOk(c.Invoke(func(conn *redis.Connection) {
			if err := conn.Open(); err != nil {
				logger.Fatal(err)
			}
		}))
		defer func() {
			assertions.AssertOk(c.Invoke(func(conn *redis.Connection) {
				if err := conn.Close(); err != nil {
					logger.Fatal(err)
				}
			}))
		}()

		assertions.AssertOk(c.Provide(
			redis.NewTokensRepository,
			dig.As(new(repositories.TokensRepository)),
		))
	} else {
		assertions.AssertOk(c.Provide(
			func() repositories.ProfilesRepository {
				return nil
			},
		))
		assertions.AssertOk(c.Provide(
			func() repositories.SessionsRepository {
				return nil
			},
		))
		assertions.AssertOk(c.Provide(
			func() repositories.TokensRepository {
				return nil
			},
		))
	}

	if config.QueuesEnabled.Value() {
		assertions.AssertOk(c.Provide(messages.NewMqClient))
		assertions.AssertOk(c.Invoke(func(q *messages.MqClient) {
			if err := q.OpenConnection(); err != nil {
				logger.Fatal(err)
			}
		}))
		defer func() {
			assertions.AssertOk(c.Invoke(func(q *messages.MqClient) {
				if err := q.CloseConnection(); err != nil {
					logger.Fatal(err)
				}
			}))
		}()
	} else {
		assertions.AssertOk(c.Provide(func() *messages.MqClient {
			return nil
		}))
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGTERM)

	assertions.AssertOk(c.Provide(api.NewServer))
	assertions.AssertOk(c.Invoke(func(srv *api.Server) {
		go func() {
			if err := srv.Run(); err != nil {
				logger.Fatal(err)
				ch <- os.Interrupt
			}
		}()
	}))
	defer func() {
		assertions.AssertOk(c.Invoke(func(srv *api.Server) {
			if err := srv.Shutdown(); err != nil {
				logger.Fatal(err)
			}
		}))
	}()

	<-ch
}
