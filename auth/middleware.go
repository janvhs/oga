package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/uptrace/bunrouter"
	"golang.org/x/oauth2"
)

var (
	ErrNoTokenHeader error = errors.New("auth: no authorization header available")
)

type UserCtxKey struct{}

type middleware struct {
	jwtMiddleware *jwtmiddleware.JWTMiddleware
	issuerURL     *url.URL
	mu            sync.Mutex
}

func New(issuerURL *url.URL, audience ...string) (*middleware, error) {
	provider := jwks.NewProvider(issuerURL)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		audience,
		validator.WithAllowedClockSkew(time.Minute),
	)

	if err != nil {
		return nil, err
	}

	// TODO: Do error handler
	// TODO: Do custom claims

	jwtMiddleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
	)

	return &middleware{
		jwtMiddleware: jwtMiddleware,
		issuerURL:     issuerURL,
	}, nil
}

func (m *middleware) oidcAuth(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			// FIXME: Return a Error with Status code like uptrace does
			return ErrNoTokenHeader
		}

		// Get the access token from the header
		accessToken := strings.TrimPrefix(authHeader, "Bearer ")

		// TODO: Verify JWT. This could be problematic if the access token is a opaque string.
		// This will rely heavily on the auth2 extension for JWT Profiles.

		// Get user info from oidc endpoint to check if the access token is valid.
		// I do this because some providers like auth0 do not support the oauth2
		// access token introspection extension.
		// Therefore, this is a workaround.
		// The access token has to be requested with the "openid" scope for this to work.
		// FIXME: Test this implementation.
		provider, err := oidc.NewProvider(
			req.Context(),
			"https://bode-fun.eu.auth0.com/", // Without "/.well-known/openid-configuration"
		)
		if err != nil {
			// FIXME: Return a Error with Status code like uptrace does
			log.Fatal(errors.Join(errors.New("oidc"), err))
		}

		// Request the user's info from the provider.
		info, err := provider.UserInfo(req.Context(), oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: accessToken,
		}))

		if err != nil {
			// FIXME: Return a Error with Status code like uptrace does
			log.Fatal(errors.Join(errors.New("user info"), err))
		}

		// FIXME: Remove this
		log.Println(info)

		// TODO: Create a user struct
		ctxWithUser := context.WithValue(req.Context(), UserCtxKey{}, info)

		// Make the user available to the request chain.
		return next(w, req.WithContext(ctxWithUser))
	}
}

func (m *middleware) CheckJWT(next bunrouter.HandlerFunc) bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		// TODO: Is that needed? Learn about mutexes
		m.mu.Lock()
		defer m.mu.Unlock()

		encounteredError := true

		var nextHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			req.Request = r
			next(w, req)
		}

		m.jwtMiddleware.CheckJWT(nextHandler).ServeHTTP(w, req.Request)

		if encounteredError {
			return errors.New("auth: invalid access token")
		}

		return nil
	}
}
