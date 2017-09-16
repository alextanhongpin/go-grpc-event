package hunter

import (
	// opentracing "github.com/opentracing/opentracing-go"

	"log"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func New() opentracing.Tracer {
	// zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	// injector := jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
	// extractor := jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)
	// zipkinSharedRPCSpan := jaeger.TracerOptions.ZipkinSharedRPCSpan(true)

	// tracer, closer := jaeger.NewTracer(
	// 	"grpc_event",
	// 	jaeger.NewConstSampler(true),
	// 	jaeger.NewNullReporter(),
	// 	injector,
	// 	extractor,
	// 	zipkinSharedRPCSpan,
	// )
	// trc := tracer.(opentracing.Tracer)
	// return trc, closer

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

	tracer, _, err := cfg.New(
		"grpc_jaeger",
		config.Logger(jaeger.StdLogger),
		// config.Observer(rpcmetrics.NewObserver(
		// 	metricsFactory.Namespace("route", nil),
		// 	rpcmetrics.DefaultNameNormalizer)),
	)
	if err != nil {
		log.Fatal(err)
	}
	return tracer
}
