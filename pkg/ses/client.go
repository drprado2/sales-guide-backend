package ses

import "context"

type (
	Client struct {
		SendEmail
		SendTemplatedEmail
	}

	SendEmail          func(ctx context.Context, builder *EmailRawInputBuilder) error
	SendTemplatedEmail func(ctx context.Context, builder EmailTemplateInputBuilder) error
)
