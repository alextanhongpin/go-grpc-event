package auth0

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
