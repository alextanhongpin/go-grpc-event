package auth0

import (
	auth0 "github.com/auth0-community/go-auth0"
	jose "gopkg.in/square/go-jose.v2"
)

// New returns a new Auth0 struct
func New(opts ...Option) *Auth0 {
	options := Options{
		audience: []string{},
		jwksURI:  "",
		issuer:   "",
	}
	for _, o := range opts {
		o(&options)
	}
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: options.jwksURI})
	audience := options.audience
	configuration := auth0.NewConfiguration(client, audience, options.issuer, jose.RS256)
	validator := auth0.NewValidator(configuration)
	return &Auth0{
		Ref: validator,
	}
}
