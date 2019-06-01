package api

import (
	"net/http"

	"github.com/go-chi/chi"

	"github.com/studtool/common/types"
)

func (srv *Server) parseHeaderAuthToken(r *http.Request) string {
	return srv.ParseHeaderAuthToken(r)
}

func (srv *Server) parseHeaderUserID(r *http.Request) types.ID {
	return srv.ParseHeaderUserID(r)
}

func (srv *Server) parseHeaderRefreshToken(r *http.Request) string {
	return srv.ParseHeaderRefreshToken(r)
}

func (srv *Server) parsePathUserID(r *http.Request) string {
	return chi.URLParam(r, "user_id")
}

func (srv *Server) parsePathSessionID(r *http.Request) string {
	return chi.URLParam(r, "session_id")
}
