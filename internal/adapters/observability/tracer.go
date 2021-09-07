package observability

import (
	"context"
	"fmt"
	"github.com/drprado2/react-redux-typescript/configs"
	apptracer2 "github.com/drprado2/react-redux-typescript/pkg/apptracer"
	logs2 "github.com/drprado2/react-redux-typescript/pkg/logs"
	"github.com/opentracing/opentracing-go"
	zipkintracer "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	"github.com/openzipkin/zipkin-go/model"
	zipkinhttp "github.com/openzipkin/zipkin-go/reporter/http"
	"time"
)

func NewTracer(ctx context.Context) (*apptracer2.TracerService, *model.Endpoint, error) {
	envs := configs.Get()

	endpoint, err := zipkin.NewEndpoint("sales-guide", fmt.Sprintf("0.0.0.0:%v", envs.ServerPort))
	if err != nil {
		logs2.Logger(ctx).Errorf("unable to create local endpoint: %+v\n", err)
		return nil, nil, err
	}

	httpReporter := zipkinhttp.NewReporter(envs.ZipkinUrl, zipkinhttp.BatchInterval(time.Second*3))
	defer httpReporter.Close()

	tracer, err := zipkin.NewTracer(
		httpReporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithTraceID128Bit(true))
	if err != nil {
		logs2.Logger(ctx).Errorf("unable to create apptracer: %+v\n", err)
		return nil, nil, err
	}

	tracerService := &apptracer2.TracerService{
		Endpoint: endpoint,
		Tracer:   tracer,
	}

	opentracing.SetGlobalTracer(zipkintracer.Wrap(tracer))

	return tracerService, endpoint, nil
}
