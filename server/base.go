package server

import (
	"fmt"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/uptrace/bunrouter"
)

// Ensure via the compiler that the Server type is satisfied
var _ Server = (*baseServer)(nil)

type baseServer struct {
	host   string
	port   uint
	logger *log.Logger // No mutex needed. Charm implements one.

	router *bunrouter.Router
}

type Option func(server *baseServer)

func WithHost(host string) Option {
	return func(server *baseServer) {
		server.host = host
	}
}

func WithPort(port uint) Option {
	return func(server *baseServer) {
		server.port = port
	}
}

func New(logger *log.Logger, opts ...Option) *baseServer {
	router := bunrouter.New()

	server := &baseServer{
		host:   "localhost",
		port:   8080,
		logger: logger,

		router: router,
	}

	for _, opt := range opts {
		opt(server)
	}

	server.registerRoutes()

	return server
}

func (server *baseServer) Serve() error {
	addr := fmt.Sprintf("%s:%d", server.host, server.port)

	server.logger.Infof("server: started http server on http://%s", addr)
	return http.ListenAndServe(addr, server.router)
}
