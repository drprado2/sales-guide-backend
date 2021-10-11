package ctxvals

import (
	"context"
	"time"
)

const (
	TimezoneCtxKey   = "timezone"
	TimeOffsetCtxKey = "timeoffset"
	DefaultTimezone  = "America/Sao_Paulo"
	DefaultTimeOff   = -3
	LocationCtxKey   = "location"
)

var (
	DefaultLocation, _ = time.LoadLocation(DefaultTimezone)
)

func WithTimezone(ctx context.Context, timezone string) context.Context {
	return context.WithValue(ctx, TimezoneCtxKey, timezone)
}

func TimezoneOrDefault(ctx context.Context) string {
	if value := ctx.Value(TimezoneCtxKey); value != nil {
		return value.(string)
	}
	return DefaultTimezone
}

func WithTimeOffset(ctx context.Context, offset int) context.Context {
	return context.WithValue(ctx, TimeOffsetCtxKey, offset)
}

func TimeOffsetOrDefault(ctx context.Context) int {
	if value := ctx.Value(TimeOffsetCtxKey); value != nil {
		return value.(int)
	}
	return DefaultTimeOff
}

func WithLocation(ctx context.Context, location *time.Location) context.Context {
	return context.WithValue(ctx, LocationCtxKey, location)
}

func LocationOrDefault(ctx context.Context) *time.Location {
	if value := ctx.Value(LocationCtxKey); value != nil {
		return value.(*time.Location)
	}
	return DefaultLocation
}
