package apptracer

import (
	"context"
	"github.com/openzipkin/zipkin-go"
	zipkinmodel "github.com/openzipkin/zipkin-go/model"
	"runtime"
	"strings"
)

type TracerServiceInterface interface {
	SpanFromContext(ctx context.Context) (zipkin.Span, context.Context)
}

type TracerService struct {
	Endpoint *zipkinmodel.Endpoint
	Tracer   *zipkin.Tracer
}

func (t *TracerService) SpanFromContext(ctx context.Context) (zipkin.Span, context.Context) {
	pc := make([]uintptr, 10)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	lastSlash := strings.LastIndex(f.Name(), "/")
	spanName := f.Name()[lastSlash+1:]

	span, spanCtx := t.Tracer.StartSpanFromContext(ctx, spanName, zipkin.RemoteEndpoint(t.Endpoint))
	span.Tag("cid", ctx.Value("cid").(string))

	return span, spanCtx
}
