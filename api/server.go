package api

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/repositories"
	"context"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/studtool/common/jwt"
	"net/http"
)

const (
	idPattern = `\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`
)

type Server struct {
	server             *http.Server
	profileValidator   *models.ProfileValidator
	jwtManager         *jwt.Manager
	profilesRepository repositories.ProfilesRepository
	sessionsRepository repositories.SessionsRepository
}

func NewServer(
	pR repositories.ProfilesRepository,
	sR repositories.SessionsRepository,
) *Server {
	srv := &Server{
		server: &http.Server{
			Addr: fmt.Sprintf(":%s", config.ServerPort),
		},
		profileValidator:   models.NewProfileValidator(),
		profilesRepository: pR,
		sessionsRepository: sR,
	}

	mx := mux.NewRouter()
	mx.Handle(`/api/auth/profiles`, handlers.MethodHandler{
		http.MethodPost: http.HandlerFunc(srv.createProfile),
	})
	mx.Handle(`/api/auth/profiles/{profile_id:`+idPattern+`}/credentials`, handlers.MethodHandler{
		http.MethodPatch: srv.withAuth(http.HandlerFunc(srv.updateCredentials)),
	})
	mx.Handle(`/api/auth/profiles/{profile_id:`+idPattern+`}`, handlers.MethodHandler{
		http.MethodDelete: srv.withAuth(http.HandlerFunc(srv.deleteProfile)),
	})
	mx.Handle(`/api/auth/sessions`, handlers.MethodHandler{
		http.MethodPost:   http.HandlerFunc(srv.startSession),
		http.MethodPatch:  http.HandlerFunc(srv.refreshSession),
		http.MethodDelete: srv.withAuth(http.HandlerFunc(srv.endSession)),
	})
	mx.Handle(`/api/auth/sessions/{profile_id:`+idPattern+`}`, handlers.MethodHandler{
		http.MethodDelete: srv.withAuth(http.HandlerFunc(srv.endAllSessions)),
	})

	srv.server.Handler = srv.withRecover(mx)
	return srv
}

func (srv *Server) Run() error {
	return srv.server.ListenAndServe()
}

func (srv *Server) Shutdown() error {
	return srv.server.Shutdown(context.TODO())
}
