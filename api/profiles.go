package api

import (
	"net/http"

	"github.com/studtool/auth-service/models"
)

func (srv *Server) createProfile(w http.ResponseWriter, r *http.Request) {
	profile := &models.Profile{}
	if err := srv.server.ParseBodyJSON(profile, r); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	if err := srv.profileValidator.ValidateOnCreate(profile); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	if err := srv.profilesRepository.AddProfile(profile); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	_ = srv.usersQueue.SendUserCreated(profile.UserId) //TODO handle error

	srv.server.WriteOk(w)
}

func (srv *Server) updateCredentials(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (srv *Server) deleteProfile(w http.ResponseWriter, r *http.Request) {
	//TODO
}
