package server

import (
	"fmt"
	"net/http"
	"net/url"

	"bode.fun/go/oga/auth"
	"github.com/uptrace/bunrouter"
)

func (server *baseServer) registerRoutes() {
	server.router.GET("/", func(w http.ResponseWriter, req bunrouter.Request) error {
		_, err := fmt.Fprint(w, "Hello, World!")
		return err
	})

	issuerURL, err := url.Parse("https://bode-fun.eu.auth0.com/")
	if err != nil {
		panic(err)
	}

	authMw, err := auth.New(issuerURL, "https://api.oga.bode.fun")
	if err != nil {
		panic(err)
	}

	server.router.Use(authMw.CheckJWT).GET("/protected", func(w http.ResponseWriter, req bunrouter.Request) error {
		_, err := fmt.Fprint(w, "Hello, Protected!")
		return err
	})
}
