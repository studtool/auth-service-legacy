package api

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/types"
	"auth-service/utils"
	"net/http"
	"time"
)

func (srv *Server) startSession(w http.ResponseWriter, r *http.Request) {
	profile := &models.Profile{}
	if err := srv.parseRequestBody(profile, r); err != nil {
		srv.writeErrJSON(w, err)
		return
	}

	if err := srv.profilesRepository.FindUserIdByCredentials(profile); err != nil {
		srv.writeErrJSON(w, err)
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
		srv.writeErrJSON(w, err)
		return
	} else {
		session.AuthToken = t
	}

	if t, err := srv.refreshTokenManager.CreateToken(); err != nil {
		srv.writeErrJSON(w, err)
		return
	} else {
		session.RefreshToken = t
	}

	if err := srv.sessionsRepository.AddSession(session); err != nil {
		srv.writeErrJSON(w, err)
		return
	}

	srv.writeBodyJSON(w, http.StatusOK, session)
}

func (srv *Server) refreshSession(w http.ResponseWriter, r *http.Request) {
	session := &models.Session{
		RefreshToken: srv.parseRefreshToken(r),
	}
	if err := srv.sessionsRepository.FindUserIdByRefreshToken(session); err != nil {
		srv.writeErrJSON(w, err)
		return
	}

	jwtClaims := &utils.JwtClaims{
		UserId:  session.UserId,
		ExpTime: session.ExpireTime,
	}
	if t, err := srv.authTokenManager.CreateToken(jwtClaims); err != nil {
		srv.writeErrJSON(w, err)
		return
	} else {
		session.AuthToken = t
	}

	srv.writeBodyJSON(w, http.StatusOK, session)
}

func (srv *Server) endSession(w http.ResponseWriter, r *http.Request) {
	token := srv.parseRefreshToken(r)

	if err := srv.sessionsRepository.DeleteSessionByRefreshToken(token); err != nil {
		srv.writeErrJSON(w, err)
		return
	}

	srv.writeOk(w)
}

func (srv *Server) endAllSessions(w http.ResponseWriter, r *http.Request) {
	token := srv.parseRefreshToken(r)

	if err := srv.sessionsRepository.DeleteAllSessionsByRefreshToken(token); err != nil {
		srv.writeErrJSON(w, err)
		return
	}

	srv.writeOk(w)
}
