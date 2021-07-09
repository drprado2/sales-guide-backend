package sm

import "context"

type (
	Client struct {
		GetValueFromKey
	}

	GetValueFromKey func(ctx context.Context, secret string, key string) (string, error)
)
