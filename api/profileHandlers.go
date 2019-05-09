package api

import (
	"net/http"

	"github.com/studtool/auth-service/models"
)

func (srv *Server) createProfile(w http.ResponseWriter, r *http.Request) {
	profile := &models.Profile{}
	if err := srv.server.ParseBodyJSON(&profile.Credentials, r); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}
	if err := srv.credentialsValidator.Validate(&profile.Credentials); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	if err := srv.profilesRepository.AddProfile(profile); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	if err := srv.usersQueue.SendUserCreated(profile.UserID); err != nil {
		//TODO find a good way to handle this
		_ = srv.profilesRepository.DeleteProfileById(profile)
		srv.server.WriteErrJSON(w, err)
		return
	}

	srv.server.WriteOkJSON(w, &profile.ProfileInfo)
}

func (srv *Server) updateCredentials(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (srv *Server) deleteProfile(w http.ResponseWriter, r *http.Request) {
	//TODO
}
