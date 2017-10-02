package auth0

import (
	"net/http"

	auth0 "github.com/auth0-community/go-auth0"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// Options is a struct that represents the available options for the auth0
type Options struct {
	audience []string
	jwksURI  string
	issuer   string
}

// Option is a function that returns a closure with reference to options
type Option func(*Options)

// Audience represents list of urls that is allowed by the openid
func Audience(audience ...string) Option {
	return func(o *Options) {
		o.audience = audience
	}
}

// JWKSURI represents the JSON Web Keys endpoint
func JWKSURI(uri string) Option {
	return func(o *Options) {
		o.jwksURI = uri
	}
}

// Issuer represents the openid issue
func Issuer(iss string) Option {
	return func(o *Options) {
		o.issuer = iss
	}
}

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

// Auth0 holds reference to the auth0 jwt validator
type Auth0 struct {
	Ref *auth0.JWTValidator
}

// Validate checks the request header for the authorization token and returns the token or error
func (a *Auth0) Validate(r *http.Request) (*jwt.JSONWebToken, error) {
	return a.Ref.ValidateRequest(r)
}
