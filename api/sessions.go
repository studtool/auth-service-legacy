package api

import (
	"net/http"
	"time"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/types"

	"github.com/studtool/auth-service/config"
	"github.com/studtool/auth-service/models"
	"github.com/studtool/auth-service/utils"
)

func (srv *Server) startSession(w http.ResponseWriter, r *http.Request) {
	profile := &models.Profile{}
	if err := srv.server.ParseBodyJSON(profile, r); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	if err := srv.profilesRepository.FindUserIdByCredentials(profile); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	}

	session := &models.Session{
		UserId:     profile.UserId,
		ExpireTime: types.DateTime(time.Now().Add(config.JwtExpTime.Value())),
	}

	jwtClaims := &utils.JwtClaims{
		UserId:  session.UserId,
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

	srv.server.WriteBodyJSON(w, http.StatusOK, session)
}

func (srv *Server) parseSession(w http.ResponseWriter, r *http.Request) {
	token := srv.server.ParseAuthToken(r)
	if token == consts.EmptyString {
		srv.server.WriteErrJSON(w, srv.noAuthTokenErr)
		return
	}

	claims, err := srv.authTokenManager.ParseToken(token)
	if err != nil {
		srv.server.WriteErrJSON(w, srv.noAuthTokenErr)
		return
	}

	srv.server.SetUserId(w, claims.UserId)
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

	jwtClaims := &utils.JwtClaims{
		UserId:  session.UserId,
		ExpTime: session.ExpireTime,
	}
	if t, err := srv.authTokenManager.CreateToken(jwtClaims); err != nil {
		srv.server.WriteErrJSON(w, err)
		return
	} else {
		session.AuthToken = t
	}

	srv.server.WriteBodyJSON(w, http.StatusOK, session)
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
