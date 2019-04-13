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
	UserIdHeader       = "X-User-Id"
	RefreshTokenHeader = "X-Refresh-Token"
)

func (srv *Server) SetUserId(w http.ResponseWriter, userId string) {
	w.Header().Set(UserIdHeader, userId)
}

func (srv *Server) ParseUserId(r *http.Request) string {
	return r.Header.Get(UserIdHeader)
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

func (srv *Server) WriteErrJSON(w http.ResponseWriter, err *errs.Error) {
	if err.Type == errs.Internal {
		srv.logger.Error(err.Message)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch err.Type {
	case errs.BadFormat:
		srv.WriteErrBodyJSON(w, http.StatusBadRequest, err)

	case errs.InvalidFormat:
		srv.WriteErrBodyJSON(w, http.StatusUnprocessableEntity, err)

	case errs.Conflict:
		srv.WriteErrBodyJSON(w, http.StatusConflict, err)

	case errs.NotFound:
		srv.WriteErrBodyJSON(w, http.StatusNotFound, err)

	case errs.NotAuthorized:
		srv.WriteErrBodyJSON(w, http.StatusUnauthorized, err)

	case errs.PermissionDenied:
		srv.WriteErrBodyJSON(w, http.StatusForbidden, err)

	default:
		panic(fmt.Sprintf("no status code for error. Type: %d, Message: %s", err.Type, err.Message))
	}
}

func (srv *Server) WriteBodyJSON(w http.ResponseWriter, status int, v easyjson.Marshaler) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	data, _ := easyjson.Marshal(v)
	_, _ = w.Write(data)
}

func (srv *Server) WriteErrBodyJSON(w http.ResponseWriter, status int, err *errs.Error) {
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
