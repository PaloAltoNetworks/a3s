package bootstrap

import (
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// CloseRecorderHandler is the type of recorder closer handler
type CloseRecorderHandler func()

// ConfigureTracerWithURL returns a jaeger backed opentracing tracer from an URL.
func ConfigureTracerWithURL(tracerURL string, serviceName string) (CloseRecorderHandler, error) {

	if tracerURL == "" {
		return nil, nil
	}

	tracer, closeFunc, err := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  tracerURL,
		},
	}.NewTracer()

	if err != nil {
		return nil, err
	}

	opentracing.InitGlobalTracer(tracer)

	return func() { _ = closeFunc.Close() }, nil // nolint: errcheck
}
