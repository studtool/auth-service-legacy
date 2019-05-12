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
	token := &models.Token{}
	if err := srv.server.ParseBodyJSON(token, r); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	if err := srv.tokensRepository.GetToken(token); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	p := &models.ProfileInfo{
		UserID: token.UserID,
	}
	if err := srv.profilesRepository.SetProfileVerified(p); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	createdUserData := &queues.CreatedUserData{
		UserID: p.UserID,
	}
	if err := srv.mqClient.SendUserCreatedMessage(createdUserData); err != nil {
		srv.server.WriteErrJSON(w, err) //TODO handle
		return
	}

	srv.server.WriteOkJSON(w, p)
}

func (srv *Server) updateEmail(w http.ResponseWriter, r *http.Request) {
	emailUpdate := &models.EmailUpdate{
		UserID: srv.parseUserId(r),
	}
	if srv.server.ParseUserID(r) != emailUpdate.UserID {
		srv.server.WriteErrJSON(w, srv.notAuthorizedErr)
		return
	}

	if err := srv.server.ParseBodyJSON(emailUpdate, r); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	if err := srv.profilesRepository.UpdateEmail(emailUpdate); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	srv.server.WriteOk(w)
}

func (srv *Server) updatePassword(w http.ResponseWriter, r *http.Request) {
	passwordUpdate := &models.PasswordUpdate{
		UserID: srv.parseUserId(r),
	}
	if srv.server.ParseUserID(r) != passwordUpdate.UserID {
		srv.server.WriteErrJSON(w, srv.notAuthorizedErr)
		return
	}

	if err := srv.server.ParseBodyJSON(passwordUpdate, r); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	if err := srv.profilesRepository.UpdatePassword(passwordUpdate); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	srv.server.WriteOk(w)
}

func (srv *Server) deleteProfile(w http.ResponseWriter, r *http.Request) {
	srv.server.WriteNotImplemented(w) //TODO
}
