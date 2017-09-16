package hunter

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

func New() (opentracing.Tracer, io.Closer) {
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	injector := jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
	extractor := jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)
	zipkinSharedRPCSpan := jaeger.TracerOptions.ZipkinSharedRPCSpan(true)

	tracer, closer := jaeger.NewTracer(
		"grpc_event",
		jaeger.NewConstSampler(true),
		jaeger.NewNullReporter(),
		injector,
		extractor,
		zipkinSharedRPCSpan,
	)
	return tracer, closer
}

// cfg := config.Configuration{
// 	Sampler: &config.SamplerConfig{
// 		Type:              "const",
// 		Param:             1,
// 		SamplingServerURL: "localhost:5775",
// 	},
// 	Reporter: &config.ReporterConfig{
// 		LogSpans:            false,
// 		BufferFlushInterval: 1 * time.Second,
// 	},
// }

// tracer, _, err := cfg.New(
// 	"grpc_jaeger",
// 	config.Logger(jaeger.StdLogger),
// 	// config.Observer(rpcmetrics.NewObserver(
// 	// 	metricsFactory.Namespace("route", nil),
// 	// 	rpcmetrics.DefaultNameNormalizer)),
// )
// if err != nil {
// 	return tracer, err
// }
// return tracer, nil
