package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (srv *Server) parseUserID(r *http.Request) string {
	return chi.URLParam(r, "user_id")
}

func (srv *Server) parseSessionID(r *http.Request) string {
	return chi.URLParam(r, "session_id")
}
