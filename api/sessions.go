package api

import (
	"auth-service/models"
	"net/http"
)

func (srv *Server) startSession(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credentials
	if err := srv.parseRequestBody(&credentials, r); err != nil {
		srv.writeErrJSON(w, err)
		return
	}

	var session models.Session
	if err := srv.sessionsRepository.AddSession(&credentials, &session); err != nil {
		srv.writeErrJSON(w, err)
		return
	}

	srv.writeBodyJSON(w, http.StatusOK, &session)
}

func (srv *Server) refreshSession(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (srv *Server) endSession(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (srv *Server) endAllSessions(w http.ResponseWriter, r *http.Request) {
	//TODO
}
