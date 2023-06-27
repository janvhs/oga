package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/uptrace/bunrouter"
)

func (server *baseServer) registerRoutes() {
	server.router.GET("/", func(w http.ResponseWriter, req bunrouter.Request) error {
		_, err := fmt.Fprint(w, "Hello, World!")
		return err
	})

	protected := server.router.NewGroup("/protected").Use(func(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
		return func(w http.ResponseWriter, req bunrouter.Request) error {
			w.WriteHeader(http.StatusUnauthorized)
			err := errors.New("auth middleware: not implemented")
			fmt.Fprint(w, err.Error())
			return err
		}
	})

	protected.GET("", func(w http.ResponseWriter, req bunrouter.Request) error {
		_, err := fmt.Fprint(w, "Hello, Protected!")
		return err
	})
}
