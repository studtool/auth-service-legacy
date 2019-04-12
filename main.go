package main

import (
	"go.uber.org/dig"
	"os"
	"os/signal"

	"github.com/studtool/auth-service/api"
	"github.com/studtool/auth-service/beans"
	"github.com/studtool/auth-service/config"
	"github.com/studtool/auth-service/discovery"
	"github.com/studtool/auth-service/mq"
	"github.com/studtool/auth-service/repositories"
	"github.com/studtool/auth-service/repositories/postgres"
)

func main() {
	c := dig.New()

	panicOnErr(c.Provide(postgres.NewConnection))
	panicOnErr(c.Provide(
		postgres.NewProfilesRepository,
		dig.As(new(repositories.ProfilesRepository)),
	))
	panicOnErr(c.Provide(
		postgres.NewSessionsRepository,
		dig.As(new(repositories.SessionsRepository)),
	))
	panicOnErr(c.Provide(mq.NewQueue))
	panicOnErr(c.Provide(discovery.NewClient))
	panicOnErr(c.Provide(api.NewServer))

	if config.RepositoriesEnabled.Value() {
		panicOnErr(c.Invoke(func(conn *postgres.Connection) {
			if err := conn.Open(); err != nil {
				beans.Logger.Fatal(err)
			} else {
				beans.Logger.Info("storage connection: open")
			}
		}))
		defer func() {
			panicOnErr(c.Invoke(func(conn *postgres.Connection) {
				if err := conn.Close(); err != nil {
					beans.Logger.Fatal(err)
				} else {
					beans.Logger.Info("connection to storage: closed")
				}
			}))
		}()
	}

	if config.QueuesEnabled.Value() {
		panicOnErr(c.Invoke(func(q *mq.MQ) {
			if err := q.OpenConnection(); err != nil {
				beans.Logger.Fatal(err)
			} else {
				beans.Logger.Info("message queue connection: open")
			}
		}))
		defer func() {
			panicOnErr(c.Invoke(func(q *mq.MQ) {
				if err := q.CloseConnection(); err != nil {
					beans.Logger.Fatal(err)
				} else {
					beans.Logger.Info("message queue connection: closed")
				}
			}))
		}()
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	panicOnErr(c.Invoke(func(srv *api.Server) {
		go func() {
			if err := srv.Run(); err != nil {
				beans.Logger.Fatal(err)
				ch <- os.Interrupt
			}
		}()
	}))
	defer func() {
		panicOnErr(c.Invoke(func(srv *api.Server) {
			if err := srv.Shutdown(); err != nil {
				beans.Logger.Fatal(err)
			} else {
				beans.Logger.Info("server: down")
			}
		}))
	}()

	if config.DiscoveryClientEnabled.Value() {
		panicOnErr(c.Invoke(func(cl *discovery.Client) {
			if err := cl.Register(); err != nil {
				beans.Logger.Fatal(err)
			} else {
				beans.Logger.Info("discovery client connection: open")
			}
		}))
		defer func() {
			panicOnErr(c.Invoke(func(cl *discovery.Client) {
				if err := cl.Unregister(); err != nil {
					beans.Logger.Fatal(err)
				} else {
					beans.Logger.Info("discovery client connection: closed")
				}
			}))
		}()
	}

	<-ch
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
