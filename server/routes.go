package server

import (
	"fmt"
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (server *baseServer) registerRoutes() {
	server.router.GET("/", func(w http.ResponseWriter, req bunrouter.Request) error {
		_, err := fmt.Fprint(w, "Hello, World!")
		return err
	})
}
