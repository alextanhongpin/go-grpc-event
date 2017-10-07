package database

// Options represents the available options for the database
type Options struct {
	host     string
	db       string
	username string
	password string
	timeout  int
}

// Option is a function that returns a new option
type Option func(*Options)

// Host is a the database host
func Host(host string) Option {
	return func(o *Options) {
		o.host = host
	}
}

// Name is the database name
func Name(db string) Option {
	return func(o *Options) {
		o.db = db
	}
}

// Username is the database name
func Username(username string) Option {
	return func(o *Options) {
		o.username = username
	}
}

// Password is the database name
func Password(password string) Option {
	return func(o *Options) {
		o.password = password
	}
}

// TimeoutInSeconds is the time taken for mongodb to connect before it fails
func TimeoutInSeconds(timeout int) Option {
	return func(o *Options) {
		o.timeout = timeout
	}
}
