package tracer

import (
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

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

// New returns a new tracer with default options if none is provided
func New(opts ...Option) (tracer opentracing.Tracer, err error) {
	options := Options{
		host:          "http://localhost:9411/api/v1/spans", // The zipkin http url
		name:          "grpc_event",                         // The namespace of the tracer we are running
		sameSpan:      true,                                 // same span can be set to true for RPC style spans (Zipkin V1) vs Node style (OpenTracing)
		traceID128Bit: true,                                 // make Tracer generate 128 bit traceID's for root spans.
	}

	for _, o := range opts {
		o(&options)
	}
	// Create a new collector
	collector, err := zipkin.NewHTTPCollector(options.host)
	if err != nil {
		return
	}

	// Create a new zipkin recorder
	recorder := zipkin.NewRecorder(collector, false, "127.0.0.1:8081", options.name)

	// Create a new tracer
	tracer, err = zipkin.NewTracer(
		recorder,
		zipkin.ClientServerSameSpan(options.sameSpan),
		zipkin.TraceID128Bit(options.traceID128Bit),
	)
	if err != nil {
		return
	}
	return
}