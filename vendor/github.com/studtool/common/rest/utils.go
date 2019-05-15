package rest

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/mailru/easyjson"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/errs"
)

func (srv *Server) ParseBodyJSON(v easyjson.Unmarshaler, r *http.Request) *errs.Error {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errs.NewBadFormatError(err.Error())
	}

	if err := easyjson.Unmarshal(b, v); err != nil {
		return errs.NewInvalidFormatError(err.Error())
	}

	return nil
}

const (
	UserIDHeader       = "X-User-Id"
	RefreshTokenHeader = "X-Refresh-Token"
)

const (
	UnauthorizedUserID = consts.EmptyString
)

func (srv *Server) SetUserID(w http.ResponseWriter, userID string) {
	w.Header().Set(UserIDHeader, userID)
}

func (srv *Server) ParseUserID(r *http.Request) string {
	return r.Header.Get(UserIDHeader)
}

func (srv *Server) ParseAuthToken(r *http.Request) string {
	t := r.Header.Get("Authorization")

	const bearerLen = len("Bearer ")

	n := len(t)
	if n <= bearerLen {
		return consts.EmptyString
	}

	return t[bearerLen:]
}

func (srv *Server) ParseRefreshToken(r *http.Request) string {
	return r.Header.Get(RefreshTokenHeader)
}

func (srv *Server) WriteOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) WriteOkJSON(w http.ResponseWriter, v easyjson.Marshaler) {
	srv.writeBodyJSON(w, http.StatusOK, v)
}

func (srv *Server) WriteErrJSON(w http.ResponseWriter, err *errs.Error) {
	if err.Type == errs.Internal {
		srv.logger.Error(err.Message)

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

	case errs.PermissionDenied:
		srv.writeErrBodyJSON(w, http.StatusForbidden, err)

	default:
		panic(fmt.Sprintf("no status code for error. Type: %d, Message: %s", err.Type, err.Message))
	}
}

func (srv *Server) WriteUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
}

func (srv *Server) WriteNotImplemented(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotImplemented)
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

type LoggingResponseWriter struct {
	status int
	writer http.ResponseWriter
}

func (w *LoggingResponseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *LoggingResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w *LoggingResponseWriter) WriteHeader(status int) {
	w.status = status
	w.writer.WriteHeader(status)
}

func (w *LoggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.writer.(http.Hijacker).Hijack()
}
