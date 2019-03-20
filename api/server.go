package api

import (
	"auth-service/config"
	"auth-service/models"
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"net/http"
)

type Server struct {
	server           *http.Server
	profileValidator *models.ProfileValidator
}

func NewServer() *Server {
	srv := &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", config.ServerPort),
		},
		profileValidator: models.NewProfileValidator(),
	}

	mx := http.NewServeMux()
	mx.Handle("/api/auth/profiles", handlers.MethodHandler{
		http.MethodPost:   http.HandlerFunc(srv.createProfile),
		http.MethodPatch:  http.HandlerFunc(srv.updateProfile),
		http.MethodDelete: http.HandlerFunc(srv.deleteProfile),
	})
	mx.Handle("/api/auth/sessions", handlers.MethodHandler{
		http.MethodPost:   http.HandlerFunc(srv.createSession),
		http.MethodPatch:  http.HandlerFunc(srv.updateSession),
		http.MethodDelete: http.HandlerFunc(srv.deleteSession),
	})

	srv.server.Handler = mx
	return srv
}

func (srv *Server) Run() {
	config.Logger.Infof("starting server on %s", srv.server.Addr)
	if err := srv.server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (srv *Server) Shutdown() {
	config.Logger.Info("server shutdown initialized")
	if err := srv.server.Shutdown(context.TODO()); err != nil {
		panic(err)
	}
}
