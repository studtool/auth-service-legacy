package api

import (
	"auth-service/beans"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/studtool/common/errs"
	"io/ioutil"
	"net/http"
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

func (srv *Server) parseUserId(r *http.Request) string {
	return r.Header.Get("X-User-Id")
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
