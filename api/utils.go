package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (srv *Server) parseProfileId(r *http.Request) string {
	return mux.Vars(r)["profile_id"]
}
