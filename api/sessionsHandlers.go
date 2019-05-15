package api

import (
	"net/http"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/rest"

	"github.com/studtool/auth-service/beans"
	"github.com/studtool/auth-service/config"
	"github.com/studtool/auth-service/models"
	"github.com/studtool/auth-service/utils"
)

func (srv *Server) startSession(w http.ResponseWriter, r *http.Request) {
	profile := &models.Profile{}
	if err := srv.server.ParseBodyJSON(&profile.Credentials, r); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	if config.VerificationRequired.Value() {
		if err := srv.profilesRepository.FindVerifiedProfile(profile); err != nil {
			srv.server.WriteErrJSON(w, err)
			return
		}
	} else {
		if err := srv.profilesRepository.FindProfile(profile); err != nil {
			srv.server.WriteErrJSON(w, err)
			return
		}
	}

	session := &models.Session{
		UserID:     profile.UserID,
		ExpireTime: srv.tokenExpTimeCalc.Calculate(),
	}

	authAttr := &utils.AuthTokenAttributes{
		UserID:  session.UserID,
		ExpTime: session.ExpireTime,
	}
	if t, err := srv.authTokenManager.CreateToken(authAttr); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	} else {
		session.AuthToken = t
	}

	refreshAttr := &utils.RefreshTokenAttributes{
		UserID: session.UserID,
	}
	if t, err := srv.refreshTokenManager.CreateToken(refreshAttr); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	} else {
		session.RefreshToken = t
	}

	if err := srv.sessionsRepository.AddSession(session); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	srv.server.WriteOkJSON(w, session)
}

func (srv *Server) parseSession(w http.ResponseWriter, r *http.Request) {
	token := srv.server.ParseAuthToken(r)
	if token == consts.EmptyString {
		srv.server.SetUserID(w, rest.UnauthorizedUserID)
		srv.server.WriteOk(w)
		return
	}

	attr, err := srv.authTokenManager.ParseToken(token)
	if err != nil {
		beans.Logger().Error(err) //TODO format error
		srv.server.WriteErrJSON(w, srv.notAuthorizedErr)
		return
	}
	if err := srv.tokenExpTimeCalc.Check(attr.ExpTime); err != nil {
		srv.server.WriteErrJSON(w, srv.notAuthorizedErr)
		return
	}

	srv.server.SetUserID(w, attr.UserID)
	srv.server.WriteOk(w)
}

func (srv *Server) refreshSession(w http.ResponseWriter, r *http.Request) {
	session := &models.Session{
		SessionID:    srv.parseSessionID(r),
		RefreshToken: srv.server.ParseRefreshToken(r),
	}
	if err := srv.sessionsRepository.FindSession(session); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	jwtClaims := &utils.AuthTokenAttributes{
		UserID:  session.UserID,
		ExpTime: srv.tokenExpTimeCalc.Calculate(),
	}
	if t, err := srv.authTokenManager.CreateToken(jwtClaims); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	} else {
		session.AuthToken = t
	}

	srv.server.WriteOkJSON(w, session)
}

func (srv *Server) endSession(w http.ResponseWriter, r *http.Request) {
	session := &models.Session{
		SessionID: srv.parseSessionID(r),
		UserID:    srv.server.ParseUserID(r),
	}
	if err := srv.sessionsRepository.DeleteSessionBySessionID(session); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	srv.server.WriteOk(w)
}

func (srv *Server) endAllSessions(w http.ResponseWriter, r *http.Request) {
	srv.server.WriteNotImplemented(w)
}
