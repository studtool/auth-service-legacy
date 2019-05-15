package api

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"go.uber.org/dig"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/errs"
	"github.com/studtool/common/rest"

	"github.com/studtool/auth-service/beans"
	"github.com/studtool/auth-service/config"
	"github.com/studtool/auth-service/messages"
	"github.com/studtool/auth-service/repositories"
	"github.com/studtool/auth-service/utils"
	"github.com/studtool/auth-service/validators"
)

type Server struct {
	server *rest.Server

	emailValidator       *validators.EmailValidator
	passwordValidator    *validators.PasswordValidator
	credentialsValidator *validators.CredentialsValidator

	authTokenManager    *utils.AuthTokenManager
	refreshTokenManager *utils.RefreshTokenManager

	notAuthorizedErr *errs.Error
	tokenExpTimeCalc *utils.ExpireTimeCalculator

	profilesRepository repositories.ProfilesRepository
	tokensRepository   repositories.TokensRepository
	sessionsRepository repositories.SessionsRepository

	mqClient *messages.QueueClient
}

type ServerParams struct {
	dig.In

	MqClient *messages.QueueClient

	ProfilesRepository repositories.ProfilesRepository
	TokensRepository   repositories.TokensRepository
	SessionsRepository repositories.SessionsRepository
}

func NewServer(params ServerParams) *Server {
	srv := &Server{
		server: rest.NewServer(
			rest.ServerConfig{
				Host: consts.EmptyString,
				Port: config.ServerPort.Value(),
			},
		),

		emailValidator:       validators.NewEmailValidator(),
		passwordValidator:    validators.NewPasswordValidator(),
		credentialsValidator: validators.NewCredentialsValidator(),

		authTokenManager:    utils.NewAuthTokenManager(),
		refreshTokenManager: utils.NewRefreshTokenManager(),

		notAuthorizedErr: errs.NewNotAuthorizedError("not authorized"),
		tokenExpTimeCalc: utils.NewExpireTimeCalculator(),

		profilesRepository: params.ProfilesRepository,
		tokensRepository:   params.TokensRepository,
		sessionsRepository: params.SessionsRepository,

		mqClient: params.MqClient,
	}

	mx := mux.NewRouter()
	mx.Handle(`/api/auth/profiles`, handlers.MethodHandler{
		http.MethodPost: http.HandlerFunc(srv.createProfile),
	})
	mx.Handle(`/api/auth/profiles/{user_id}`, handlers.MethodHandler{
		http.MethodPatch: http.HandlerFunc(srv.verifyProfile),
	})
	mx.Handle(`/api/auth/profiles/{user_id}/email`, handlers.MethodHandler{
		http.MethodPatch: srv.server.WithAuth(http.HandlerFunc(srv.updateEmail)),
	})
	mx.Handle(`/api/auth/profiles/{user_id}/password`, handlers.MethodHandler{
		http.MethodPatch: srv.server.WithAuth(http.HandlerFunc(srv.updatePassword)),
	})
	mx.Handle(`/api/auth/profiles/{user_id}`, handlers.MethodHandler{
		http.MethodDelete: srv.server.WithAuth(http.HandlerFunc(srv.deleteProfile)),
	})
	mx.Handle(`/api/auth/sessions`, handlers.MethodHandler{
		http.MethodPost:   http.HandlerFunc(srv.startSession),
		http.MethodDelete: srv.server.WithAuth(http.HandlerFunc(srv.endAllSessions)),
	})
	mx.Handle(`/api/auth/sessions/{session_id}`, handlers.MethodHandler{
		http.MethodPatch:  http.HandlerFunc(srv.refreshSession),
		http.MethodDelete: srv.server.WithAuth(http.HandlerFunc(srv.endSession)),
	})
	mx.Handle(`/api/private/auth/session`, http.HandlerFunc(srv.parseSession))

	srv.server.SetLogger(beans.Logger())

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
