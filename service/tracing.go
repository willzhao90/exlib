package service

import (
	"io"
	"log"

	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

// InitOpenTracing loads environment and returns a Jaeger tracer client
func InitOpenTracing() (opentracing.Tracer, io.Closer) {
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		// parsing errors might happen here, such as when we get a string where we expect a number
		log.Fatalf("Could not parse Jaeger env vars: %s", err.Error())
	}

	// TODO close tracer
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.Fatalf("Could not initialize jaeger tracer: %s", err.Error())
	}
	return tracer, closer
}
