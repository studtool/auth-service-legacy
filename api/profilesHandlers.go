package api

import (
	"net/http"

	"github.com/studtool/common/queues"

	"github.com/studtool/auth-service/beans"
	"github.com/studtool/auth-service/models"
)

func (srv *Server) createProfile(w http.ResponseWriter, r *http.Request) {
	profile := &models.Profile{}
	if err := srv.ParseBodyJSON(&profile.Credentials, r); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}
	if err := srv.credentialsValidator.Validate(&profile.Credentials); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}

	if err := srv.profilesRepository.AddProfile(profile); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}

	createdUserData := &queues.CreatedUserData{
		UserID: profile.UserID,
	}
	if err := srv.mqClient.PublishUserCreatedMessage(createdUserData); err != nil {
		srv.WriteErrJSON(w, err) //TODO handle
		return
	}

	token := &models.Token{
		UserID: profile.UserID,
	}
	if err := srv.tokensRepository.SetToken(token); err != nil {
		srv.WriteErrJSON(w, err) //TODO handle
		return
	}

	regEmailData := &queues.RegistrationEmailData{
		Email: profile.Email,
		Token: token.Token,
	}
	if err := srv.mqClient.PublishRegistrationEmailMessage(regEmailData); err != nil {
		srv.WriteErrJSON(w, err) //TODO handle
		return
	}

	srv.WriteOkJSON(w, &profile.ProfileInfo)
}

func (srv *Server) verifyProfile(w http.ResponseWriter, r *http.Request) {
	token := &models.Token{}
	if err := srv.ParseBodyJSON(token, r); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}

	if err := srv.tokensRepository.GetToken(token); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}
	if err := srv.tokensRepository.DeleteToken(token); err != nil {
		beans.Logger().Error(err) //TODO format error
	}

	p := &models.ProfileInfo{
		UserID: token.UserID,
	}
	if err := srv.profilesRepository.SetProfileVerified(p); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}

	createdUserData := &queues.CreatedUserData{
		UserID: p.UserID,
	}
	if err := srv.mqClient.PublishUserCreatedMessage(createdUserData); err != nil {
		srv.WriteErrJSON(w, err) //TODO handle
		return
	}

	srv.WriteOkJSON(w, p)
}

func (srv *Server) updateEmail(w http.ResponseWriter, r *http.Request) {
	emailUpdate := &models.EmailUpdate{
		UserID: srv.parsePathUserID(r),
	}
	if string(srv.parseHeaderUserID(r)) != emailUpdate.UserID { //TODO
		srv.WriteErrJSON(w, srv.notAuthorizedErr)
		return
	}

	if err := srv.ParseBodyJSON(emailUpdate, r); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}

	if err := srv.profilesRepository.UpdateEmail(emailUpdate); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}

	srv.WriteOk(w)
}

func (srv *Server) updatePassword(w http.ResponseWriter, r *http.Request) {
	passwordUpdate := &models.PasswordUpdate{
		UserID: srv.parsePathUserID(r),
	}
	if string(srv.parseHeaderUserID(r)) != passwordUpdate.UserID { //TODO
		srv.WriteErrJSON(w, srv.notAuthorizedErr)
		return
	}

	if err := srv.ParseBodyJSON(passwordUpdate, r); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}

	if err := srv.profilesRepository.UpdatePassword(passwordUpdate); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}

	srv.WriteOk(w)
}

func (srv *Server) deleteProfile(w http.ResponseWriter, r *http.Request) {
	srv.WriteNotImplemented(w) //TODO
}
