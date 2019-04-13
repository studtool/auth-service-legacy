package api

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/studtool/auth-service/beans"
	"github.com/studtool/auth-service/mq"
	"net/http"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/errs"
	"github.com/studtool/common/rest"

	"github.com/studtool/auth-service/config"
	"github.com/studtool/auth-service/models"
	"github.com/studtool/auth-service/repositories"
	"github.com/studtool/auth-service/utils"
)

const (
	idPattern = `\w{8}-\w{4}-\w{4}-\w{4}-\w{12}`
)

type Server struct {
	server *rest.Server

	profileValidator *models.ProfileValidator

	authTokenManager    *utils.AuthTokenManager
	refreshTokenManager *utils.RefreshTokenManager

	noAuthTokenErr *errs.Error

	profilesRepository repositories.ProfilesRepository
	sessionsRepository repositories.SessionsRepository

	usersQueue *mq.MQ
}

func NewServer(pRepo repositories.ProfilesRepository,
	sRepo repositories.SessionsRepository, uQueue *mq.MQ) *Server {

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
	mx.Handle(`/api/auth/profiles/{profile_id:`+idPattern+`}/credentials`, handlers.MethodHandler{
		http.MethodPatch: srv.server.WithAuth(http.HandlerFunc(srv.updateCredentials)),
	})
	mx.Handle(`/api/auth/profiles/{profile_id:`+idPattern+`}`, handlers.MethodHandler{
		http.MethodDelete: srv.server.WithAuth(http.HandlerFunc(srv.deleteProfile)),
	})
	mx.Handle(`/api/auth/sessions`, handlers.MethodHandler{
		http.MethodPost:   http.HandlerFunc(srv.startSession),
		http.MethodDelete: http.HandlerFunc(srv.endAllSessions),
	})
	mx.Handle(`/api/auth/session`, handlers.MethodHandler{
		http.MethodGet:    http.HandlerFunc(srv.parseSession),
		http.MethodPatch:  http.HandlerFunc(srv.refreshSession),
		http.MethodDelete: http.HandlerFunc(srv.endSession),
	})

	srv.server.SetLogger(beans.Logger)
	srv.server.SetHandler(srv.server.WithLogs(srv.server.WithRecover(mx)))

	return srv
}

func (srv *Server) Run() error {
	return srv.server.Run()
}

func (srv *Server) Shutdown() error {
	return srv.server.Shutdown()
}
