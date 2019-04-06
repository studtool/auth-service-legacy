package main

import (
	"auth-service/api"
	"auth-service/beans"
	"auth-service/config"
	"auth-service/discovery"
	"auth-service/mq"
	"auth-service/repositories"
	"auth-service/repositories/postgres"
	"go.uber.org/dig"
	"os"
	"os/signal"
)

func main() {
	c := dig.New()

	_ = c.Provide(postgres.NewConnection)
	_ = c.Provide(
		postgres.NewProfilesRepository,
		dig.As(new(repositories.ProfilesRepository)),
	)
	_ = c.Provide(
		postgres.NewSessionsRepository,
		dig.As(new(repositories.SessionsRepository)),
	)
	_ = c.Provide(mq.NewQueue)
	_ = c.Provide(discovery.NewClient)
	_ = c.Provide(api.NewServer)

	if config.RepositoriesEnabled.Value() {
		_ = c.Invoke(func(conn *postgres.Connection) {
			if err := conn.Open(); err != nil {
				beans.Logger.Fatal(err)
			} else {
				beans.Logger.Info("storage connection: open")
			}
		})
		defer func() {
			_ = c.Invoke(func(conn *postgres.Connection) {
				if err := conn.Close(); err != nil {
					beans.Logger.Fatal(err)
				} else {
					beans.Logger.Info("connection to storage: closed")
				}
			})
		}()
	}

	if config.QueuesEnabled.Value() {
		_ = c.Invoke(func(q *mq.MQ) {
			if err := q.OpenConnection(); err != nil {
				beans.Logger.Fatal(err)
			} else {
				beans.Logger.Info("message queue connection: open")
			}
		})
		defer func() {
			_ = c.Invoke(func(q *mq.MQ) {
				if err := q.CloseConnection(); err != nil {
					beans.Logger.Fatal(err)
				} else {
					beans.Logger.Info("message queue connection: closed")
				}
			})
		}()
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	_ = c.Invoke(func(srv *api.Server) {
		go func() {
			if err := srv.Run(); err != nil {
				beans.Logger.Fatal(err)
				ch <- os.Interrupt
			}
		}()

		beans.Logger.Infof("server: started; [port: %s]", config.ServerPort.Value())
	})
	defer func() {
		_ = c.Invoke(func(srv *api.Server) {
			if err := srv.Shutdown(); err != nil {
				beans.Logger.Fatal(err)
			} else {
				beans.Logger.Info("server: down")
			}
		})
	}()

	if config.DiscoveryClientEnabled.Value() {
		_ = c.Invoke(func(cl *discovery.Client) {
			if err := cl.Register(); err != nil {
				beans.Logger.Fatal(err)
			} else {
				beans.Logger.Info("discovery client connection: open")
			}
		})
		defer func() {
			_ = c.Invoke(func(cl *discovery.Client) {
				if err := cl.Unregister(); err != nil {
					beans.Logger.Fatal(err)
				} else {
					beans.Logger.Info("discovery client connection: closed")
				}
			})
		}()
	}

	<-ch
}
