package adapters

import (
	"context"
	"github.com/drprado2/react-redux-typescript/pkg/logs"
)

type (
	LogrusLogger struct{}
)

func (*LogrusLogger) Infof(ctx context.Context, text string, args ...interface{}) {
	logs.Logger(ctx).Infof(text, args)
}

func (*LogrusLogger) InfoWithFieldsf(ctx context.Context, fields map[string]interface{}, text string, args ...interface{}) {
	logs.Logger(ctx).WithFields(fields).Infof(text, args)
}

func (*LogrusLogger) Warnf(ctx context.Context, text string, args ...interface{}) {
	logs.Logger(ctx).Warnf(text, args)
}

func (*LogrusLogger) WarnWithFieldsf(ctx context.Context, fields map[string]interface{}, text string, args ...interface{}) {
	logs.Logger(ctx).WithFields(fields).Warnf(text, args)
}

func (*LogrusLogger) Errorf(ctx context.Context, text string, args ...interface{}) {
	logs.Logger(ctx).Errorf(text, args)
}

func (*LogrusLogger) ErrorWithFieldsf(ctx context.Context, fields map[string]interface{}, text string, args ...interface{}) {
	logs.Logger(ctx).WithFields(fields).Errorf(text, args)
}

func (*LogrusLogger) Debugf(ctx context.Context, text string, args ...interface{}) {
	logs.Logger(ctx).Debugf(text, args)
}

func (*LogrusLogger) DebugWithFieldsf(ctx context.Context, fields map[string]interface{}, text string, args ...interface{}) {
	logs.Logger(ctx).WithFields(fields).Debugf(text, args)
}

func (*LogrusLogger) Fatalf(ctx context.Context, text string, args ...interface{}) {
	logs.Logger(ctx).Fatalf(text, args)
}

func (*LogrusLogger) FatalWithFieldsf(ctx context.Context, fields map[string]interface{}, text string, args ...interface{}) {
	logs.Logger(ctx).WithFields(fields).Fatalf(text, args)
}
