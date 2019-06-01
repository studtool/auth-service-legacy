package api

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi"

	"github.com/studtool/common/rest"
	"github.com/studtool/common/types"
)

func (srv *Server) parseHeaderAuthToken(r *http.Request) string {
	return srv.ParseHeaderAuthToken(r)
}

func (srv *Server) parseHeaderUserID(r *http.Request) types.ID {
	return srv.ParseHeaderUserID(r)
}

func (srv *Server) parseHeaderRefreshToken(r *http.Request) string {
	return srv.ParseHeaderRefreshToken(r)
}

func (srv *Server) parsePathUserID(r *http.Request) string {
	return chi.URLParam(r, "user_id")
}

func (srv *Server) parsePathSessionID(r *http.Request) string {
	return chi.URLParam(r, "session_id")
}

type PathAPIClassifier struct {
	commonClassifier *rest.PathAPIClassifier
}

func NewPathAPIClassifier() *PathAPIClassifier {
	return &PathAPIClassifier{
		commonClassifier: rest.NewPathAPIClassifier(),
	}
}

func (c *PathAPIClassifier) GetType(r *http.Request) string {
	if strings.HasPrefix(r.RequestURI, "/api/internal/auth/session/") {
		return rest.APITypeInternal
	}
	return c.commonClassifier.GetType(r)
}
