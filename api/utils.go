package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (srv *Server) parseUserID(r *http.Request) string {
	return mux.Vars(r)["user_id"]
}

func (srv *Server) parseSessionID(r *http.Request) string {
	return mux.Vars(r)["session_id"]
}
