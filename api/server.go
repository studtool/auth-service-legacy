package api

import (
	"auth-service/config"
	"fmt"
	"net/http"
)

type Server struct {
	server http.Server
}

func NewServer() *Server {
	srv := http.Server{
		Addr: fmt.Sprintf(":%s", config.ServerPort),
	}

	return &Server{
		server: srv,
	}
}

func (srv *Server) Run() {
	//TODO
}

func (srv *Server) Shutdown() {
	//TODO
}
