package main

import (
	"auth-service/api"
	"auth-service/beans"
	"auth-service/config"
	"auth-service/consul"
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
	_ = c.Provide(func(conn *postgres.Connection) repositories.ProfilesRepository {
		return postgres.NewProfilesRepository(conn)
	})
	_ = c.Provide(mq.NewQueue)
	_ = c.Provide(consul.NewClient)
	_ = c.Provide(api.NewServer)

	_ = c.Invoke(func(conn *postgres.Connection) {
		if err := conn.Open(); err != nil {
			beans.Logger.Fatal(err)
		}
	})
	defer func() {
		_ = c.Invoke(func(conn *postgres.Connection) {
			if err := conn.Close(); err != nil {
				beans.Logger.Fatal(err)
			}
		})
	}()

	if config.ShouldInitStorage {
		_ = c.Invoke(func(r *postgres.ProfilesRepository) {
			if err := r.Init(); err != nil {
				beans.Logger.Fatal(err)
			}
		})
	}

	_ = c.Invoke(func(q *mq.MQ) {
		if err := q.OpenConnection(); err != nil {
			beans.Logger.Fatal(err)
		}
	})
	defer func() {
		_ = c.Invoke(func(q *mq.MQ) {
			if err := q.CloseConnection(); err != nil {
				beans.Logger.Fatal(err)
			}
		})
	}()

	_ = c.Invoke(func(srv *api.Server) {
		if err := srv.Run(); err != nil {
			beans.Logger.Fatal(err)
		}
	})
	defer func() {
		_ = c.Invoke(func(srv *api.Server) {
			if err := srv.Shutdown(); err != nil {
				beans.Logger.Fatal(err)
			}
		})
	}()

	if config.ConsulClientEnabled {
		_ = c.Invoke(func(cl *consul.Client) {
			if err := cl.Register(); err != nil {
				beans.Logger.Fatal(err)
			}
		})
		defer func() {
			_ = c.Invoke(func(cl *consul.Client) {
				if err := cl.Unregister(); err != nil {
					beans.Logger.Fatal(err)
				}
			})
		}()
	}

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)

	<-ch
}
