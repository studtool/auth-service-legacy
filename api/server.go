package api

import (
	"fmt"
	"net/http"

	"go.uber.org/dig"

	"github.com/go-http-utils/headers"

	"github.com/go-chi/chi"
	"github.com/gorilla/handlers"

	"github.com/studtool/common/errs"
	"github.com/studtool/common/logs"
	"github.com/studtool/common/rest"

	"github.com/studtool/auth-service/config"
	"github.com/studtool/auth-service/messages"
	"github.com/studtool/auth-service/repositories"
	"github.com/studtool/auth-service/srvutils"
	"github.com/studtool/auth-service/validators"
)

type Server struct {
	rest.Server

	emailValidator       *validators.EmailValidator
	passwordValidator    *validators.PasswordValidator
	credentialsValidator *validators.CredentialsValidator

	authTokenManager    *srvutils.AuthTokenManager
	refreshTokenManager *srvutils.RefreshTokenManager

	notAuthorizedErr *errs.Error
	tokenExpTimeCalc *srvutils.ExpireTimeCalculator

	profilesRepository repositories.ProfilesRepository
	tokensRepository   repositories.TokensRepository
	sessionsRepository repositories.SessionsRepository

	structLogger  logs.Logger
	reflectLogger logs.Logger

	mqClient *messages.MqClient
}

type ServerParams struct {
	dig.In

	MqClient *messages.MqClient

	ProfilesRepository repositories.ProfilesRepository
	TokensRepository   repositories.TokensRepository
	SessionsRepository repositories.SessionsRepository
}

func NewServer(params ServerParams) *Server {
	srv := &Server{
		emailValidator:       validators.NewEmailValidator(),
		passwordValidator:    validators.NewPasswordValidator(),
		credentialsValidator: validators.NewCredentialsValidator(),

		authTokenManager:    srvutils.NewAuthTokenManager(),
		refreshTokenManager: srvutils.NewRefreshTokenManager(),

		tokenExpTimeCalc: srvutils.NewExpireTimeCalculator(),
		notAuthorizedErr: errs.NewNotAuthorizedError("not authorized"),

		profilesRepository: params.ProfilesRepository,
		tokensRepository:   params.TokensRepository,
		sessionsRepository: params.SessionsRepository,

		mqClient: params.MqClient,
	}

	v := rest.ParseAPIVersion(config.ComponentVersion)
	srvPublicPath := rest.MakeAPIPath(v, rest.APITypePublic, "/auth")
	srvProtectedPath := rest.MakeAPIPath(v, rest.APITypeProtected, "/auth")
	srvInternalPath := rest.MakeAPIPath(v, rest.APITypeInternal, "/auth")

	r := chi.NewRouter()

	r.Handle(srvPublicPath+"/profiles", handlers.MethodHandler{
		http.MethodPost: http.HandlerFunc(srv.createProfile),
	})
	r.Handle(srvPublicPath+"/sessions", handlers.MethodHandler{
		http.MethodPost: http.HandlerFunc(srv.startSession),
	})
	r.Handle(srvPublicPath+"/sessions/{session_id}", handlers.MethodHandler{
		http.MethodPatch: http.HandlerFunc(srv.refreshSession),
	})

	r.Handle(srvProtectedPath+"/profiles/{user_id}/email", handlers.MethodHandler{
		http.MethodPatch: srv.WithAuth(http.HandlerFunc(srv.updateEmail)),
	})
	r.Handle(srvProtectedPath+"/profiles/{user_id}/password", handlers.MethodHandler{
		http.MethodPatch: srv.WithAuth(http.HandlerFunc(srv.updatePassword)),
	})
	r.Handle(srvProtectedPath+"/profiles/{user_id}", handlers.MethodHandler{
		http.MethodDelete: srv.WithAuth(http.HandlerFunc(srv.deleteProfile)),
	})
	r.Handle(srvProtectedPath+"/sessions/{session_id}", handlers.MethodHandler{
		http.MethodDelete: srv.WithAuth(http.HandlerFunc(srv.endSession)),
	})
	r.Handle(srvProtectedPath+"/sessions", handlers.MethodHandler{
		http.MethodDelete: srv.WithAuth(http.HandlerFunc(srv.endAllSessions)),
	})

	r.Handle(srvInternalPath+"/session/*", http.HandlerFunc(srv.parseSession))

	r.Handle(`/pprof`, rest.GetProfilerHandler())
	r.Handle(`/metrics`, rest.GetMetricsHandler())

	reqHandler := srv.WithRecover(r)
	if config.RequestsLogsEnabled.Value() {
		reqHandler = srv.WithLogs(reqHandler)
	}
	if config.CorsAllowed.Value() {
		reqHandler = srv.WithCORS(reqHandler, rest.CORS{
			Origins: []string{"*"},
			Methods: []string{
				http.MethodGet, http.MethodHead,
				http.MethodPost, http.MethodPatch,
				http.MethodDelete, http.MethodOptions,
			},
			Headers: []string{
				headers.Authorization, headers.UserAgent,
				headers.ContentType, headers.ContentLength,
				headers.ContentEncoding, headers.ContentLanguage,
			},
			Credentials: false,
		})
	}

	srv.structLogger = srvutils.MakeStructLogger(srv)
	srv.reflectLogger = srvutils.MakeReflectLogger(srv)

	srv.Server = *rest.NewServer(
		rest.ServerParams{
			Address: fmt.Sprintf(":%d", config.ServerPort.Value()),
			Handler: reqHandler,

			StructLogger:  srv.structLogger,
			ReflectLogger: srv.reflectLogger,
			RequestLogger: srvutils.MakeRequestLogger(srv),

			APIClassifier: rest.NewPathAPIClassifier(),
		},
	)
	srv.structLogger.Info("initialized")

	return srv
}
