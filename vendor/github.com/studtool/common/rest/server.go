package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/studtool/common/logs"
)

type ServerConfig struct {
	Host string
	Port string
}

type Server struct {
	server *http.Server
	logger *logs.Logger
}

func NewServer(c ServerConfig) *Server {
	return &Server{
		server: &http.Server{
			Addr: fmt.Sprintf("%s:%s", c.Host, c.Port),
		},
	}
}

func (srv *Server) SetHandler(h http.Handler) {
	srv.server.Handler = h
}

func (srv *Server) SetLogger(log *logs.Logger) {
	srv.logger = log
}

func (srv *Server) Run() error {
	return srv.server.ListenAndServe()
}

func (srv *Server) Shutdown() error {
	return srv.server.Shutdown(context.TODO())
}
