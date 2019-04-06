package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (srv *Server) parseProfileId(r *http.Request) string {
	return mux.Vars(r)["profile_id"]
}
