package tracer

// Options represents the configurable option for the tracer
type Options struct {
	host          string
	name          string
	sameSpan      bool
	traceID128Bit bool
}

// Option is a function that returns a pointer to the options
type Option func(*Options)

// Host sets the tracer host
func Host(host string) Option {
	return func(o *Options) {
		o.host = host
	}
}

// Name sets the tracer span name
func Name(name string) Option {
	return func(o *Options) {
		o.name = name
	}
}
