package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/studtool/common/consts"
)

func (srv *Server) WithLogs(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			wr := &LoggingResponseWriter{
				writer: w,
			}
			h.ServeHTTP(wr, r)

			srv.logger.Info(
				fmt.Sprintf("%s %s %d", r.Method, r.RequestURI, wr.status),
			)
		},
	)
}

func (srv *Server) WithRecover(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					srv.logger.Error(fmt.Sprintf("panic: %v", r))

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
			if srv.ParseUserID(r) == consts.EmptyString {
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
	origins := strings.Join(cors.Origins, ",")
	methods := strings.Join(cors.Methods, ",")
	headers := strings.Join(cors.Headers, ",")
	credentials := strconv.FormatBool(cors.Credentials)

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", origins)
			w.Header().Set("Access-Control-Allow-Credentials", credentials)
			w.Header().Set("Access-Control-Allow-Methods", methods)
			w.Header().Set("Access-Control-Allow-Headers", headers)
			h.ServeHTTP(w, r)
		},
	)
}
