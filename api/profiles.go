package api

import (
	"auth-service/models"
	"net/http"
)

func (srv *Server) createProfile(w http.ResponseWriter, r *http.Request) {
	var profile models.Profile
	if err := srv.parseRequestBody(&profile, r); err != nil {
		srv.writeErrJSON(w, err)
		return
	}

	srv.writeOk(w)
}

func (srv *Server) updateProfile(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (srv *Server) deleteProfile(w http.ResponseWriter, r *http.Request) {
	//TODO
}
