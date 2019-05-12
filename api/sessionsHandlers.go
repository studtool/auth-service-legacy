package api

import (
	"github.com/studtool/auth-service/beans"
	"github.com/studtool/auth-service/models"
	"github.com/studtool/auth-service/utils"
	"github.com/studtool/common/consts"
	"net/http"
)

func (srv *Server) startSession(w http.ResponseWriter, r *http.Request) {
	profile := &models.Profile{}
	if err := srv.server.ParseBodyJSON(&profile.Credentials, r); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	if err := srv.profilesRepository.FindUserIdByCredentials(profile); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	session := &models.Session{
		UserID:     profile.UserID,
		ExpireTime: srv.tokenExpTimeCalc.Calculate(),
	}

	jwtClaims := &utils.AuthTokenAttributes{
		UserId:  session.UserID,
		ExpTime: session.ExpireTime,
	}
	if t, err := srv.authTokenManager.CreateToken(jwtClaims); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	} else {
		session.AuthToken = t
	}

	if t, err := srv.refreshTokenManager.CreateToken(); err != nil {
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
		srv.server.WriteErrJSON(w, srv.notAuthorizedErr)
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

	srv.server.SetUserID(w, attr.UserId)
	srv.server.WriteOk(w)
}

func (srv *Server) refreshSession(w http.ResponseWriter, r *http.Request) {
	session := &models.Session{
		RefreshToken: srv.server.ParseRefreshToken(r),
	}
	if err := srv.sessionsRepository.FindUserIdByRefreshToken(session); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	jwtClaims := &utils.AuthTokenAttributes{
		UserId:  session.UserID,
		ExpTime: session.ExpireTime,
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
	token := srv.server.ParseRefreshToken(r)

	if err := srv.sessionsRepository.DeleteSessionByRefreshToken(token); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	srv.server.WriteOk(w)
}

func (srv *Server) endAllSessions(w http.ResponseWriter, r *http.Request) {
	token := srv.server.ParseRefreshToken(r)

	if err := srv.sessionsRepository.DeleteAllSessionsByRefreshToken(token); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	srv.server.WriteOk(w)
}
