package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"io/ioutil"
	"net/http"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/errs"

	"github.com/studtool/auth-service/beans"
)

func (srv *Server) parseRequestBody(v easyjson.Unmarshaler, r *http.Request) *errs.Error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errs.NewBadFormatError(err.Error())
	}

	if err := easyjson.Unmarshal(b, v); err != nil {
		return errs.NewInvalidFormatError(err.Error())
	}

	return nil
}

func (srv *Server) setUserId(w http.ResponseWriter, userId string) {
	w.Header().Set("X-User-Id", userId)
}

func (srv *Server) parseUserId(r *http.Request) string {
	return r.Header.Get("X-User-Id")
}

func (srv *Server) parseAuthToken(r *http.Request) string {
	t := r.Header.Get("Authorization")

	const bearerLen = len("Bearer ")

	n := len(t)
	if n <= bearerLen {
		return consts.EmptyString
	}

	return t[bearerLen:]
}

func (srv *Server) parseRefreshToken(r *http.Request) string {
	return r.Header.Get("X-Refresh-Token")
}

func (srv *Server) parseProfileId(r *http.Request) string {
	return mux.Vars(r)["profile_id"]
}

func (srv *Server) writeOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) writeErrJSON(w http.ResponseWriter, err *errs.Error) {
	if err.Type == errs.Internal {
		beans.Logger.Error(err.Message)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch err.Type {
	case errs.BadFormat:
		srv.writeErrBodyJSON(w, http.StatusBadRequest, err)

	case errs.InvalidFormat:
		srv.writeErrBodyJSON(w, http.StatusUnprocessableEntity, err)

	case errs.Conflict:
		srv.writeErrBodyJSON(w, http.StatusConflict, err)

	case errs.NotFound:
		srv.writeErrBodyJSON(w, http.StatusNotFound, err)

	case errs.NotAuthorized:
		srv.writeErrBodyJSON(w, http.StatusUnauthorized, err)

	default:
		panic(fmt.Sprintf("no status code for error. Type: %d, Message: %s", err.Type, err.Message))
	}
}

func (srv *Server) writeBodyJSON(w http.ResponseWriter, status int, v easyjson.Marshaler) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	data, _ := easyjson.Marshal(v)
	_, _ = w.Write(data)
}

func (srv *Server) writeErrBodyJSON(w http.ResponseWriter, status int, err *errs.Error) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(err.JSON())
}
