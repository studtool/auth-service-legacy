package api

import (
	"auth-service/config"
	"auth-service/errs"
	"github.com/mailru/easyjson"
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

func (srv *Server) writeOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) writeErrJSON(w http.ResponseWriter, err *errs.Error) {
	if err.Type == errs.Internal {
		config.Logger.Error(err.Message)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch err.Type {
	case errs.BadFormat:
		srv.writeBodyJSON(w, http.StatusBadRequest, err)

	case errs.InvalidFormat:
		srv.writeBodyJSON(w, http.StatusUnprocessableEntity, err)

	case errs.Conflict:
		srv.writeBodyJSON(w, http.StatusConflict, err)
	}
}

func (srv *Server) writeBodyJSON(w http.ResponseWriter, status int, v easyjson.Marshaler) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	data, _ := easyjson.Marshal(v)
	_, _ = w.Write(data)
}