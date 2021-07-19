package domain

import (
	"context"
	apptracer2 "github.com/drprado2/react-redux-typescript/pkg/apptracer"
)

type (
	Logger interface {
		Infof(ctx context.Context, text string, args ...interface{})
		InfoWithFieldsf(ctx context.Context, fields map[string]interface{}, text string, args ...interface{})
		Warnf(ctx context.Context, text string, args ...interface{})
		WarnWithFieldsf(ctx context.Context, fields map[string]interface{}, text string, args ...interface{})
		Errorf(ctx context.Context, text string, args ...interface{})
		ErrorWithFieldsf(ctx context.Context, fields map[string]interface{}, text string, args ...interface{})
		Debugf(ctx context.Context, text string, args ...interface{})
		DebugWithFieldsf(ctx context.Context, fields map[string]interface{}, text string, args ...interface{})
		Fatalf(ctx context.Context, text string, args ...interface{})
		FatalWithFieldsf(ctx context.Context, fields map[string]interface{}, text string, args ...interface{})
	}
)

var (
	LoggerSvc Logger
	TracerSvc apptracer2.TracerService
)

func Setup(
	loggerSvc Logger,
	tracer apptracer2.TracerService,
) {
	LoggerSvc = loggerSvc
	TracerSvc = tracer
}
