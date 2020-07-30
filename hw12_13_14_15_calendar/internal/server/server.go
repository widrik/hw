package server

import (
	"context"
	"net/http"
	"time"
)

type Settings struct {
	host string
	port int
}

type Server struct {
	server *http.Server
}

func CreateSettings(host string, port int) (*Settings, error) {
	return &Settings{host, port}, nil
}

func NewWebServer(handler http.Handler, listenAddress string) *Server {
	server := &http.Server{
		Addr:    listenAddress,
		Handler: handler,
	}
	return &Server{server: server}
}

func (srv Server) Start() error {
	err := srv.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (srv Server) Stop(timeout time.Duration) error {
	if err := srv.server.Shutdown(context.Background()); err != nil {
		return err
	}

	return nil
}
