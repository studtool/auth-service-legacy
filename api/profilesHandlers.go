package api

import (
	"net/http"

	"github.com/studtool/common/queues"

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

	token := &models.Token{
		UserID: profile.UserID,
	}
	if err := srv.tokensRepository.SetToken(token); err != nil {
		srv.server.WriteErrJSON(w, err) //TODO handle
		return
	}

	regEmailData := &queues.RegistrationEmailData{
		Email: profile.Email,
		Token: token.Token,
	}
	if err := srv.mqClient.SendRegEmailMessage(regEmailData); err != nil {
		srv.server.WriteErrJSON(w, err) //TODO handle
		return
	}

	srv.server.WriteOkJSON(w, &profile.ProfileInfo)
}

func (srv *Server) verifyProfile(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (srv *Server) updateEmail(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (srv *Server) updatePassword(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (srv *Server) deleteProfile(w http.ResponseWriter, r *http.Request) {
	//TODO
}
