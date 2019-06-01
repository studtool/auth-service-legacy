package api

import (
	"net/http"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/types"

	"github.com/studtool/auth-service/beans"
	"github.com/studtool/auth-service/config"
	"github.com/studtool/auth-service/models"
	"github.com/studtool/auth-service/srvutils"
)

func (srv *Server) startSession(w http.ResponseWriter, r *http.Request) {
	profile := &models.Profile{}
	if err := srv.ParseBodyJSON(&profile.Credentials, r); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}

	if config.VerificationRequired.Value() {
		if err := srv.profilesRepository.FindVerifiedProfile(profile); err != nil {
			srv.WriteErrJSON(w, err)
			return
		}
	} else {
		if err := srv.profilesRepository.FindProfile(profile); err != nil {
			srv.WriteErrJSON(w, err)
			return
		}
	}

	session := &models.Session{
		UserID:     profile.UserID,
		ExpireTime: srv.tokenExpTimeCalc.Calculate(),
	}

	authAttr := &srvutils.AuthTokenAttributes{
		UserID:  session.UserID,
		ExpTime: session.ExpireTime,
	}
	if t, err := srv.authTokenManager.CreateToken(authAttr); err != nil {
		srv.WriteErrJSON(w, err)
		return
	} else {
		session.AuthToken = t
	}

	refreshAttr := &srvutils.RefreshTokenAttributes{
		UserID: session.UserID,
	}
	if t, err := srv.refreshTokenManager.CreateToken(refreshAttr); err != nil {
		srv.WriteErrJSON(w, err)
		return
	} else {
		session.RefreshToken = t
	}

	if err := srv.sessionsRepository.AddSession(session); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}

	srv.WriteOkJSON(w, session)
}

func (srv *Server) parseSession(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		srv.WriteOk(w)
	}

	token := srv.parseHeaderAuthToken(r)
	if token == consts.EmptyString {
		srv.WriteErrJSON(w, srv.notAuthorizedErr)
		return
	}

	attr, err := srv.authTokenManager.ParseToken(token)
	if err != nil {
		beans.Logger().Error(err) //TODO format error
		srv.WriteErrJSON(w, srv.notAuthorizedErr)
		return
	}
	if err := srv.tokenExpTimeCalc.Check(attr.ExpTime); err != nil {
		srv.WriteErrJSON(w, srv.notAuthorizedErr)
		return
	}

	srv.SetHeaderUserID(w, types.ID(attr.UserID)) //TODO
	srv.WriteOk(w)
}

func (srv *Server) refreshSession(w http.ResponseWriter, r *http.Request) {
	session := &models.Session{
		SessionID:    srv.parsePathSessionID(r),
		RefreshToken: srv.parseHeaderRefreshToken(r),
	}
	if err := srv.sessionsRepository.FindSession(session); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}

	jwtClaims := &srvutils.AuthTokenAttributes{
		UserID:  session.UserID,
		ExpTime: srv.tokenExpTimeCalc.Calculate(),
	}
	if t, err := srv.authTokenManager.CreateToken(jwtClaims); err != nil {
		srv.WriteErrJSON(w, err)
		return
	} else {
		session.AuthToken = t
	}

	srv.WriteOkJSON(w, session)
}

func (srv *Server) endSession(w http.ResponseWriter, r *http.Request) {
	session := &models.Session{
		SessionID: srv.parsePathSessionID(r),
		UserID:    string(srv.parseHeaderUserID(r)), //TODO
	}
	if err := srv.sessionsRepository.DeleteSessionBySessionID(session); err != nil {
		srv.WriteErrJSON(w, err)
		return
	}

	srv.WriteOk(w)
}

func (srv *Server) endAllSessions(w http.ResponseWriter, r *http.Request) {
	srv.WriteNotImplemented(w)
}
