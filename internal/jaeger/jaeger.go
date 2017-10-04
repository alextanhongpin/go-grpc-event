package hunter

import (
	// opentracing "github.com/opentracing/opentracing-go"

	"context"
	"io"
	"log"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func New(tracerKind string) (opentracing.Tracer, io.Closer) {
	// zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	// injector := jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
	// extractor := jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)
	// zipkinSharedRPCSpan := jaeger.TracerOptions.ZipkinSharedRPCSpan(true)

	// return jaeger.NewTracer(
	// 	tracerKind,
	// 	jaeger.NewConstSampler(true),
	// 	jaeger.NewNullReporter(),
	// 	injector,
	// 	extractor,
	// 	zipkinSharedRPCSpan,
	// )

	// works!
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:              "const",
			Param:             1,
			SamplingServerURL: "localhost:5775",
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            false,
			BufferFlushInterval: 1 * time.Second,
		},
	}

	tracer, closer, err := cfg.New(
		tracerKind,
		config.Logger(jaeger.StdLogger),
		// config.Observer(rpcmetrics.NewObserver(
		// 	metricsFactory.Namespace("route", nil),
		// 	rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		log.Fatal(err)
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

// NewSpanFromContext reads the parent context and return a new child context
func NewSpanFromContext(ctx context.Context, name string) opentracing.Span {
	var parentCtx opentracing.SpanContext
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan != nil {
		parentCtx = parentSpan.Context()
	}
	return opentracing.GlobalTracer().StartSpan(name, opentracing.ChildOf(parentCtx))
}
