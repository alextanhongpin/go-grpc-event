// package jaeger

// import (
// 	"time"

// 	opentracing "github.com/opentracing/opentracing-go"
// 	jaeger "github.com/uber/jaeger-client-go"

// 	"github.com/uber/jaeger-client-go/config"
// )

// func New() (opentracing.Tracer, error) {
// 	var trc opentracing.Tracer
// 	var err error
// 	cfg := config.Configuration{
// 		Sampler: &config.SamplerConfig{
// 			Type:              "const",
// 			Param:             1,
// 			SamplingServerURL: "localhost:5775",
// 		},
// 		Reporter: &config.ReporterConfig{
// 			LogSpans:            false,
// 			BufferFlushInterval: 1 * time.Second,
// 		},
// 	}

// 	tracer, _, err := cfg.New(
// 		"grpc_jaeger",
// 		config.Logger(jaeger.StdLogger),
// 		// config.Observer(rpcmetrics.NewObserver(
// 		// 	metricsFactory.Namespace("route", nil),
// 		// 	rpcmetrics.DefaultNameNormalizer)),
// 	)
// 	trc = tracer.(*opentracing.Tracer)
// 	if err != nil {
// 		return trc, err
// 	}
// 	return trc, nil
// }
