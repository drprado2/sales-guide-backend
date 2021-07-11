package ses

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/drprado2/react-redux-typescript/pkg/awsconfig"
)

const (
	DefaultCharSet = "UTF-8"
)

var (
	client *ses.Client
)

func Setup(ctx context.Context) error {
	cfg, err := awsconfig.GetDefault(ctx)
	if err != nil {
		return err
	}

	client = ses.NewFromConfig(cfg)
	return nil
}

func SendEmailSvc(ctx context.Context, input *ses.SendEmailInput) error {
	_, err := client.SendEmail(ctx, input)
	return err
}

func SendTemplatedEmailSvc(ctx context.Context, input *ses.SendTemplatedEmailInput) error {
	_, err := client.SendTemplatedEmail(ctx, input)
	return err
}
