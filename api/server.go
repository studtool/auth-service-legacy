package api

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/repositories"
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"net/http"
)

const (
	idPattern = `\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`
)

type Server struct {
	server             *http.Server
	profileValidator   *models.ProfileValidator
	profilesRepository repositories.ProfilesRepository
}

func NewServer(pRepo repositories.ProfilesRepository) *Server {
	srv := &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", config.ServerPort),
		},
		profileValidator:   models.NewProfileValidator(),
		profilesRepository: pRepo,
	}

	mx := http.NewServeMux()
	mx.Handle(`/api/auth/profiles`, handlers.MethodHandler{
		http.MethodPost:  http.HandlerFunc(srv.createProfile),
		http.MethodPatch: http.HandlerFunc(srv.updateProfile),
	})
	mx.Handle(`/api/auth/profiles/{id:`+idPattern+`}`, handlers.MethodHandler{
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

func (srv *Server) Run() error {
	return srv.server.ListenAndServe()
}

func (srv *Server) Shutdown() error {
	return srv.server.Shutdown(context.TODO())
}
