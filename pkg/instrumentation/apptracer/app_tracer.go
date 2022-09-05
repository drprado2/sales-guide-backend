package apptracer

import (
	"context"
	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"runtime"
	"strings"
)

type TracerServiceInterface interface {
	SpanFromContext(ctx context.Context) (zipkin.Span, context.Context)
}

type TracerService struct {
	Tracer *zipkin.Tracer
	Client *zipkinhttp.Client
}

func (t *TracerService) SpanFromContext(ctx context.Context) (zipkin.Span, context.Context) {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	lastSlash := strings.LastIndex(f.Name(), "/")
	spanName := f.Name()[lastSlash+1:]

	span, spanCtx := t.Tracer.StartSpanFromContext(ctx, spanName)

	return span, spanCtx
}
