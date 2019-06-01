package rest

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/pprof"
	"strings"

	"github.com/go-http-utils/headers"
	"github.com/gorilla/mux"
	"github.com/mailru/easyjson"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/studtool/common/consts"
	"github.com/studtool/common/errs"
	"github.com/studtool/common/types"
)

func GetMetricsHandler() http.Handler {
	return promhttp.Handler()
}

func GetProfilerHandler() http.Handler {
	router := mux.NewRouter()

	router.HandleFunc("/", pprof.Index)
	router.HandleFunc("/cmdline", pprof.Cmdline)
	router.HandleFunc("/profile", pprof.Profile)
	router.HandleFunc("/symbol", pprof.Symbol)

	router.Handle("/goroutine", pprof.Handler("goroutine"))
	router.Handle("/heap", pprof.Handler("heap"))
	router.Handle("/threadcreate", pprof.Handler("threadcreate"))
	router.Handle("/block", pprof.Handler("block"))

	return router
}

func (srv *Server) GetRawBody(r *http.Request) ([]byte, *errs.Error) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errs.New(err)
	}
	return b, nil
}

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

func (srv *Server) ParseHeaderAuthToken(r *http.Request) string {
	t := r.Header.Get(headers.Authorization)

	const bearerLen = len("Bearer ")

	n := len(t)
	if n <= bearerLen {
		return consts.EmptyString
	}

	return t[bearerLen:]
}

func (srv *Server) ParseHeaderUserAgent(r *http.Request) string {
	return r.Header.Get(headers.UserAgent)
}

func (srv *Server) ParseHeaderXRealIP(r *http.Request) string {
	return r.Header.Get(headers.XRealIP)
}

func (srv *Server) ParseHeaderRefreshToken(r *http.Request) string {
	return r.Header.Get(RefreshTokenHeader)
}

func (srv *Server) ParseHeaderUserID(r *http.Request) types.ID {
	return types.ID(r.Header.Get(UserIDHeader))
}

func (srv *Server) SetHeaderUserID(w http.ResponseWriter, userID types.ID) {
	w.Header().Set(UserIDHeader, string(userID))
}

func (srv *Server) WriteOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (srv *Server) WriteOkJSON(w http.ResponseWriter, v easyjson.Marshaler) {
	srv.writeBodyJSON(w, http.StatusOK, v)
}

func (srv *Server) WriteOkRaw(w http.ResponseWriter, data []byte) {
	srv.writeBodyRaw(w, http.StatusOK, data)
}

func (srv *Server) WriteErrJSON(w http.ResponseWriter, err *errs.Error) {
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

	case errs.NotImplemented:
		srv.structLogger.Errorf("not implemented: %s", err.Message)
		w.WriteHeader(http.StatusInternalServerError)

	case errs.Internal:
		w.WriteHeader(http.StatusInternalServerError)

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
	w.Header().Set(headers.ContentType, "application/json")
	data, _ := easyjson.Marshal(v)
	_, _ = w.Write(data)
}

func (srv *Server) writeBodyRaw(w http.ResponseWriter, status int, data []byte) {
	w.WriteHeader(status)
	_, _ = w.Write(data)
}

func (srv *Server) writeErrBodyJSON(w http.ResponseWriter, status int, err *errs.Error) {
	w.WriteHeader(status)
	w.Header().Set(headers.ContentType, "application/json")
	_, _ = w.Write(err.JSON())
}

type LoggingResponseWriter struct {
	status   int
	severity ErrorSeverity
	writer   http.ResponseWriter
}

func (w *LoggingResponseWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *LoggingResponseWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w *LoggingResponseWriter) WriteHeader(status int) {
	w.status = status

	switch w.status {
	case http.StatusOK:
		w.severity = SeverityNone

	case http.StatusBadRequest:
		fallthrough
	case http.StatusNotFound:
		fallthrough
	case http.StatusUnprocessableEntity:
		fallthrough
	case http.StatusConflict:
		fallthrough
	case http.StatusForbidden:
		w.severity = SeverityLow

	default:
		w.severity = SeverityHigh
	}

	w.writer.WriteHeader(status)
}

func (w *LoggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.writer.(http.Hijacker).Hijack()
}

type ErrorSeverity int

const (
	SeverityNone = iota
	SeverityLow
	SeverityHigh
)

type APIClassifier interface {
	GetType(r *http.Request) string
}

const (
	APITypeNone = consts.EmptyString

	APITypePublic    = "public"
	APITypeProtected = "protected"
	APITypePrivate   = "private"
	APITypeInternal  = "internal"
)

type PathAPIClassifier struct{}

func NewPathAPIClassifier() *PathAPIClassifier {
	return &PathAPIClassifier{}
}

func (c *PathAPIClassifier) GetType(r *http.Request) string {
	path := r.RequestURI[len("/api/v"):]

	idx := 0
	for ; idx < len(path); idx++ {
		if path[idx] == '/' {
			break
		}
	}
	if idx == len(r.RequestURI) {
		return APITypeNone
	}

	path = path[idx+1:]
	if strings.HasPrefix(path, APITypePublic) {
		return APITypePublic
	}
	if strings.HasPrefix(path, APITypeProtected) {
		return APITypeProtected
	}
	if strings.HasPrefix(path, APITypePublic) {
		return APITypePrivate
	}
	if strings.HasPrefix(path, APITypeInternal) {
		return APITypeInternal
	}

	return APITypeNone
}
