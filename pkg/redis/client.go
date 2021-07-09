package redis

import (
	"context"
	"time"
)

type (
	Client struct {
		SetKey
		UnmarshalFromKey
		GetFromKey
		DeleteKey
		SetKeyJson
	}

	SetKey           func(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	UnmarshalFromKey func(ctx context.Context, key string, target interface{}) error
	GetFromKey       func(ctx context.Context, key string) (*string, error)
	DeleteKey        func(ctx context.Context, key string) (bool, error)
	SetKeyJson       func(ctx context.Context, key string, value interface{}, ttl time.Duration) error
)
