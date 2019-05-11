package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (srv *Server) parseUserId(r *http.Request) string {
	return mux.Vars(r)["user_id"]
}
