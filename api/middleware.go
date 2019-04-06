package api

import (
	"fmt"
	"net/http"

	"github.com/studtool/auth-service/beans"
)

func (srv *Server) withRecover(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					beans.Logger.Error(fmt.Sprintf("panic: %v", r))
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			h.ServeHTTP(w, r)
		},
	)
}

func (srv *Server) withAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userId := srv.parseUserId(r)
			if userId == "" {
				w.WriteHeader(http.StatusUnauthorized)
			}
			h.ServeHTTP(w, r)
		},
	)
}
