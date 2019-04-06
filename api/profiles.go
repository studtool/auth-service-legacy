package api

import (
	"net/http"

	"github.com/studtool/auth-service/models"
)

func (srv *Server) createProfile(w http.ResponseWriter, r *http.Request) {
	profile := &models.Profile{}
	if err := srv.parseRequestBody(profile, r); err != nil {
		srv.writeErrJSON(w, err)
		return
	}

	if err := srv.profileValidator.ValidateOnCreate(profile); err != nil {
		srv.writeErrJSON(w, err)
		return
	}

	if err := srv.profilesRepository.AddProfile(profile); err != nil {
		srv.writeErrJSON(w, err)
		return
	}

	srv.writeOk(w)
}

func (srv *Server) updateCredentials(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (srv *Server) deleteProfile(w http.ResponseWriter, r *http.Request) {
	//TODO
}
