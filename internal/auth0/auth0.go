package auth0

import (
	"net/http"

	auth0 "github.com/auth0-community/go-auth0"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Auth0 holds reference to the auth0 jwt validator
type Auth0 struct {
	Ref *auth0.JWTValidator
}

// Validate checks the request header for the authorization token and returns the token or error
func (a *Auth0) Validate(r *http.Request) (*jwt.JSONWebToken, error) {
	return a.Ref.ValidateRequest(r)
}
