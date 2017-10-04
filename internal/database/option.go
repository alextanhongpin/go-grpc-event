package database

// Options represents the available options for the database
type Options struct {
	host string
	name string
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
func Name(name string) Option {
	return func(o *Options) {
		o.name = name
	}
}
