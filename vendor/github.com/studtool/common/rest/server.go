package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	mx := http.NewServeMux()

	mx.Handle("/metrics", promhttp.Handler())
	mx.Handle("/", h)

	srv.server.Handler = mx
}

func (srv *Server) SetLogger(log *logs.Logger) {
	srv.logger = log
}

func (srv *Server) Run() error {
	srv.logger.Info(fmt.Sprintf("started [%s]", srv.server.Addr))

	return srv.server.ListenAndServe()
}

func (srv *Server) Shutdown() error {
	srv.logger.Info("shutdown")

	return srv.server.Shutdown(context.TODO())
}
