package server

import (
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
		Addr:              listenAddress,
		Handler:           handler,
	}
	return &Server{server: server}
}

func (s Server) Start() error {
	err := s.server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

func (s Server) Stop(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}