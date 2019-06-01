package rest

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-http-utils/headers"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/logs"
)

func (srv *Server) WithLogs(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			wr := &LoggingResponseWriter{
				writer: w,
			}

			t := time.Now()
			h.ServeHTTP(wr, r)
			rt := time.Since(t)

			logFunc := srv.requestLogger.Info

			switch wr.severity {
			case SeverityNone:
				break
			case SeverityLow:
				logFunc = srv.requestLogger.Warning
			case SeverityHigh:
				logFunc = srv.requestLogger.Error
			default:
				logFunc = srv.requestLogger.Fatal
			}

			ip := srv.ParseHeaderXRealIP(r)
			if ip == consts.EmptyString {
				ip = r.RemoteAddr
			}

			reqParams := logs.RequestParams{
				Method:      r.Method,
				Path:        r.RequestURI,
				Status:      wr.status,
				UserID:      srv.ParseHeaderUserID(r),
				Type:        srv.apiClassifier.GetType(r),
				IP:          ip,
				UserAgent:   srv.ParseHeaderUserAgent(r),
				RequestTime: rt,
			}

			logFunc(&reqParams)
		},
	)
}

func (srv *Server) WithRecover(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					srv.reflectLogger.Error(r)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			h.ServeHTTP(w, r)
		},
	)
}

func (srv *Server) WithAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if srv.ParseHeaderUserID(r) == consts.EmptyString {
				srv.WriteUnauthorized(w)
				return
			}
			h.ServeHTTP(w, r)
		},
	)
}

type CORS struct {
	Origins     []string
	Methods     []string
	Headers     []string
	Credentials bool
}

func (srv *Server) WithCORS(h http.Handler, cors CORS) http.Handler {
	origin := strings.Join(cors.Origins, ",")
	method := strings.Join(cors.Methods, ",")
	header := strings.Join(cors.Headers, ",")
	credentials := strconv.FormatBool(cors.Credentials)

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(headers.AccessControlAllowOrigin, origin)
			w.Header().Set(headers.AccessControlAllowCredentials, credentials)
			w.Header().Set(headers.AccessControlAllowMethods, method)
			w.Header().Set(headers.AccessControlAllowHeaders, header)
			h.ServeHTTP(w, r)
		},
	)
}
