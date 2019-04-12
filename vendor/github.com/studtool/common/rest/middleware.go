package rest

import (
	"fmt"
	"net/http"
)

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
			userId := srv.ParseUserId(r)
			if userId == "" {
				w.WriteHeader(http.StatusUnauthorized)
			}
			h.ServeHTTP(w, r)
		},
	)
}
