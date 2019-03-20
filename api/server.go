package api

import (
	"auth-service/config"
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer() *Server {
	srv := &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", config.ServerPort),
		},
	}

	mx := http.NewServeMux()
	mx.Handle("/api/auth/profiles", handlers.MethodHandler{
		http.MethodPost:   http.HandlerFunc(srv.createProfile),
		http.MethodPatch:  http.HandlerFunc(srv.updateProfile),
		http.MethodDelete: http.HandlerFunc(srv.deleteProfile),
	})

	srv.server.Handler = mx
	return srv
}

func (srv *Server) Run() {
	if err := srv.server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (srv *Server) Shutdown() {
	if err := srv.server.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
