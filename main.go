package main

import (
	"os"
	"os/signal"

	"go.uber.org/dig"

	"github.com/studtool/common/utils"

	"github.com/studtool/auth-service/api"
	"github.com/studtool/auth-service/beans"
	"github.com/studtool/auth-service/config"
	"github.com/studtool/auth-service/mq"
	"github.com/studtool/auth-service/repositories"
	"github.com/studtool/auth-service/repositories/postgres"
)

func main() {
	c := dig.New()

	if config.RepositoriesEnabled.Value() {
		utils.AssertOk(c.Provide(postgres.NewConnection))
		utils.AssertOk(c.Invoke(func(conn *postgres.Connection) {
			if err := conn.Open(); err != nil {
				beans.Logger.Fatal(err)
			}
		}))
		defer func() {
			utils.AssertOk(c.Invoke(func(conn *postgres.Connection) {
				if err := conn.Close(); err != nil {
					beans.Logger.Fatal(err)
				}
			}))
		}()

		utils.AssertOk(c.Provide(
			postgres.NewProfilesRepository,
			dig.As(new(repositories.ProfilesRepository)),
		))
		utils.AssertOk(c.Provide(
			postgres.NewSessionsRepository,
			dig.As(new(repositories.SessionsRepository)),
		))
	} else {
		utils.AssertOk(c.Provide(
			func() repositories.ProfilesRepository {
				return nil
			},
		))
		utils.AssertOk(c.Provide(
			func() repositories.SessionsRepository {
				return nil
			},
		))
	}

	if config.QueuesEnabled.Value() {
		utils.AssertOk(c.Provide(mq.NewClient))
		utils.AssertOk(c.Invoke(func(q *mq.Client) {
			if err := q.OpenConnection(); err != nil {
				beans.Logger.Fatal(err)
			}
		}))
		defer func() {
			utils.AssertOk(c.Invoke(func(q *mq.Client) {
				if err := q.CloseConnection(); err != nil {
					beans.Logger.Fatal(err)
				}
			}))
		}()
	} else {
		utils.AssertOk(c.Provide(func() *mq.Client {
			return nil
		}))
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill)
	signal.Notify(ch, os.Interrupt)

	utils.AssertOk(c.Provide(api.NewServer))
	utils.AssertOk(c.Invoke(func(srv *api.Server) {
		go func() {
			if err := srv.Run(); err != nil {
				beans.Logger.Fatal(err)
				ch <- os.Interrupt
			}
		}()
	}))
	defer func() {
		utils.AssertOk(c.Invoke(func(srv *api.Server) {
			if err := srv.Shutdown(); err != nil {
				beans.Logger.Fatal(err)
			}
		}))
	}()

	<-ch
}
