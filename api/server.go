package api

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/errs"
	"github.com/studtool/common/rest"

	"github.com/studtool/auth-service/beans"
	"github.com/studtool/auth-service/config"
	"github.com/studtool/auth-service/models"
	"github.com/studtool/auth-service/mq"
	"github.com/studtool/auth-service/repositories"
	"github.com/studtool/auth-service/utils"
)

type Server struct {
	server *rest.Server

	profileValidator *models.ProfileValidator

	authTokenManager    *utils.AuthTokenManager
	refreshTokenManager *utils.RefreshTokenManager

	noAuthTokenErr *errs.Error

	profilesRepository repositories.ProfilesRepository
	sessionsRepository repositories.SessionsRepository

	usersQueue *mq.Client
}

func NewServer(pRepo repositories.ProfilesRepository,
	sRepo repositories.SessionsRepository, uQueue *mq.Client) *Server {

	srv := &Server{
		server: rest.NewServer(
			rest.ServerConfig{
				Host: consts.EmptyString,
				Port: config.ServerPort.Value(),
			},
		),

		profileValidator: models.NewProfileValidator(),

		authTokenManager:    utils.NewAuthTokenManager(),
		refreshTokenManager: utils.NewRefreshTokenManager(),

		noAuthTokenErr: errs.NewNotAuthorizedError("authorization token required"),

		profilesRepository: pRepo,
		sessionsRepository: sRepo,

		usersQueue: uQueue,
	}

	mx := mux.NewRouter()
	mx.Handle(`/api/auth/profiles`, handlers.MethodHandler{
		http.MethodPost: http.HandlerFunc(srv.createProfile),
	})
	mx.Handle(`/api/auth/profiles/{profile_id}/credentials`, handlers.MethodHandler{
		http.MethodPatch: srv.server.WithAuth(http.HandlerFunc(srv.updateCredentials)),
	})
	mx.Handle(`/api/auth/profiles/{profile_id}`, handlers.MethodHandler{
		http.MethodDelete: srv.server.WithAuth(http.HandlerFunc(srv.deleteProfile)),
	})
	mx.Handle(`/api/auth/sessions`, handlers.MethodHandler{
		http.MethodPost:   http.HandlerFunc(srv.startSession),
		http.MethodDelete: http.HandlerFunc(srv.endAllSessions),
	})
	mx.Handle(`/api/auth/session`, handlers.MethodHandler{
		http.MethodPatch:  http.HandlerFunc(srv.refreshSession),
		http.MethodDelete: http.HandlerFunc(srv.endSession),
	})
	mx.Handle(`/api/private/auth/session`, handlers.MethodHandler{
		http.MethodGet: http.HandlerFunc(srv.parseSession),
	})

	srv.server.SetLogger(beans.Logger)

	h := srv.server.WithRecover(mx)
	if config.RequestsLogsEnabled.Value() {
		h = srv.server.WithLogs(h)
	}
	if config.CorsAllowed.Value() {
		h = srv.server.WithCORS(h, rest.CORS{
			Origins: []string{"*"},
			Methods: []string{
				http.MethodGet, http.MethodHead,
				http.MethodPost, http.MethodPatch,
				http.MethodDelete, http.MethodOptions,
			},
			Headers: []string{
				"User-Agent", "Authorization",
				"Content-Type", "Content-Length",
			},
			Credentials: false,
		})
	}

	srv.server.SetHandler(h)

	return srv
}

func (srv *Server) Run() error {
	return srv.server.Run()
}

func (srv *Server) Shutdown() error {
	return srv.server.Shutdown()
}
