package tracer

// Usage
//
// trc, err := tracer.New(
// 	tracer.Name(*tracerKind),
// 	tracer.Host(*tracerHost), // "http://localhost:9411/api/v1/spans"
// )
// if err != nil {
// 	fmt.Printf("unable to create Zipkin tracer: %+v\n", err)
// 	os.Exit(-1)
// }
import (
	opentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
)

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
